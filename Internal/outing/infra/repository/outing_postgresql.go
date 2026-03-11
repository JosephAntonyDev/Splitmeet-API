package repository

import (
	"database/sql"
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/outing/domain/entities"
)

type OutingPostgresql struct {
	db *sql.DB
}

func NewOutingPostgresql(db *sql.DB) *OutingPostgresql {
	return &OutingPostgresql{db: db}
}

// Outing operations

func (r *OutingPostgresql) Save(outing *entities.Outing) error {
	query := `
		INSERT INTO outings (name, description, category_id, group_id, creator_id, outing_date, split_type, total_amount, status, is_editable)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		outing.Name,
		outing.Description,
		outing.CategoryID,
		outing.GroupID,
		outing.CreatorID,
		outing.OutingDate,
		outing.SplitType,
		outing.TotalAmount,
		outing.Status,
		outing.IsEditable,
	).Scan(&outing.ID, &outing.CreatedAt, &outing.UpdatedAt)
}

func (r *OutingPostgresql) GetByID(id int64) (*entities.Outing, error) {
	query := `
		SELECT id, name, description, category_id, group_id, creator_id, outing_date, split_type, total_amount, status, is_editable, created_at, updated_at
		FROM outings WHERE id = $1
	`
	outing := &entities.Outing{}
	err := r.db.QueryRow(query, id).Scan(
		&outing.ID,
		&outing.Name,
		&outing.Description,
		&outing.CategoryID,
		&outing.GroupID,
		&outing.CreatorID,
		&outing.OutingDate,
		&outing.SplitType,
		&outing.TotalAmount,
		&outing.Status,
		&outing.IsEditable,
		&outing.CreatedAt,
		&outing.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("outing not found")
	}
	if err != nil {
		return nil, err
	}
	return outing, nil
}

func (r *OutingPostgresql) GetByIDWithDetails(id int64) (*entities.OutingWithDetails, error) {
	query := `
		SELECT o.id, o.name, o.description, o.category_id, o.group_id, o.creator_id, o.outing_date, o.split_type, o.total_amount, o.status, o.is_editable, o.created_at, o.updated_at,
			COALESCE(c.name, '') as category_name, COALESCE(g.name, '') as group_name, u.username as creator_username,
			(SELECT COUNT(*) FROM outing_participants WHERE outing_id = o.id) as participant_count,
			(SELECT COUNT(*) FROM outing_participants WHERE outing_id = o.id AND status = 'confirmed') as paid_count
		FROM outings o
		LEFT JOIN categories c ON o.category_id = c.id
		LEFT JOIN groups g ON o.group_id = g.id
		JOIN users u ON o.creator_id = u.id
		WHERE o.id = $1
	`
	outing := &entities.OutingWithDetails{}
	err := r.db.QueryRow(query, id).Scan(
		&outing.ID,
		&outing.Name,
		&outing.Description,
		&outing.CategoryID,
		&outing.GroupID,
		&outing.CreatorID,
		&outing.OutingDate,
		&outing.SplitType,
		&outing.TotalAmount,
		&outing.Status,
		&outing.IsEditable,
		&outing.CreatedAt,
		&outing.UpdatedAt,
		&outing.CategoryName,
		&outing.GroupName,
		&outing.CreatorUsername,
		&outing.ParticipantCount,
		&outing.PaidCount,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("outing not found")
	}
	if err != nil {
		return nil, err
	}
	return outing, nil
}

func (r *OutingPostgresql) GetByUserID(userID int64) ([]entities.OutingWithDetails, error) {
	query := `
		SELECT DISTINCT o.id, o.name, o.description, o.category_id, o.group_id, o.creator_id, o.outing_date, o.split_type, o.total_amount, o.status, o.is_editable, o.created_at, o.updated_at,
			COALESCE(c.name, '') as category_name, COALESCE(g.name, '') as group_name, u.username as creator_username,
			(SELECT COUNT(*) FROM outing_participants WHERE outing_id = o.id) as participant_count,
			(SELECT COUNT(*) FROM outing_participants WHERE outing_id = o.id AND status = 'confirmed') as paid_count
		FROM outings o
		LEFT JOIN categories c ON o.category_id = c.id
		LEFT JOIN groups g ON o.group_id = g.id
		JOIN users u ON o.creator_id = u.id
		LEFT JOIN outing_participants op ON o.id = op.outing_id AND op.user_id = $1
		WHERE o.creator_id = $1 
		   OR (op.user_id = $1 AND op.status != 'declined')
		ORDER BY o.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outings []entities.OutingWithDetails
	for rows.Next() {
		var outing entities.OutingWithDetails
		err := rows.Scan(
			&outing.ID,
			&outing.Name,
			&outing.Description,
			&outing.CategoryID,
			&outing.GroupID,
			&outing.CreatorID,
			&outing.OutingDate,
			&outing.SplitType,
			&outing.TotalAmount,
			&outing.Status,
			&outing.IsEditable,
			&outing.CreatedAt,
			&outing.UpdatedAt,
			&outing.CategoryName,
			&outing.GroupName,
			&outing.CreatorUsername,
			&outing.ParticipantCount,
			&outing.PaidCount,
		)
		if err != nil {
			return nil, err
		}
		outings = append(outings, outing)
	}
	return outings, nil
}

func (r *OutingPostgresql) GetByGroupID(groupID int64) ([]entities.OutingWithDetails, error) {
	query := `
		SELECT o.id, o.name, o.description, o.category_id, o.group_id, o.creator_id, o.outing_date, o.split_type, o.total_amount, o.status, o.is_editable, o.created_at, o.updated_at,
			COALESCE(c.name, '') as category_name, COALESCE(g.name, '') as group_name, u.username as creator_username,
			(SELECT COUNT(*) FROM outing_participants WHERE outing_id = o.id) as participant_count,
			(SELECT COUNT(*) FROM outing_participants WHERE outing_id = o.id AND status = 'confirmed') as paid_count
		FROM outings o
		LEFT JOIN categories c ON o.category_id = c.id
		LEFT JOIN groups g ON o.group_id = g.id
		JOIN users u ON o.creator_id = u.id
		WHERE o.group_id = $1
		ORDER BY o.created_at DESC
	`
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outings []entities.OutingWithDetails
	for rows.Next() {
		var outing entities.OutingWithDetails
		err := rows.Scan(
			&outing.ID,
			&outing.Name,
			&outing.Description,
			&outing.CategoryID,
			&outing.GroupID,
			&outing.CreatorID,
			&outing.OutingDate,
			&outing.SplitType,
			&outing.TotalAmount,
			&outing.Status,
			&outing.IsEditable,
			&outing.CreatedAt,
			&outing.UpdatedAt,
			&outing.CategoryName,
			&outing.GroupName,
			&outing.CreatorUsername,
			&outing.ParticipantCount,
			&outing.PaidCount,
		)
		if err != nil {
			return nil, err
		}
		outings = append(outings, outing)
	}
	return outings, nil
}

func (r *OutingPostgresql) Update(outing *entities.Outing) error {
	query := `
		UPDATE outings
		SET name = $1, description = $2, category_id = $3, outing_date = $4, split_type = $5, status = $6, is_editable = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`
	return r.db.QueryRow(query,
		outing.Name,
		outing.Description,
		outing.CategoryID,
		outing.OutingDate,
		outing.SplitType,
		outing.Status,
		outing.IsEditable,
		outing.ID,
	).Scan(&outing.UpdatedAt)
}

func (r *OutingPostgresql) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM outings WHERE id = $1", id)
	return err
}

func (r *OutingPostgresql) UpdateTotalAmount(outingID int64, amount float64) error {
	_, err := r.db.Exec("UPDATE outings SET total_amount = $1, updated_at = NOW() WHERE id = $2", amount, outingID)
	return err
}

func (r *OutingPostgresql) MarkAsCompleted(outingID int64) error {
	_, err := r.db.Exec("UPDATE outings SET status = 'completed', is_editable = false, updated_at = NOW() WHERE id = $1", outingID)
	return err
}

// Participant operations

func (r *OutingPostgresql) AddParticipant(participant *entities.OutingParticipant) error {
	query := `
		INSERT INTO outing_participants (outing_id, user_id, invited_by, status, amount_owed, custom_amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, joined_at
	`
	return r.db.QueryRow(query,
		participant.OutingID,
		participant.UserID,
		participant.InvitedBy,
		participant.Status,
		participant.AmountOwed,
		participant.CustomAmount,
	).Scan(&participant.ID, &participant.JoinedAt)
}

func (r *OutingPostgresql) GetParticipantByOutingAndUser(outingID, userID int64) (*entities.OutingParticipant, error) {
	query := `
		SELECT id, outing_id, user_id, status, amount_owed, custom_amount, joined_at
		FROM outing_participants WHERE outing_id = $1 AND user_id = $2
	`
	p := &entities.OutingParticipant{}
	err := r.db.QueryRow(query, outingID, userID).Scan(
		&p.ID,
		&p.OutingID,
		&p.UserID,
		&p.Status,
		&p.AmountOwed,
		&p.CustomAmount,
		&p.JoinedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("participant not found")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *OutingPostgresql) GetParticipantByID(participantID int64) (*entities.OutingParticipant, error) {
	query := `
		SELECT id, outing_id, user_id, status, amount_owed, custom_amount, joined_at
		FROM outing_participants WHERE id = $1
	`
	p := &entities.OutingParticipant{}
	err := r.db.QueryRow(query, participantID).Scan(
		&p.ID,
		&p.OutingID,
		&p.UserID,
		&p.Status,
		&p.AmountOwed,
		&p.CustomAmount,
		&p.JoinedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("participant not found")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *OutingPostgresql) GetParticipantsByOutingID(outingID int64) ([]entities.OutingParticipantWithUser, error) {
	query := `
		SELECT op.id, op.outing_id, op.user_id, op.status, op.amount_owed, op.custom_amount, op.joined_at,
			u.username, u.name, u.email
		FROM outing_participants op
		JOIN users u ON op.user_id = u.id
		WHERE op.outing_id = $1
	`
	rows, err := r.db.Query(query, outingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []entities.OutingParticipantWithUser
	for rows.Next() {
		var p entities.OutingParticipantWithUser
		err := rows.Scan(
			&p.ID,
			&p.OutingID,
			&p.UserID,
			&p.Status,
			&p.AmountOwed,
			&p.CustomAmount,
			&p.JoinedAt,
			&p.Username,
			&p.Name,
			&p.Email,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *OutingPostgresql) GetConfirmedParticipants(outingID int64) ([]entities.OutingParticipantWithUser, error) {
	query := `
		SELECT op.id, op.outing_id, op.user_id, op.status, op.amount_owed, op.custom_amount, op.joined_at,
			u.username, u.name, u.email
		FROM outing_participants op
		JOIN users u ON op.user_id = u.id
		WHERE op.outing_id = $1 AND op.status = 'confirmed'
	`
	rows, err := r.db.Query(query, outingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []entities.OutingParticipantWithUser
	for rows.Next() {
		var p entities.OutingParticipantWithUser
		err := rows.Scan(
			&p.ID,
			&p.OutingID,
			&p.UserID,
			&p.Status,
			&p.AmountOwed,
			&p.CustomAmount,
			&p.JoinedAt,
			&p.Username,
			&p.Name,
			&p.Email,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *OutingPostgresql) UpdateParticipantStatus(outingID, userID int64, status entities.ParticipantStatus) error {
	_, err := r.db.Exec("UPDATE outing_participants SET status = $1 WHERE outing_id = $2 AND user_id = $3", status, outingID, userID)
	return err
}

func (r *OutingPostgresql) UpdateParticipantAmountOwed(participantID int64, amount float64) error {
	_, err := r.db.Exec("UPDATE outing_participants SET amount_owed = $1 WHERE id = $2", amount, participantID)
	return err
}

func (r *OutingPostgresql) RemoveParticipant(outingID, userID int64) error {
	_, err := r.db.Exec("DELETE FROM outing_participants WHERE outing_id = $1 AND user_id = $2", outingID, userID)
	return err
}

// Item operations

func (r *OutingPostgresql) AddItem(item *entities.OutingItem) error {
	query := `
		INSERT INTO outing_items (outing_id, product_id, custom_name, custom_presentation, quantity, unit_price, is_shared)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, subtotal, created_at
	`
	return r.db.QueryRow(query,
		item.OutingID,
		item.ProductID,
		item.CustomName,
		item.CustomPresentation,
		item.Quantity,
		item.UnitPrice,
		item.IsShared,
	).Scan(&item.ID, &item.Subtotal, &item.CreatedAt)
}

func (r *OutingPostgresql) GetItemByID(itemID int64) (*entities.OutingItem, error) {
	query := `
		SELECT id, outing_id, product_id, custom_name, custom_presentation, quantity, unit_price, subtotal, is_shared, created_at
		FROM outing_items WHERE id = $1
	`
	item := &entities.OutingItem{}
	err := r.db.QueryRow(query, itemID).Scan(
		&item.ID,
		&item.OutingID,
		&item.ProductID,
		&item.CustomName,
		&item.CustomPresentation,
		&item.Quantity,
		&item.UnitPrice,
		&item.Subtotal,
		&item.IsShared,
		&item.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("item not found")
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *OutingPostgresql) GetItemsByOutingID(outingID int64) ([]entities.OutingItemWithProduct, error) {
	query := `
		SELECT i.id, i.outing_id, i.product_id, i.custom_name, i.custom_presentation, i.quantity, i.unit_price, i.subtotal, i.is_shared, i.created_at,
			COALESCE(p.name, '') as product_name, COALESCE(p.presentation, '') as product_presentation, COALESCE(p.size, '') as product_size
		FROM outing_items i
		LEFT JOIN products p ON i.product_id = p.id
		WHERE i.outing_id = $1
		ORDER BY i.created_at
	`
	rows, err := r.db.Query(query, outingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entities.OutingItemWithProduct
	for rows.Next() {
		var item entities.OutingItemWithProduct
		err := rows.Scan(
			&item.ID,
			&item.OutingID,
			&item.ProductID,
			&item.CustomName,
			&item.CustomPresentation,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.IsShared,
			&item.CreatedAt,
			&item.ProductName,
			&item.ProductPresentation,
			&item.ProductSize,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OutingPostgresql) UpdateItem(item *entities.OutingItem) error {
	query := `
		UPDATE outing_items
		SET custom_name = $1, custom_presentation = $2, quantity = $3, unit_price = $4, subtotal = $5, is_shared = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(query, item.CustomName, item.CustomPresentation, item.Quantity, item.UnitPrice, item.Subtotal, item.IsShared, item.ID)
	return err
}

func (r *OutingPostgresql) DeleteItem(itemID int64) error {
	_, err := r.db.Exec("DELETE FROM outing_items WHERE id = $1", itemID)
	return err
}

// Split operations

func (r *OutingPostgresql) AddItemSplit(split *entities.ItemSplit) error {
	query := `
		INSERT INTO item_splits (outing_item_id, participant_id, split_amount, percentage)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRow(query,
		split.OutingItemID,
		split.ParticipantID,
		split.SplitAmount,
		split.Percentage,
	).Scan(&split.ID)
}

func (r *OutingPostgresql) GetSplitsByItemID(itemID int64) ([]entities.ItemSplitWithUser, error) {
	query := `
		SELECT s.id, s.outing_item_id, s.participant_id, s.split_amount, s.percentage,
			op.user_id, u.username, u.name
		FROM item_splits s
		JOIN outing_participants op ON s.participant_id = op.id
		JOIN users u ON op.user_id = u.id
		WHERE s.outing_item_id = $1
	`
	rows, err := r.db.Query(query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []entities.ItemSplitWithUser
	for rows.Next() {
		var split entities.ItemSplitWithUser
		err := rows.Scan(
			&split.ID,
			&split.OutingItemID,
			&split.ParticipantID,
			&split.SplitAmount,
			&split.Percentage,
			&split.UserID,
			&split.Username,
			&split.Name,
		)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}
	return splits, nil
}

func (r *OutingPostgresql) GetSplitsByParticipantID(participantID int64) ([]entities.ItemSplit, error) {
	query := `
		SELECT id, outing_item_id, participant_id, split_amount, percentage
		FROM item_splits WHERE participant_id = $1
	`
	rows, err := r.db.Query(query, participantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []entities.ItemSplit
	for rows.Next() {
		var split entities.ItemSplit
		err := rows.Scan(
			&split.ID,
			&split.OutingItemID,
			&split.ParticipantID,
			&split.SplitAmount,
			&split.Percentage,
		)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}
	return splits, nil
}

func (r *OutingPostgresql) DeleteSplitsByItemID(itemID int64) error {
	_, err := r.db.Exec("DELETE FROM item_splits WHERE outing_item_id = $1", itemID)
	return err
}

func (r *OutingPostgresql) DeleteSplit(splitID int64) error {
	_, err := r.db.Exec("DELETE FROM item_splits WHERE id = $1", splitID)
	return err
}
