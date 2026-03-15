package cmd

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IsaacDSC/migrations/migration"
)

func init() {
	log.SetOutput(os.NewFile(0, os.DevNull)) // suppress log output in tests
}

func TestUp_NoMigrationsToApply(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	state := []migration.Migrate{
		{Version: 1, Up: func(*sql.Tx) error { return nil }, Down: func(*sql.Tx) error { return nil }},
	}
	dbVersion := 1 // already at version 1, so state[1:] is empty

	insertCalled := false
	insertFn := func(tx *sql.Tx, version int) error {
		insertCalled = true
		return nil
	}

	err = Up(db, insertFn, dbVersion, state)
	if err != nil {
		t.Fatalf("Up: %v", err)
	}

	if insertCalled {
		t.Error("insertFn should not be called when there are no migrations to apply")
	}
	mock.ExpectationsWereMet()
}

func TestUp_AppliesMigrations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	state := []migration.Migrate{
		{Version: 1, Up: func(*sql.Tx) error { return nil }, Down: func(*sql.Tx) error { return nil }},
	}
	dbVersion := 0

	mock.ExpectBegin()
	mock.ExpectCommit()

	var insertVersions []int
	insertFn := func(tx *sql.Tx, version int) error {
		insertVersions = append(insertVersions, version)
		return nil
	}

	err = Up(db, insertFn, dbVersion, state)
	if err != nil {
		t.Fatalf("Up: %v", err)
	}

	if len(insertVersions) != 1 || insertVersions[0] != 1 {
		t.Errorf("insertFn called with %v, want [1]", insertVersions)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock: %v", err)
	}
}

func TestUp_AppliesMultipleMigrationsInOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	state := []migration.Migrate{
		{Version: 1, Up: nil, Down: nil},
		{Version: 3, Up: func(*sql.Tx) error { return nil }, Down: nil},
		{Version: 2, Up: func(*sql.Tx) error { return nil }, Down: nil},
	}
	dbVersion := 1

	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectCommit()

	var insertVersions []int
	insertFn := func(tx *sql.Tx, version int) error {
		insertVersions = append(insertVersions, version)
		return nil
	}

	err = Up(db, insertFn, dbVersion, state)
	if err != nil {
		t.Fatalf("Up: %v", err)
	}

	want := []int{2, 3}
	if len(insertVersions) != len(want) {
		t.Fatalf("insertFn called %d times, want %d", len(insertVersions), len(want))
	}
	for i := range want {
		if insertVersions[i] != want[i] {
			t.Errorf("insertVersions[%d] = %d, want %d", i, insertVersions[i], want[i])
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock: %v", err)
	}
}

func TestUp_DuplicateVersionsReturnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	state := []migration.Migrate{
		{Version: 1, Up: func(*sql.Tx) error { return nil }, Down: func(*sql.Tx) error { return nil }},
		{Version: 1, Up: func(*sql.Tx) error { return nil }, Down: func(*sql.Tx) error { return nil }},
	}
	dbVersion := 0

	insertFn := func(tx *sql.Tx, version int) error { return nil }

	err = Up(db, insertFn, dbVersion, state)
	if err == nil {
		t.Fatal("Up should return error when state has duplicate versions")
	}
	if err.Error() != "migration 1 is repeated" {
		t.Errorf("Up() error = %q, want %q", err.Error(), "migration 1 is repeated")
	}

	// No DB operations should be performed when duplicate versions are detected
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock: %v", err)
	}
}

func TestUp_InsertFnErrorPanics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	state := []migration.Migrate{
		{Version: 1, Up: func(*sql.Tx) error { return nil }, Down: nil},
	}
	dbVersion := 0

	mock.ExpectBegin()
	mock.ExpectRollback()

	insertFn := func(tx *sql.Tx, version int) error {
		return sql.ErrTxDone
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Up should panic when insertFn returns error")
		}
	}()
	_ = Up(db, insertFn, dbVersion, state)
}
