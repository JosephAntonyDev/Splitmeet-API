package repository

import (
	"database/sql"
	"errors"

	"github.com/JosephAntonyDev/splitmeet-api/internal/payment/domain/entities"
)

type PaymentPostgresql struct {
	db *sql.DB
}

func NewPaymentPostgresql(db *sql.DB) *PaymentPostgresql {
	return &PaymentPostgresql{db: db}
}

func (r *PaymentPostgresql) Create(payment *entities.Payment) error {
	query := `
		INSERT INTO payments (outing_id, participant_id, amount, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	return r.db.QueryRow(
		query,
		payment.OutingID,
		payment.ParticipantID,
		payment.Amount,
		payment.Status,
		payment.Notes,
		payment.CreatedAt,
		payment.UpdatedAt,
	).Scan(&payment.ID)
}

func (r *PaymentPostgresql) GetByID(id int64) (*entities.Payment, error) {
	query := `
		SELECT id, outing_id, participant_id, amount, status, paid_at, confirmed_by, notes, created_at, updated_at
		FROM payments
		WHERE id = $1`

	payment := &entities.Payment{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.OutingID,
		&payment.ParticipantID,
		&payment.Amount,
		&payment.Status,
		&payment.PaidAt,
		&payment.ConfirmedBy,
		&payment.Notes,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("payment not found")
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentPostgresql) GetByIDWithDetails(id int64) (*entities.PaymentWithDetails, error) {
	query := `
		SELECT 
			p.id, p.outing_id, p.participant_id, p.amount, p.status, p.paid_at, p.confirmed_by, p.notes, p.created_at, p.updated_at,
			o.name as outing_name,
			u.username as participant_username, u.name as participant_name,
			COALESCE(c.username, '') as confirmed_by_username
		FROM payments p
		JOIN outings o ON p.outing_id = o.id
		JOIN outing_participants op ON p.participant_id = op.id
		JOIN users u ON op.user_id = u.id
		LEFT JOIN users c ON p.confirmed_by = c.id
		WHERE p.id = $1`

	payment := &entities.PaymentWithDetails{}
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.OutingID,
		&payment.ParticipantID,
		&payment.Amount,
		&payment.Status,
		&payment.PaidAt,
		&payment.ConfirmedBy,
		&payment.Notes,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.OutingName,
		&payment.ParticipantUsername,
		&payment.ParticipantName,
		&payment.ConfirmedByUsername,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("payment not found")
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentPostgresql) Update(payment *entities.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, paid_at = $2, confirmed_by = $3, notes = $4, updated_at = $5
		WHERE id = $6`

	_, err := r.db.Exec(query, payment.Status, payment.PaidAt, payment.ConfirmedBy, payment.Notes, payment.UpdatedAt, payment.ID)
	return err
}

func (r *PaymentPostgresql) Delete(id int64) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PaymentPostgresql) GetByOutingID(outingID int64) ([]entities.PaymentWithDetails, error) {
	query := `
		SELECT 
			p.id, p.outing_id, p.participant_id, p.amount, p.status, p.paid_at, p.confirmed_by, p.notes, p.created_at, p.updated_at,
			o.name as outing_name,
			u.username as participant_username, u.name as participant_name,
			COALESCE(c.username, '') as confirmed_by_username
		FROM payments p
		JOIN outings o ON p.outing_id = o.id
		JOIN outing_participants op ON p.participant_id = op.id
		JOIN users u ON op.user_id = u.id
		LEFT JOIN users c ON p.confirmed_by = c.id
		WHERE p.outing_id = $1
		ORDER BY p.created_at DESC`

	return r.queryPaymentsWithDetails(query, outingID)
}

func (r *PaymentPostgresql) GetByParticipantID(participantID int64) ([]entities.PaymentWithDetails, error) {
	query := `
		SELECT 
			p.id, p.outing_id, p.participant_id, p.amount, p.status, p.paid_at, p.confirmed_by, p.notes, p.created_at, p.updated_at,
			o.name as outing_name,
			u.username as participant_username, u.name as participant_name,
			COALESCE(c.username, '') as confirmed_by_username
		FROM payments p
		JOIN outings o ON p.outing_id = o.id
		JOIN outing_participants op ON p.participant_id = op.id
		JOIN users u ON op.user_id = u.id
		LEFT JOIN users c ON p.confirmed_by = c.id
		WHERE p.participant_id = $1
		ORDER BY p.created_at DESC`

	return r.queryPaymentsWithDetails(query, participantID)
}

func (r *PaymentPostgresql) GetPendingByOutingID(outingID int64) ([]entities.PaymentWithDetails, error) {
	query := `
		SELECT 
			p.id, p.outing_id, p.participant_id, p.amount, p.status, p.paid_at, p.confirmed_by, p.notes, p.created_at, p.updated_at,
			o.name as outing_name,
			u.username as participant_username, u.name as participant_name,
			COALESCE(c.username, '') as confirmed_by_username
		FROM payments p
		JOIN outings o ON p.outing_id = o.id
		JOIN outing_participants op ON p.participant_id = op.id
		JOIN users u ON op.user_id = u.id
		LEFT JOIN users c ON p.confirmed_by = c.id
		WHERE p.outing_id = $1 AND p.status = 'pending'
		ORDER BY p.created_at DESC`

	return r.queryPaymentsWithDetails(query, outingID)
}

func (r *PaymentPostgresql) GetPendingByOutingAndParticipant(outingID, participantID int64) (*entities.Payment, error) {
	query := `
		SELECT id, outing_id, participant_id, amount, status, paid_at, confirmed_by, notes, created_at, updated_at
		FROM payments
		WHERE outing_id = $1 AND participant_id = $2 AND status = 'pending'
		ORDER BY created_at ASC
		LIMIT 1`

	payment := &entities.Payment{}
	err := r.db.QueryRow(query, outingID, participantID).Scan(
		&payment.ID,
		&payment.OutingID,
		&payment.ParticipantID,
		&payment.Amount,
		&payment.Status,
		&payment.PaidAt,
		&payment.ConfirmedBy,
		&payment.Notes,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no pending payment found for this participant in this outing")
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentPostgresql) queryPaymentsWithDetails(query string, arg interface{}) ([]entities.PaymentWithDetails, error) {
	rows, err := r.db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []entities.PaymentWithDetails
	for rows.Next() {
		var p entities.PaymentWithDetails
		err := rows.Scan(
			&p.ID,
			&p.OutingID,
			&p.ParticipantID,
			&p.Amount,
			&p.Status,
			&p.PaidAt,
			&p.ConfirmedBy,
			&p.Notes,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.OutingName,
			&p.ParticipantUsername,
			&p.ParticipantName,
			&p.ConfirmedByUsername,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *PaymentPostgresql) GetSummaryByOutingID(outingID int64) (*entities.PaymentSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0) as paid_amount,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending_amount,
			COUNT(*) as payments_count,
			COUNT(CASE WHEN status = 'paid' THEN 1 END) as paid_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count
		FROM payments
		WHERE outing_id = $1`

	summary := &entities.PaymentSummary{OutingID: outingID}
	err := r.db.QueryRow(query, outingID).Scan(
		&summary.TotalAmount,
		&summary.PaidAmount,
		&summary.PendingAmount,
		&summary.PaymentsCount,
		&summary.PaidCount,
		&summary.PendingCount,
	)

	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *PaymentPostgresql) GetParticipantIDByOutingAndUser(outingID, userID int64) (int64, error) {
	query := `SELECT id FROM outing_participants WHERE outing_id = $1 AND user_id = $2`
	var participantID int64
	err := r.db.QueryRow(query, outingID, userID).Scan(&participantID)
	if err == sql.ErrNoRows {
		return 0, errors.New("participant not found in outing")
	}
	return participantID, err
}

func (r *PaymentPostgresql) IsParticipantInOuting(outingID, participantID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM outing_participants WHERE outing_id = $1 AND id = $2)`
	var exists bool
	err := r.db.QueryRow(query, outingID, participantID).Scan(&exists)
	return exists, err
}

func (r *PaymentPostgresql) GetOutingTotalAmount(outingID int64) (float64, error) {
	query := `SELECT COALESCE(total_amount, 0) FROM outings WHERE id = $1`
	var total float64
	err := r.db.QueryRow(query, outingID).Scan(&total)
	if err == sql.ErrNoRows {
		return 0, errors.New("outing not found")
	}
	return total, err
}

func (r *PaymentPostgresql) GetTotalConfirmedPayments(outingID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE outing_id = $1 AND status = 'paid'`
	var total float64
	err := r.db.QueryRow(query, outingID).Scan(&total)
	return total, err
}

func (r *PaymentPostgresql) GetParticipantAmountOwed(outingID, participantID int64) (float64, error) {
	query := `SELECT COALESCE(amount_owed, 0) FROM outing_participants WHERE outing_id = $1 AND id = $2`
	var amount float64
	err := r.db.QueryRow(query, outingID, participantID).Scan(&amount)
	if err == sql.ErrNoRows {
		return 0, errors.New("participant not found in outing")
	}
	return amount, err
}

func (r *PaymentPostgresql) GetConfirmedParticipantCount(outingID int64) (int, error) {
	query := `SELECT COUNT(*) FROM outing_participants WHERE outing_id = $1 AND status = 'confirmed'`
	var count int
	err := r.db.QueryRow(query, outingID).Scan(&count)
	return count, err
}

func (r *PaymentPostgresql) CancelPendingPayments(outingID int64) error {
	query := `UPDATE payments SET status = 'cancelled', updated_at = NOW() WHERE outing_id = $1 AND status = 'pending'`
	_, err := r.db.Exec(query, outingID)
	return err
}
