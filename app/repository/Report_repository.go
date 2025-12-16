package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"context"
	"database/sql"
)

type ReportRepository interface {
	GetStatics(ctx context.Context, filter StatisticsFilter) (*model.AchievementStatistics, error)
	GetStudentReport(ctx context.Context, studentID string) (*model.StudentReport, error)
	IsAdvisor(
		ctx context.Context,
		lecturerID string,
		studentID string,
	) (bool, error)
}

type StaticsReport struct {
	DB *sql.DB
}

type StatisticsFilter struct {
	StudentID  *string
	LecturerID *string
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &StaticsReport{DB: db}
}

func (r *StaticsReport) GetStatics(ctx context.Context, filter StatisticsFilter) (*model.AchievementStatistics, error) {
	where := ""
	args := []interface{}{}

	if filter.StudentID != nil {
		where = "WHERE ar.student_id = $1"
		args = append(args, *filter.StudentID)
	}

	if filter.LecturerID != nil {
		where = `
			JOIN students s ON s.id = ar.student_id
			WHERE s.advisor_id = $1
		`
		args = append(args, *filter.LecturerID)
	}

	query := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE ar.status = 'submitted') AS submitted,
			COUNT(*) FILTER (WHERE ar.status = 'verified') AS verified,
			COUNT(*) FILTER (WHERE ar.status = 'rejected') AS rejected
		FROM achievement_references ar
		` + where

	row := r.DB.QueryRowContext(ctx, query, args...)

	stats := new(model.AchievementStatistics)

	err := row.Scan(
		&stats.Total,
		&stats.Submitted,
		&stats.Verified,
		&stats.Rejected,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *StaticsReport) GetStudentReport(ctx context.Context, studentID string) (*model.StudentReport, error) {
	statsQuery := `
		SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE status = 'submitted') AS submitted,
			COUNT(*) FILTER (WHERE status = 'verified') AS verified,
			COUNT(*) FILTER (WHERE status = 'rejected') AS rejected
		FROM achievement_references
		WHERE student_id = $1
	`

	report := &model.StudentReport{
		StudentID: studentID,
	}

	err := r.DB.QueryRowContext(
		ctx,
		statsQuery,
		studentID,
	).Scan(
		&report.Total,
		&report.Submitted,
		&report.Verified,
		&report.Rejected,
	)

	if err != nil {
		return nil, err
	}

	listQuery := `
		SELECT
			ar.mongo_achievement_id,
			a.title,
			a.achievement_type,
			ar.status,
			ar.submitted_at,
			ar.verified_at
		FROM achievement_references ar
		JOIN achievements a
		  ON a.id = ar.mongo_achievement_id
		WHERE ar.student_id = $1
		ORDER BY ar.submitted_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, listQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.StudentAchievementItem
		if err := rows.Scan(
			&item.MongoAchievementID,
			&item.Title,
			&item.AchievementType,
			&item.Status,
			&item.SubmittedAt,
			&item.VerifiedAt,
		); err != nil {
			return nil, err
		}

		report.Achievements = append(report.Achievements, item)
	}

	return report, nil
}

func (r *StaticsReport) IsAdvisor(
	ctx context.Context,
	lecturerID string,
	studentID string,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM students
			WHERE id = $1
			  AND advisor_id = $2
		)
	`

	var exists bool
	err := r.DB.QueryRowContext(
		ctx,
		query,
		studentID,
		lecturerID,
	).Scan(&exists)

	return exists, err
}
