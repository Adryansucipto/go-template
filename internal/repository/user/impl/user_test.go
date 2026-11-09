package impl

import (
	"context"
	"errors"
	"regexp"
	"testing"

	repositoryModel "go-template/internal/repository/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

// Skenario 1: Query GetUserRepository Sukses
// Ekspektasi: GORM menjalankan SELECT * FROM "public"."MST_USER" dan mengembalikan rows.
func TestGetUserRepository_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	repo := &Repository{
		DbGorm: gormDB,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "username", "password"}).
		AddRow(1, "Budi Santoso", "budi123", "hashed_pwd_1").
		AddRow(2, "Siti Rahma", "siti456", "hashed_pwd_2")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."MST_USER"`)).
		WillReturnRows(rows)

	users, err := repo.GetUserRepository(context.Background())

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Budi Santoso", users[0].Name)
	assert.Equal(t, "siti456", users[1].Username)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Skenario 2: Query GetUserRepository Database Error
// Ekspektasi: Database mengembalikan error -> repository me-return error.
func TestGetUserRepository_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	repo := &Repository{
		DbGorm: gormDB,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."MST_USER"`)).
		WillReturnError(errors.New("connection reset by peer"))

	users, err := repo.GetUserRepository(context.Background())

	assert.Error(t, err)
	assert.Empty(t, users)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Skenario 3: Insert CreateUserRepository Sukses
// Ekspektasi: GORM menjalankan INSERT INTO "public"."MST_USER" dan mengembalikan id baru.
func TestCreateUserRepository_Success(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	repo := &Repository{
		DbGorm: gormDB,
	}

	newUser := repositoryModel.User{
		Name:     "Andi Wijaya",
		Username: "andi123",
		Password: "hashed_password",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "public"."MST_USER"`)).
		WithArgs(newUser.Name, newUser.Username, newUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectCommit()

	created, err := repo.CreateUserRepository(context.Background(), newUser)

	assert.NoError(t, err)
	assert.Equal(t, 10, created.ID)
	assert.Equal(t, newUser.Name, created.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Skenario 4: Insert CreateUserRepository Gagal
// Ekspektasi: GORM menjalankan INSERT dan database mengembalikan error (misal duplicate username).
func TestCreateUserRepository_DBError(t *testing.T) {
	gormDB, mock := setupMockDB(t)

	repo := &Repository{
		DbGorm: gormDB,
	}

	newUser := repositoryModel.User{
		Name:     "Andi Wijaya",
		Username: "andi123",
		Password: "hashed_password",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "public"."MST_USER"`)).
		WithArgs(newUser.Name, newUser.Username, newUser.Password).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	mock.ExpectRollback()

	created, err := repo.CreateUserRepository(context.Background(), newUser)

	assert.Error(t, err)
	assert.Zero(t, created.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

