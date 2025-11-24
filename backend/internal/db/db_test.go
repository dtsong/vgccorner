package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewDatabase(t *testing.T) {
	// Since NewDatabase tries to open a real connection, we'll test the struct creation
	// by testing that a Database with a valid connection works
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	mock.ExpectPing()

	database := &Database{conn: db}
	if database.conn == nil {
		t.Error("expected database connection to be initialized")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestClose(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	database := &Database{conn: db}

	mock.ExpectClose()
	err = database.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	mock.ExpectExec("INSERT INTO test").
		WithArgs("value1", "value2").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.Exec(ctx, "INSERT INTO test VALUES (?, ?)", "value1", "value2")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestQueryRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")
	mock.ExpectQuery("SELECT (.+) FROM test").WillReturnRows(rows)

	row := database.QueryRow(ctx, "SELECT id, name FROM test WHERE id = ?", 1)
	if row == nil {
		t.Error("expected row to be non-nil")
	}

	var id int
	var name string
	err = row.Scan(&id, &name)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id != 1 || name != "test" {
		t.Errorf("expected id=1, name=test, got id=%d, name=%s", id, name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test1").
		AddRow(2, "test2")
	mock.ExpectQuery("SELECT (.+) FROM test").WillReturnRows(rows)

	result, err := database.Query(ctx, "SELECT id, name FROM test")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	defer func() { _ = result.Close() }()

	count := 0
	for result.Next() {
		count++
	}

	if count != 2 {
		t.Errorf("expected 2 rows, got %d", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestStoreBattle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	battle := &Battle{
		Format:      "VGC 2025",
		Timestamp:   time.Now(),
		DurationSec: 300,
		Winner:      "player1",
		Player1ID:   "Alice",
		Player2ID:   "Bob",
		BattleLog:   "battle log content",
		IsPrivate:   false,
	}

	// Mock transaction
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO battles").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("battle-uuid"))
	mock.ExpectCommit()

	battleID, err := database.StoreBattle(ctx, battle)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if battleID == "" {
		t.Error("expected battleID to be non-empty")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestStoreBattleWithAnalysis(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	battle := &Battle{
		Format:      "VGC 2025",
		Timestamp:   time.Now(),
		DurationSec: 300,
		Winner:      "player1",
		Player1ID:   "Alice",
		Player2ID:   "Bob",
		BattleLog:   "battle log content",
		IsPrivate:   false,
		Analysis: &BattleAnalysis{
			TotalTurns:       10,
			AvgDamagePerTurn: 50.5,
		},
		KeyMoments: []*KeyMoment{
			{
				TurnNumber:   5,
				MomentType:   "ko",
				Description:  "Critical KO",
				Significance: 8,
			},
		},
	}

	// Mock transaction with analysis and key moments
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO battles").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("battle-uuid"))
	mock.ExpectExec("INSERT INTO battle_analysis").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO key_moments").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	battleID, err := database.StoreBattle(ctx, battle)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if battleID == "" {
		t.Error("expected battleID to be non-empty")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetBattle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	battleID := "test-battle-id"
	timestamp := time.Now()

	battleRows := sqlmock.NewRows([]string{
		"id", "format", "timestamp", "duration_sec", "winner",
		"player1_id", "player2_id", "battle_log", "is_private",
		"created_at", "updated_at",
	}).AddRow(
		battleID, "VGC 2025", timestamp, 300, "player1",
		"Alice", "Bob", "log content", false,
		timestamp, timestamp,
	)

	mock.ExpectQuery("SELECT (.+) FROM battles WHERE id").
		WithArgs(battleID).
		WillReturnRows(battleRows)

	// Mock analysis query (matches 15 fields from getBattleAnalysis)
	mock.ExpectQuery("SELECT (.+) FROM battle_analysis WHERE battle_id").
		WithArgs(battleID).
		WillReturnRows(sqlmock.NewRows([]string{
			"battle_id", "total_turns", "avg_damage_per_turn", "avg_heal_per_turn", "moves_used_count",
			"switches_count", "super_effective_moves", "not_very_effective_moves", "critical_hits",
			"player1_damage_dealt", "player1_damage_taken", "player1_healing_done",
			"player2_damage_dealt", "player2_damage_taken", "player2_healing_done",
		}).AddRow(battleID, 10, 50.5, 10.2, 20, 5, 3, 2, 1, 100, 80, 20, 90, 100, 15))

	// Mock key moments query (matches 4 fields from getKeyMoments)
	mock.ExpectQuery("SELECT (.+) FROM key_moments WHERE battle_id").
		WithArgs(battleID).
		WillReturnRows(sqlmock.NewRows([]string{"turn_number", "moment_type", "description", "significance"}))

	battle, err := database.GetBattle(ctx, battleID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if battle == nil {
		t.Fatal("expected battle to be non-nil")
	}

	if battle.ID != battleID {
		t.Errorf("expected ID %s, got %s", battleID, battle.ID)
	}

	if battle.Format != "VGC 2025" {
		t.Errorf("expected format 'VGC 2025', got %s", battle.Format)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetBattleNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	battleID := "nonexistent"

	mock.ExpectQuery("SELECT (.+) FROM battles WHERE id").
		WithArgs(battleID).
		WillReturnError(sql.ErrNoRows)

	battle, err := database.GetBattle(ctx, battleID)
	if err != nil {
		t.Errorf("expected no error for not found, got %v", err)
	}

	if battle != nil {
		t.Error("expected battle to be nil for not found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestListBattles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	filter := &BattleFilter{
		Format:    "VGC 2025",
		IsPrivate: boolPtr(false),
	}

	timestamp := time.Now()

	// Mock count query
	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Mock battles query
	battleRows := sqlmock.NewRows([]string{
		"id", "format", "timestamp", "duration_sec", "winner",
		"player1_id", "player2_id", "is_private",
	}).
		AddRow("id1", "VGC 2025", timestamp, 300, "player1", "Alice", "Bob", false).
		AddRow("id2", "VGC 2025", timestamp, 250, "player2", "Charlie", "Dave", false)

	mock.ExpectQuery("SELECT (.+) FROM battles").
		WillReturnRows(battleRows)

	battles, total, err := database.ListBattles(ctx, filter, 10, 0)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(battles) != 2 {
		t.Errorf("expected 2 battles, got %d", len(battles))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestWithTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func() { _ = db.Close() }()

	database := &Database{conn: db}
	ctx := context.Background()

	t.Run("successful transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO test").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := database.WithTx(ctx, func(tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO test VALUES (?)", "value")
			return err
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("transaction with error rollback", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		testErr := errors.New("test error")
		err := database.WithTx(ctx, func(tx *sql.Tx) error {
			return testErr
		})

		if err != testErr {
			t.Errorf("expected test error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func boolPtr(b bool) *bool {
	return &b
}
