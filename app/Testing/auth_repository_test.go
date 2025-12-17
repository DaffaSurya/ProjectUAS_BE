package testing

import (
	"PROJECTUAS_BE/app/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFindByEmail_Success(t *testing.T) {

	// ===== Setup mock DB =====
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	email := "test@gmail.com"

	rows := sqlmock.NewRows([]string{
		"id",
		"email",
		"password_hash",
		"role_id",
		"full_name",
	}).AddRow(
		"user-123",
		email,
		"hashedpassword",
		"role-1",
		"Test User",
	)

	mock.ExpectQuery(
		`SELECT id, email, password_hash, role_id, full_name
		 FROM users
		 WHERE email = \$1
		 LIMIT 1;`,
	).
		WithArgs(email).
		WillReturnRows(rows)

	// ===== Execute =====
	user, err := repo.FindByEmail(email)

	// ===== Assert =====
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}

	if user.Fullname != "Test User" {
		t.Errorf("expected fullname 'Test User', got %s", user.Fullname)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet mock expectations: %v", err)
	}
}

func TestFindByEmail_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	email := "notfound@gmail.com"

	mock.ExpectQuery(
		`SELECT id, email, password_hash, role_id, full_name
		 FROM users
		 WHERE email = \$1
		 LIMIT 1;`,
	).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByEmail(email)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user != nil {
		t.Fatal("expected user to be nil")
	}
}

func TestGetRoleByUserID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	userID := "user-123"
	expectedRole := "lecturer"

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(expectedRole)

	mock.ExpectQuery(`SELECT ro.name`).
		WithArgs(userID).
		WillReturnRows(rows)

	role, err := repo.GetRoleByUserID(userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if role != expectedRole {
		t.Fatalf("expected role %s, got %s", expectedRole, role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetRoleByUserID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	userID := "unknown-user"

	mock.ExpectQuery(`SELECT ro.name`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	role, err := repo.GetRoleByUserID(userID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if role != "" {
		t.Fatalf("expected empty role, got %s", role)
	}
}

func TestGetProfile_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	userID := "user-123"

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name",
	}).AddRow(
		userID,
		"testuser",
		"test@gmail.com",
		"Test User",
	)

	mock.ExpectQuery(`SELECT id, username, email, full_name`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetProfile(userID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != userID {
		t.Fatalf("expected id %s, got %s", userID, user.ID)
	}

	if user.Email != "test@gmail.com" {
		t.Fatalf("unexpected email: %s", user.Email)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewAuthRepository(db)

	userID := "not-found"

	mock.ExpectQuery(`SELECT id, username, email, full_name`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetProfile(userID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user != nil {
		t.Fatal("expected nil user")
	}
}
