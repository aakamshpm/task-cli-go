package task

import (
	"errors"
	"testing"
	"time"
)

func TestAdd_Success(t *testing.T) {
	fixed := time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)
	svc := NewService(func() time.Time { return fixed })

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
}

func TestAdd_EmptyTitle(t *testing.T) {
	svc := NewService(nil)

	_, err := svc.Add(" ")

	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestList_ReturnsCopy(t *testing.T) {
	svc := NewService(nil)
	_, _ = svc.Add("Task A")

	list1 := svc.List()
	list1[0].Title = "Mutated"

	list2 := svc.List()

	if list2[0].Title != "Task A" {
		t.Fatalf("expected internal state unchanged, got %q", list2[0].Title)
	}
}
