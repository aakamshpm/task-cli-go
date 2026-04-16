package jsonstore

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aakamshpm/task-cli-go/internal/task"
)

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	baseDir := t.TempDir()
	filePath := filepath.Join(baseDir, "data", "tasks.json")
	store := NewStore(filePath)
	got, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty tasks, got %d", len(got))
	}
}
func TestSaveAndLoad_RoundTrip(t *testing.T) {
	baseDir := t.TempDir()
	filePath := filepath.Join(baseDir, "data", "tasks.json")
	store := NewStore(filePath)
	fixed := time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC)
	in := []task.Task{
		{ID: 1, Title: "Task A", Done: false, CreatedAt: fixed},
		{ID: 2, Title: "Task B", Done: true, CreatedAt: fixed.Add(time.Minute)},
	}
	if err := store.Save(in); err != nil {
		t.Fatalf("expected no error on save, got %v", err)
	}
	out, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error on load, got %v", err)
	}
	if len(out) != len(in) {
		t.Fatalf("expected %d tasks, got %d", len(in), len(out))
	}
	for i := range in {
		if out[i].ID != in[i].ID ||
			out[i].Title != in[i].Title ||
			out[i].Done != in[i].Done ||
			!out[i].CreatedAt.Equal(in[i].CreatedAt) {
			t.Fatalf("task mismatch at index %d: got %+v, want %+v", i, out[i], in[i])
		}
	}
}
func TestSave_CreatesParentDirectory(t *testing.T) {
	baseDir := t.TempDir()
	filePath := filepath.Join(baseDir, "nested", "dir", "tasks.json")
	store := NewStore(filePath)
	if err := store.Save([]task.Task{}); err != nil {
		t.Fatalf("expected no error on save, got %v", err)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected tasks file to exist, got stat error: %v", err)
	}
}
func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	baseDir := t.TempDir()
	filePath := filepath.Join(baseDir, "data", "tasks.json")
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filePath, []byte("{invalid json"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	store := NewStore(filePath)
	_, err := store.Load()
	if err == nil {
		t.Fatalf("expected error for invalid json, got nil")
	}
}
func TestSave_Overwrite_ReplacesContent(t *testing.T) {
	baseDir := t.TempDir()
	filePath := filepath.Join(baseDir, "data", "tasks.json")
	store := NewStore(filePath)
	first := []task.Task{
		{ID: 1, Title: "Old Task", Done: false, CreatedAt: time.Date(2026, 4, 16, 9, 0, 0, 0, time.UTC)},
	}
	second := []task.Task{
		{ID: 2, Title: "New Task", Done: true, CreatedAt: time.Date(2026, 4, 16, 9, 30, 0, 0, time.UTC)},
	}
	if err := store.Save(first); err != nil {
		t.Fatalf("first save failed: %v", err)
	}
	if err := store.Save(second); err != nil {
		t.Fatalf("second save failed: %v", err)
	}
	out, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 task after overwrite, got %d", len(out))
	}
	if out[0].ID != 2 || out[0].Title != "New Task" || !out[0].Done {
		t.Fatalf("unexpected task after overwrite: %+v", out[0])
	}
}
