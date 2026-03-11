package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/group/domain/entities"
)

type GroupPostgreSQLRepository struct {
	conn *core.Conn_PostgreSQL
}

func NewGroupPostgreSQLRepository(conn *core.Conn_PostgreSQL) *GroupPostgreSQLRepository {
	return &GroupPostgreSQLRepository{conn: conn}
}

// ==================== GROUP OPERATIONS ====================

func (r *GroupPostgreSQLRepository) Save(group *entities.Group) error {
	query := `
		INSERT INTO groups (name, description, owner_id, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`

	err := r.conn.DB.QueryRow(
		query,
		group.Name,
		group.Description,
		group.OwnerID,
		group.IsActive,
		group.CreatedAt,
		group.UpdatedAt,
	).Scan(&group.ID)

	if err != nil {
		return fmt.Errorf("error al insertar grupo: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) GetByID(id int64) (*entities.Group, error) {
	query := `SELECT id, name, description, owner_id, is_active, created_at, updated_at 
			  FROM groups WHERE id = $1 AND is_active = true`

	row := r.conn.DB.QueryRow(query, id)

	var group entities.Group
	var description sql.NullString

	err := row.Scan(
		&group.ID,
		&group.Name,
		&description,
		&group.OwnerID,
		&group.IsActive,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al buscar grupo por ID: %v", err)
	}

	if description.Valid {
		group.Description = description.String
	}

	return &group, nil
}

func (r *GroupPostgreSQLRepository) GetByOwnerID(ownerID int64) ([]entities.Group, error) {
	query := `SELECT id, name, description, owner_id, is_active, created_at, updated_at 
			  FROM groups WHERE owner_id = $1 AND is_active = true ORDER BY created_at DESC`

	rows, err := r.conn.DB.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener grupos por owner: %v", err)
	}
	defer rows.Close()

	return r.scanGroups(rows)
}

func (r *GroupPostgreSQLRepository) GetByUserID(userID int64, limit, offset int, search string) ([]entities.GroupWithDetails, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1 AND gm.status = 'accepted' AND g.is_active = true`
	args := []interface{}{userID}

	if search != "" {
		countQuery += ` AND g.name ILIKE $2`
		args = append(args, "%"+search+"%")
	}

	var total int
	err := r.conn.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error al contar grupos: %v", err)
	}

	dataQuery := `
		SELECT g.id, g.name, g.description, g.owner_id, g.is_active, g.created_at, g.updated_at,
			   u.username as owner_username,
			   (SELECT COUNT(*) FROM group_members gm2 WHERE gm2.group_id = g.id AND gm2.status = 'accepted') as member_count
		FROM groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		INNER JOIN users u ON g.owner_id = u.id
		WHERE gm.user_id = $1 AND gm.status = 'accepted' AND g.is_active = true`

	dataArgs := []interface{}{userID}
	paramIdx := 2

	if search != "" {
		dataQuery += fmt.Sprintf(` AND g.name ILIKE $%d`, paramIdx)
		dataArgs = append(dataArgs, "%"+search+"%")
		paramIdx++
	}

	dataQuery += fmt.Sprintf(` ORDER BY g.created_at DESC LIMIT $%d OFFSET $%d`, paramIdx, paramIdx+1)
	dataArgs = append(dataArgs, limit, offset)

	rows, err := r.conn.DB.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener grupos: %v", err)
	}
	defer rows.Close()

	var groups []entities.GroupWithDetails
	for rows.Next() {
		var g entities.GroupWithDetails
		var description sql.NullString

		err := rows.Scan(
			&g.ID, &g.Name, &description, &g.OwnerID, &g.IsActive,
			&g.CreatedAt, &g.UpdatedAt, &g.OwnerUsername, &g.MemberCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error al escanear grupo: %v", err)
		}
		if description.Valid {
			g.Description = description.String
		}
		groups = append(groups, g)
	}

	return groups, total, nil
}

func (r *GroupPostgreSQLRepository) Update(group *entities.Group) error {
	query := `
		UPDATE groups 
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.conn.DB.Exec(query, group.Name, group.Description, group.UpdatedAt, group.ID)
	if err != nil {
		return fmt.Errorf("error al actualizar grupo: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) Delete(id int64) error {
	query := `UPDATE groups SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err := r.conn.DB.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar grupo: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) TransferOwnership(groupID, newOwnerID int64) error {
	tx, err := r.conn.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	// Update the group owner
	_, err = tx.Exec(`UPDATE groups SET owner_id = $1, updated_at = $2 WHERE id = $3`, newOwnerID, time.Now(), groupID)
	if err != nil {
		return fmt.Errorf("error al transferir propiedad: %v", err)
	}

	// Update the new owner's role
	_, err = tx.Exec(`UPDATE group_members SET role = 'owner' WHERE group_id = $1 AND user_id = $2`, groupID, newOwnerID)
	if err != nil {
		return fmt.Errorf("error al actualizar rol del nuevo owner: %v", err)
	}

	// Demote the old owner to admin
	_, err = tx.Exec(`UPDATE group_members SET role = 'admin' WHERE group_id = $1 AND user_id != $2 AND role = 'owner'`, groupID, newOwnerID)
	if err != nil {
		return fmt.Errorf("error al actualizar rol del antiguo owner: %v", err)
	}

	return tx.Commit()
}

// ==================== MEMBER OPERATIONS ====================

func (r *GroupPostgreSQLRepository) AddMember(member *entities.GroupMember) error {
	query := `
		INSERT INTO group_members (group_id, user_id, role, status, invited_by, invited_at, responded_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	err := r.conn.DB.QueryRow(
		query,
		member.GroupID,
		member.UserID,
		member.Role,
		member.Status,
		member.InvitedBy,
		member.InvitedAt,
		member.RespondedAt,
	).Scan(&member.ID)

	if err != nil {
		return fmt.Errorf("error al agregar miembro: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) GetMemberByGroupAndUser(groupID, userID int64) (*entities.GroupMember, error) {
	query := `SELECT id, group_id, user_id, role, status, invited_by, invited_at, responded_at 
			  FROM group_members WHERE group_id = $1 AND user_id = $2`

	row := r.conn.DB.QueryRow(query, groupID, userID)

	var member entities.GroupMember
	var invitedBy sql.NullInt64
	var respondedAt sql.NullTime
	var status string
	var role string

	err := row.Scan(
		&member.ID,
		&member.GroupID,
		&member.UserID,
		&role,
		&status,
		&invitedBy,
		&member.InvitedAt,
		&respondedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al buscar miembro: %v", err)
	}

	member.Status = entities.MemberStatus(status)
	member.Role = entities.MemberRole(role)
	if invitedBy.Valid {
		member.InvitedBy = &invitedBy.Int64
	}
	if respondedAt.Valid {
		member.RespondedAt = &respondedAt.Time
	}

	return &member, nil
}

func (r *GroupPostgreSQLRepository) GetMembersByGroupID(groupID int64) ([]entities.GroupMemberWithUser, error) {
	query := `
		SELECT gm.id, gm.group_id, gm.user_id, gm.role, gm.status, gm.invited_by, gm.invited_at, gm.responded_at,
			   u.username, u.name, u.email
		FROM group_members gm
		INNER JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = $1 AND gm.status = 'accepted'
		ORDER BY gm.role ASC, gm.invited_at ASC`

	rows, err := r.conn.DB.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener miembros: %v", err)
	}
	defer rows.Close()

	var members []entities.GroupMemberWithUser
	for rows.Next() {
		var member entities.GroupMemberWithUser
		var invitedBy sql.NullInt64
		var respondedAt sql.NullTime
		var status string
		var role string

		err := rows.Scan(
			&member.ID, &member.GroupID, &member.UserID, &role, &status,
			&invitedBy, &member.InvitedAt, &respondedAt,
			&member.Username, &member.Name, &member.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear miembro: %v", err)
		}

		member.Status = entities.MemberStatus(status)
		member.Role = entities.MemberRole(role)
		if invitedBy.Valid {
			member.InvitedBy = &invitedBy.Int64
		}
		if respondedAt.Valid {
			member.RespondedAt = &respondedAt.Time
		}

		members = append(members, member)
	}

	return members, nil
}

func (r *GroupPostgreSQLRepository) UpdateMemberStatus(groupID, userID int64, status entities.MemberStatus) error {
	query := `UPDATE group_members SET status = $1, responded_at = $2 WHERE group_id = $3 AND user_id = $4`
	_, err := r.conn.DB.Exec(query, status, time.Now(), groupID, userID)
	if err != nil {
		return fmt.Errorf("error al actualizar estado del miembro: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) UpdateMemberRole(groupID, userID int64, role entities.MemberRole) error {
	query := `UPDATE group_members SET role = $1 WHERE group_id = $2 AND user_id = $3`
	_, err := r.conn.DB.Exec(query, role, groupID, userID)
	if err != nil {
		return fmt.Errorf("error al actualizar rol del miembro: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) RemoveMember(groupID, userID int64) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.conn.DB.Exec(query, groupID, userID)
	if err != nil {
		return fmt.Errorf("error al remover miembro: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) GetPendingInvitations(userID int64) ([]entities.GroupMember, error) {
	query := `
		SELECT id, group_id, user_id, role, status, invited_by, invited_at, responded_at 
		FROM group_members 
		WHERE user_id = $1 AND status = 'pending'
		ORDER BY invited_at DESC`

	rows, err := r.conn.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener invitaciones: %v", err)
	}
	defer rows.Close()

	var members []entities.GroupMember
	for rows.Next() {
		var member entities.GroupMember
		var invitedBy sql.NullInt64
		var respondedAt sql.NullTime
		var status string
		var role string

		err := rows.Scan(
			&member.ID, &member.GroupID, &member.UserID, &role, &status,
			&invitedBy, &member.InvitedAt, &respondedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear invitación: %v", err)
		}

		member.Status = entities.MemberStatus(status)
		member.Role = entities.MemberRole(role)
		if invitedBy.Valid {
			member.InvitedBy = &invitedBy.Int64
		}
		if respondedAt.Valid {
			member.RespondedAt = &respondedAt.Time
		}

		members = append(members, member)
	}

	return members, nil
}

func (r *GroupPostgreSQLRepository) GetAcceptedMemberIDs(groupID int64) ([]int64, error) {
	query := `SELECT user_id FROM group_members WHERE group_id = $1 AND status = 'accepted'`
	rows, err := r.conn.DB.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener miembros aceptados: %v", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error al escanear ID: %v", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ==================== HELPERS ====================

func (r *GroupPostgreSQLRepository) scanGroups(rows *sql.Rows) ([]entities.Group, error) {
	var groups []entities.Group

	for rows.Next() {
		var group entities.Group
		var description sql.NullString

		err := rows.Scan(
			&group.ID, &group.Name, &description, &group.OwnerID,
			&group.IsActive, &group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear grupo: %v", err)
		}
		if description.Valid {
			group.Description = description.String
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar grupos: %v", err)
	}

	return groups, nil
}
