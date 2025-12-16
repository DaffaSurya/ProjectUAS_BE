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
	GetallLectures(ctx context.Context) ([]*model.LecturerResponse, error)
	Getadvisees(ctx context.Context, lecturerID string) ([]*model.AdviseeResponse, error)
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

func (r *lecturePostGres) GetallLectures(ctx context.Context) ([]*model.LecturerResponse, error) {
	query := `
		SELECT
			l.id,
			l.user_id,
			l.department,
			u.username,
			u.email,
			u.full_name
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY u.full_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []*model.LecturerResponse

	for rows.Next() {
		var l model.LecturerResponse
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.Department,
			&l.Username,
			&l.Email,
			&l.FullName,
		); err != nil {
			return nil, err
		}
		lecturers = append(lecturers, &l)
	}

	return lecturers, nil
}

func (r *lecturePostGres) Getadvisees(ctx context.Context, lecturerID string) ([]*model.AdviseeResponse, error) {

	query := `
		SELECT 
			s.id,
			s.username,
			u.email,
			u.fullname
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
		ORDER BY u.fullname ASC
	`

	rows, err := r.db.QueryContext(ctx, query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.AdviseeResponse

	for rows.Next() {
		var a model.AdviseeResponse
		if err := rows.Scan(
			&a.ID,
			&a.Username,
			&a.Email,
			&a.Fullname,
		); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}

	return result, nil
}
