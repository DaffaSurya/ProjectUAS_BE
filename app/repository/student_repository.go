package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"context"
	"database/sql"
	"errors"
)

type StudentRepository interface {
	CreateStudent(userID string) error
	GetStudentByUserID(userID string) (*model.Student, error)
	GetAllStudents() ([]model.Student, error)
	Submit(ctx context.Context, ref *model.AchievementReference) error
	GetStudentIDByUserID(ctx context.Context, userID string) (string, error)
	UpdateAdvisor(ctx context.Context, studentID string, advisorID string) error
}

type StudentPostgres struct {
	DB *sql.DB
}

	func NewStudentRepository(db *sql.DB) StudentRepository {
		return &StudentPostgres{DB: db}
	}

func (r *StudentPostgres) CreateStudent(userID string) error {
	query := `
		INSERT INTO student (user_id)
		VALUES ($1)
	`
	_, err := r.DB.Exec(query, userID)
	return err
}

func (r *StudentPostgres) GetAllStudents() ([]model.Student, error) {
	query := `
		SELECT
			s.id AS student_id,
			s.user_id,
			s.academic_year,
			s.program_study,
			u.username,
			u.full_name,
			u.email
		FROM students s	
		JOIN users u ON u.id = s.user_id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student

	for rows.Next() {
		student := model.Student{}
		err := rows.Scan(
			&student.StudentID,
			&student.UserID,
			&student.AcademicYear,
			&student.ProgramStudy,
			&student.Username,
			&student.Fullname,
			&student.Email,
		)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentPostgres) GetStudentByUserID(userID string) (*model.Student, error) {
	query := `
		SELECT
			s.id AS student_id,
			s.user_id,
			s.academic_year,
			s.program_study,
			u.username,
			u.email,
			u.full_name
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.user_id = $1
	`

	row := r.DB.QueryRow(query, userID)

	student := model.Student{}
	err := row.Scan(
		&student.StudentID,
		&student.UserID,
		&student.AcademicYear,
		&student.ProgramStudy,
		&student.Username,
		&student.Email,
		&student.Fullname,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("student not found")
	}

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *StudentPostgres) Submit(ctx context.Context, ref *model.AchievementReference) error {
	query := `
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, submitted_at)
        VALUES  ($1, $2, $3, $4, $5)`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.SubmittedAt,
	)

	return err
}

func (r *StudentPostgres) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	var studentID string

	query := `
		SELECT id
		FROM students
		WHERE user_id = $1
	`

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&studentID)
	if err != nil {
		return "", err
	}

	return studentID, nil
}

func (r *StudentPostgres) UpdateAdvisor(ctx context.Context, studentID string, advisorID string) error {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		advisorID,
		studentID,
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


