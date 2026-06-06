package db

import (
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"testing"
	"time"
)

func TestLoadConfig_NoConfigReturnsDefault(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM config`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	config, err := LoadConfig(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Theme != "sunset" {
		t.Errorf("expected theme sunset, got %s", config.Theme)
	}
}

func TestLoadConfig_ValidConfig(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM config`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO config (id, config)
		VALUES (1, '{"theme":"dracula"}')
	`)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	config, err := LoadConfig(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Theme != "dracula" {
		t.Errorf("expected dracula, got %s", config.Theme)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM config`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO config (id, config)
		VALUES (1, '{invalid json}')
	`)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	_, err = LoadConfig(db)

	if err == nil {
		t.Errorf("expected JSON unmarshal error, got nil")
	}
}

func TestFlush_SkipsInvalidJSON(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM offline_logs`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	// invalid JSON
	_, err = db.Exec(
		`INSERT INTO offline_logs (data) VALUES (?)`,
		`{invalid json}`,
	)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	// valid JSON
	_, err = db.Exec(
		`INSERT INTO offline_logs (data) VALUES (?)`,
		`{"name":"main.go","project":"test","duration":60}`,
	)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	uploaded := false

	result := Flush(func(batch []types.LogEntry) error {
		uploaded = true

		if len(batch) != 1 {
			t.Errorf("expected 1 valid entry, got %d", len(batch))
		}

		return nil
	}, db)

	if !uploaded {
		t.Errorf("expected upload function to be called")
	}

	if !result {
		t.Errorf("expected Flush to return true")
	}

	var count int

	err = db.QueryRow(
		`SELECT COUNT(*) FROM offline_logs`,
	).Scan(&count)

	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 row remaining, got %d", count)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM config`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	expected := types.CacheData{
		Theme: "dracula",
	}

	err = SaveConfig(db, expected)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := LoadConfig(db)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if got.Theme != expected.Theme {
		t.Errorf("expected theme %s, got %s", expected.Theme, got.Theme)
	}
}

func TestSaveAndGetDashboardCache(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM dashboard_cache`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	expected := types.CacheData{
		Theme: "dracula",
	}

	err = SaveDashboardCache(db, expected)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := GetDashboardCache(db)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if got.Theme != expected.Theme {
		t.Errorf("expected theme %s, got %s", expected.Theme, got.Theme)
	}
}

func TestGetDashboardCache_Expired(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM dashboard_cache`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	_, err = db.Exec(
		`INSERT INTO dashboard_cache (id, data, updated_at) VALUES ('main', '{}', ?)`,
		time.Now().Add(-10*time.Minute),
	)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	_, err = GetDashboardCache(db)

	if err == nil {
		t.Errorf("expected cache expired error")
	}
}

func TestEnqueue(t *testing.T) {
	db, err := InitSQLite()
	if err != nil {
		t.Fatalf("failed to init sqlite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM offline_logs`)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	entry := types.LogEntry{
		FileName: "main.go",
		Project:  "test-project",
		Duration: 60,
	}

	err = Enqueue(entry, db)
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM offline_logs`).Scan(&count)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 row, got %d", count)
	}
}
