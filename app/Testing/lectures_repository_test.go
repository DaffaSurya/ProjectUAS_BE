package testing

import (
	"PROJECTUAS_BE/app/repository"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestVerifyAchievement_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewLecturesRepository(db)

	mock.ExpectExec(`UPDATE achievement_references`).
		WithArgs(
			"verified",
			"lecturer-1",
			nil,
			"mongo-1",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Verify(
		context.Background(),
		"mongo-1",
		"verified",
		"lecturer-1",
		nil,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyAchievement_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewLecturesRepository(db)

	mock.ExpectExec(`UPDATE achievement_references`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Verify(
		context.Background(),
		"mongo-x",
		"verified",
		"lecturer-1",
		nil,
	)

	if err != sql.ErrNoRows {
		t.Fatal("expected sql.ErrNoRows")
	}
}

func TestRejectAchievement_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewLecturesRepository(db)

	mock.ExpectExec(`UPDATE achievement_references`).
		WithArgs(
			"lecturer-1",
			"Dokumen tidak valid",
			"mongo-1",
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Reject(
		context.Background(),
		"mongo-1",
		"Dokumen tidak valid",
		"lecturer-1",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetHistory_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewLecturesRepository(db)

	rows := sqlmock.NewRows([]string{
		"id",
		"mongo_achievement_id",
		"student_id",
		"status",
		"submitted_at",
		"verified_at",
		"verified_by",
		"rejection_reason",
	}).AddRow(
		"ref-1",
		"mongo-1",
		"student-1",
		"submitted",
		nil,
		nil,
		nil,
		nil,
	)

	mock.ExpectQuery(`FROM achievement_references`).
		WithArgs("mongo-1").
		WillReturnRows(rows)

	history, err := repo.GetHistory(context.Background(), "mongo-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("expected 1 history")
	}
}

func TestGetAllLectures_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewLecturesRepository(db)

	rows := sqlmock.NewRows([]string{
		"id",
		"user_id",
		"department",
		"username",
		"email",
		"full_name",
	}).AddRow(
		"lect-1",
		"user-1",
		"Informatics",
		"lecturer1",
		"lect@gmail.com",
		"Lecturer One",
	)

	mock.ExpectQuery(`FROM lecturers`).
		WillReturnRows(rows)

	result, err := repo.GetallLectures(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 lecturer")
	}
}

func TestGetAdvisees_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewLecturesRepository(db)

	rows := sqlmock.NewRows([]string{
		"id",
		"username",
		"email",
		"fullname",
	}).AddRow(
		"student-1",
		"student1",
		"student@gmail.com",
		"Student One",
	)

	mock.ExpectQuery(`FROM students`).
		WithArgs("lecturer-1").
		WillReturnRows(rows)

	result, err := repo.Getadvisees(context.Background(), "lecturer-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 advisee")
	}
}
