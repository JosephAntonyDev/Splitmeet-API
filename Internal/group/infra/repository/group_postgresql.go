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

func (r *GroupPostgreSQLRepository) GetByUserID(userID int64) ([]entities.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.owner_id, g.is_active, g.created_at, g.updated_at 
		FROM groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1 AND gm.status = 'accepted' AND g.is_active = true
		ORDER BY g.created_at DESC`

	rows, err := r.conn.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener grupos del usuario: %v", err)
	}
	defer rows.Close()

	return r.scanGroups(rows)
}

func (r *GroupPostgreSQLRepository) Update(group *entities.Group) error {
	query := `
		UPDATE groups 
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4`

	_, err := r.conn.DB.Exec(
		query,
		group.Name,
		group.Description,
		group.UpdatedAt,
		group.ID,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar grupo: %v", err)
	}
	return nil
}

func (r *GroupPostgreSQLRepository) Delete(id int64) error {
	// Soft delete - solo marca como inactivo
	query := `UPDATE groups SET is_active = false, updated_at = $1 WHERE id = $2`

	_, err := r.conn.DB.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error al eliminar grupo: %v", err)
	}
	return nil
}

// ==================== MEMBER OPERATIONS ====================

func (r *GroupPostgreSQLRepository) AddMember(member *entities.GroupMember) error {
	query := `
		INSERT INTO group_members (group_id, user_id, status, invited_by, invited_at, responded_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`

	err := r.conn.DB.QueryRow(
		query,
		member.GroupID,
		member.UserID,
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
	query := `SELECT id, group_id, user_id, status, invited_by, invited_at, responded_at 
			  FROM group_members WHERE group_id = $1 AND user_id = $2`

	row := r.conn.DB.QueryRow(query, groupID, userID)

	var member entities.GroupMember
	var invitedBy sql.NullInt64
	var respondedAt sql.NullTime
	var status string

	err := row.Scan(
		&member.ID,
		&member.GroupID,
		&member.UserID,
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
		SELECT gm.id, gm.group_id, gm.user_id, gm.status, gm.invited_by, gm.invited_at, gm.responded_at,
			   u.username, u.name, u.email
		FROM group_members gm
		INNER JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = $1
		ORDER BY gm.invited_at ASC`

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

		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&status,
			&invitedBy,
			&member.InvitedAt,
			&respondedAt,
			&member.Username,
			&member.Name,
			&member.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear miembro: %v", err)
		}

		member.Status = entities.MemberStatus(status)
		if invitedBy.Valid {
			member.InvitedBy = &invitedBy.Int64
		}
		if respondedAt.Valid {
			member.RespondedAt = &respondedAt.Time
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar miembros: %v", err)
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
		SELECT id, group_id, user_id, status, invited_by, invited_at, responded_at 
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

		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&status,
			&invitedBy,
			&member.InvitedAt,
			&respondedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear invitación: %v", err)
		}

		member.Status = entities.MemberStatus(status)
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

// ==================== HELPERS ====================

func (r *GroupPostgreSQLRepository) scanGroups(rows *sql.Rows) ([]entities.Group, error) {
	var groups []entities.Group

	for rows.Next() {
		var group entities.Group
		var description sql.NullString

		err := rows.Scan(
			&group.ID,
			&group.Name,
			&description,
			&group.OwnerID,
			&group.IsActive,
			&group.CreatedAt,
			&group.UpdatedAt,
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
