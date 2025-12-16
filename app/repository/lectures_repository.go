package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"context"
	"database/sql"
)

type LecturesRepository interface {
	Verify(ctx context.Context, mongoAchievementID string, status string, verifiedBy string, reason *string) error
	Reject(ctx context.Context, mongoAchievementID string, reason string, rejectedBy string) error
	GetHistory(ctx context.Context, mongoAchievementID string) ([]*model.AchievementHistory, error)
}

type lecturePostGres struct {
	db *sql.DB
}

func NewLecturesRepository(db *sql.DB) LecturesRepository {
	return &lecturePostGres{db}
}

func (r *lecturePostGres) Verify(ctx context.Context, mongoAchievementID string, status string, verifiedBy string, reason *string) error {
	query := `
		UPDATE achievement_references
		SET status = $1,
		    verified_at = NOW(),
		    verified_by = $2,
		    rejection_reason = $3
		WHERE mongo_achievement_id = $4
		  AND status = 'submitted'
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		status,
		verifiedBy,
		reason,
		mongoAchievementID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *lecturePostGres) Reject(ctx context.Context, mongoAchievementID string, reason string, rejectedBy string) error {
	query := `
		UPDATE achievement_references
		SET status = 'rejected',
		    verified_at = NOW(),
		    verified_by = $1,
		    rejection_reason = $2
		WHERE mongo_achievement_id = $3
		  AND status = 'submitted'
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		rejectedBy,
		reason,
		mongoAchievementID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *lecturePostGres) GetHistory(ctx context.Context, mongoAchievementID string) ([]*model.AchievementHistory, error) {
	query := `
		SELECT
			id,
			mongo_achievement_id,
			student_id,
			status,
			submitted_at,
			verified_at,
			verified_by,
			rejection_reason
		FROM achievement_references
		WHERE mongo_achievement_id = $1
		ORDER BY submitted_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, mongoAchievementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*model.AchievementHistory

	for rows.Next() {
		var h model.AchievementHistory
		if err := rows.Scan(
			&h.ID,
			&h.MongoAchievementID,
			&h.StudentID,
			&h.Status,
			&h.SubmittedAt,
			&h.VerifiedAt,
			&h.VerifiedBy,
			&h.RejectionReason,
		); err != nil {
			return nil, err
		}
		histories = append(histories, &h)
	}

	return histories, nil

}


