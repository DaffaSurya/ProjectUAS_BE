package testing

import (
	model "PROJECTUAS_BE/app/Model"
	"PROJECTUAS_BE/app/repository"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateStudent_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	userID := "user-123"

	mock.ExpectExec(`INSERT INTO student`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateStudent(userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentByUserID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	userID := "user-123"

	rows := sqlmock.NewRows([]string{
		"student_id", "user_id", "academic_year",
		"program_study", "username", "email", "full_name",
	}).AddRow(
		"student-1",
		userID,
		"2022",
		"Informatics",
		"testuser",
		"test@gmail.com",
		"Test User",
	)

	mock.ExpectQuery(`FROM students`).
		WithArgs(userID).
		WillReturnRows(rows)

	student, err := repo.GetStudentByUserID(userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if student.UserID != userID {
		t.Fatalf("expected userID %s", userID)
	}
}

func TestGetStudentByUserID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	userID := "not-found"

	mock.ExpectQuery(`FROM students`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	student, err := repo.GetStudentByUserID(userID)

	if err == nil {
		t.Fatal("expected error")
	}

	if student != nil {
		t.Fatal("expected nil student")
	}
}

func TestGetAllStudents_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	rows := sqlmock.NewRows([]string{
		"student_id", "user_id", "academic_year",
		"program_study", "username", "full_name", "email",
	}).AddRow(
		"student-1",
		"user-1",
		"2022",
		"Informatics",
		"user1",
		"User One",
		"user1@gmail.com",
	)

	mock.ExpectQuery(`FROM students`).
		WillReturnRows(rows)

	students, err := repo.GetAllStudents()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(students) != 1 {
		t.Fatalf("expected 1 student")
	}
}

func TestSubmitAchievement_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	ref := &model.AchievementReference{
		ID:                 "ref-1",
		StudentID:          "student-1",
		MongoAchievementID: "mongo-1",
		Status:             "submitted",
		SubmittedAt:        nil,
	}

	mock.ExpectExec(`INSERT INTO achievement_references`).
		WithArgs(
			ref.ID,
			ref.StudentID,
			ref.MongoAchievementID,
			ref.Status,
			ref.SubmittedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Submit(context.Background(), ref)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetStudentIDByUserID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	userID := "user-123"
	studentID := "student-123"

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(studentID)

	mock.ExpectQuery(`SELECT id FROM students`).
		WithArgs(userID).
		WillReturnRows(rows)

	result, err := repo.GetStudentIDByUserID(context.Background(), userID)

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if result != studentID {
		t.Fatalf("expected %s", studentID)
	}
}

func TestUpdateAdvisor_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewStudentRepository(db)

	mock.ExpectExec(`UPDATE students`).
		WithArgs("lecturer-1", "student-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.UpdateAdvisor(
		context.Background(),
		"student-1",
		"lecturer-1",
	)

	if err != nil {
		t.Fatalf("unexpected error")
	}
}
