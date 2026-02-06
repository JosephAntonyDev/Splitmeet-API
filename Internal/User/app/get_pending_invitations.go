package app

import (
	"database/sql"
	"time"
)

type GetPendingInvitations struct {
	db *sql.DB
}

func NewGetPendingInvitations(db *sql.DB) *GetPendingInvitations {
	return &GetPendingInvitations{db: db}
}

type OutingInvitation struct {
	InvitationID   int64      `json:"invitation_id"`
	OutingID       int64      `json:"outing_id"`
	OutingName     string     `json:"outing_name"`
	Description    string     `json:"description,omitempty"`
	OutingDate     *time.Time `json:"outing_date,omitempty"`
	CategoryName   string     `json:"category_name,omitempty"`
	CreatorID      int64      `json:"creator_id"`
	CreatorName    string     `json:"creator_name"`
	InvitedAt      time.Time  `json:"invited_at"`
	OutingStatus   string     `json:"outing_status"`
	IsAvailable    bool       `json:"is_available"`
	UnavailableMsg string     `json:"unavailable_msg,omitempty"`
}

type GroupInvitation struct {
	InvitationID   int64     `json:"invitation_id"`
	GroupID        int64     `json:"group_id"`
	GroupName      string    `json:"group_name"`
	Description    string    `json:"description,omitempty"`
	OwnerID        int64     `json:"owner_id"`
	OwnerName      string    `json:"owner_name"`
	InvitedAt      time.Time `json:"invited_at"`
	IsAvailable    bool      `json:"is_available"`
	UnavailableMsg string    `json:"unavailable_msg,omitempty"`
}

type PendingInvitations struct {
	OutingInvitations []OutingInvitation `json:"outing_invitations"`
	GroupInvitations  []GroupInvitation  `json:"group_invitations"`
	TotalPending      int                `json:"total_pending"`
}

func (uc *GetPendingInvitations) Execute(userID int64) (*PendingInvitations, error) {
	result := &PendingInvitations{
		OutingInvitations: []OutingInvitation{},
		GroupInvitations:  []GroupInvitation{},
	}

	// Obtener invitaciones a salidas pendientes
	outingQuery := `
		SELECT 
			op.id, op.outing_id, o.name, COALESCE(o.description, ''), o.outing_date,
			COALESCE(c.name, '') as category_name,
			o.creator_id, u.name as creator_name,
			op.joined_at, o.status
		FROM outing_participants op
		JOIN outings o ON op.outing_id = o.id
		JOIN users u ON o.creator_id = u.id
		LEFT JOIN categories c ON o.category_id = c.id
		WHERE op.user_id = $1 AND op.status = 'pending'
		ORDER BY op.joined_at DESC
	`

	rows, err := uc.db.Query(outingQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var inv OutingInvitation
		var outingDate sql.NullTime
		err := rows.Scan(
			&inv.InvitationID,
			&inv.OutingID,
			&inv.OutingName,
			&inv.Description,
			&outingDate,
			&inv.CategoryName,
			&inv.CreatorID,
			&inv.CreatorName,
			&inv.InvitedAt,
			&inv.OutingStatus,
		)
		if err != nil {
			return nil, err
		}

		if outingDate.Valid {
			inv.OutingDate = &outingDate.Time
		}

		// Verificar disponibilidad
		inv.IsAvailable = inv.OutingStatus == "active"
		if !inv.IsAvailable {
			switch inv.OutingStatus {
			case "completed":
				inv.UnavailableMsg = "Esta salida ya ha sido completada"
			case "cancelled":
				inv.UnavailableMsg = "Esta salida fue cancelada"
			default:
				inv.UnavailableMsg = "Esta salida ya no está disponible"
			}
		}

		result.OutingInvitations = append(result.OutingInvitations, inv)
	}

	// Obtener invitaciones a grupos pendientes
	groupQuery := `
		SELECT 
			gm.id, gm.group_id, g.name, COALESCE(g.description, ''),
			g.owner_id, u.name as owner_name, gm.invited_at
		FROM group_members gm
		JOIN groups g ON gm.group_id = g.id
		JOIN users u ON g.owner_id = u.id
		WHERE gm.user_id = $1 AND gm.status = 'pending'
		ORDER BY gm.invited_at DESC
	`

	rows2, err := uc.db.Query(groupQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var inv GroupInvitation
		err := rows2.Scan(
			&inv.InvitationID,
			&inv.GroupID,
			&inv.GroupName,
			&inv.Description,
			&inv.OwnerID,
			&inv.OwnerName,
			&inv.InvitedAt,
		)
		if err != nil {
			return nil, err
		}

		inv.IsAvailable = true // Los grupos siempre están disponibles si existen
		result.GroupInvitations = append(result.GroupInvitations, inv)
	}

	result.TotalPending = len(result.OutingInvitations) + len(result.GroupInvitations)

	return result, nil
}
