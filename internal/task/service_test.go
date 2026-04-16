package task

import (
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	tasks    []Task
	loadErr  error
	saveErr  error
	saveCall int
}

func (f *fakeStore) Load() ([]Task, error) {
	if f.loadErr != nil {
		return nil, f.loadErr
	}
	out := make([]Task, len(f.tasks))
	copy(out, f.tasks)
	return out, nil
}
func (f *fakeStore) Save(tasks []Task) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saveCall++
	f.tasks = make([]Task, len(tasks))
	copy(f.tasks, tasks)
	return nil
}
func TestAdd_Success(t *testing.T) {
	fixed := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{}
	svc := NewService(store, func() time.Time { return fixed })
	got, err := svc.Add("Learn Go")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("expected ID 1, got %d", got.ID)
	}
	if got.Title != "Learn Go" {
		t.Fatalf("expected title Learn Go, got %q", got.Title)
	}
	if got.Done {
		t.Fatalf("expected Done=false")
	}
	if !got.CreatedAt.Equal(fixed) {
		t.Fatalf("expected CreatedAt %v, got %v", fixed, got.CreatedAt)
	}
	if len(store.tasks) != 1 {
		t.Fatalf("expected store to have 1 task, got %d", len(store.tasks))
	}
}
func TestAdd_EmptyTitle(t *testing.T) {
	store := &fakeStore{}
	svc := NewService(store, nil)
	_, err := svc.Add("   ")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
	if store.saveCall != 0 {
		t.Fatalf("expected no save calls, got %d", store.saveCall)
	}
}
func TestList_ReturnsCopy(t *testing.T) {
	store := &fakeStore{
		tasks: []Task{
			{ID: 1, Title: "Task A", Done: false, CreatedAt: time.Now().UTC()},
		},
	}
	svc := NewService(store, nil)
	list1, err := svc.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	list1[0].Title = "Mutated"
	list2, err := svc.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if list2[0].Title != "Task A" {
		t.Fatalf("expected internal state unchanged, got %q", list2[0].Title)
	}
}

func TestMarkDone_Success(t *testing.T) {
	store := &fakeStore{
		tasks: []Task{
			{ID: 1, Title: "Task A", Done: false, CreatedAt: time.Now().UTC()},
		},
	}
	svc := NewService(store, nil)
	got, err := svc.MarkDone(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !got.Done {
		t.Fatalf("expected task to be marked done")
	}
	if !store.tasks[0].Done {
		t.Fatalf("expected persisted task to be marked done")
	}
}
func TestMarkDone_NotFound(t *testing.T) {
	store := &fakeStore{}
	svc := NewService(store, nil)
	_, err := svc.MarkDone(99)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}
func TestDelete_Success(t *testing.T) {
	store := &fakeStore{
		tasks: []Task{
			{ID: 1, Title: "Task A", CreatedAt: time.Now().UTC()},
			{ID: 2, Title: "Task B", CreatedAt: time.Now().UTC()},
		},
	}
	svc := NewService(store, nil)
	if err := svc.Delete(1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store.tasks) != 1 {
		t.Fatalf("expected 1 task left, got %d", len(store.tasks))
	}
	if store.tasks[0].ID != 2 {
		t.Fatalf("expected remaining task ID=2, got %d", store.tasks[0].ID)
	}
}
func TestDelete_NotFound(t *testing.T) {
	store := &fakeStore{}
	svc := NewService(store, nil)
	err := svc.Delete(99)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}
