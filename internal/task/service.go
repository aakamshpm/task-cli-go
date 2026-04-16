package task

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var ErrEmptyTitle = errors.New("title cannot be empty")
var ErrTaskNotFound = errors.New("task not found")

// custom time type named 'Clock'
// this is handy in test cases where we would create a fixed time and pass it to instance of service during task creation or any other operation
// we are using a type function becaues the field 'now' in 'Service' struct uses a function call to get current time.
// so in real usage, we pass real current time. but in test, we create a fixed time and use this function to pass it.
// in short, this is just a line of code to eliminate writing 'func() time.Time' everytime in function params
type Clock func() time.Time

type Service struct {
	mu    sync.Mutex // since we are handling with json files, read/write operations will be slower, hence more strict. so we only allow one person at a time in the instance using Mutex
	store Store
	now   Clock
}

// in GO, we don't have class, new keyword, and typical constructor creation and usage with "this" like other languages
// hence we need to create a custom constuctor function starting with 'New' keyword for convention followed by the type name. Eg: NewService
func NewService(store Store, now Clock) *Service {
	if now == nil {
		now = time.Now // we assign the Now() funtion to the field 'now', this will be called each time on task creation. NOTE: Service.now is a function
	}

	return &Service{
		store: store,
		now:   now,
	}
}

func nextID(tasks []Task) int {
	maxID := 0
	// range in Go returns two things: index and the value at that index
	// since we dont care about index this time, we use '_' to cover it. '_' means we dont use it, hence throw it away
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

func indexByID(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}

	}
	return -1
}

// *Service ensures that Add function is under type Service and usage outside
// the 's' here acts as a 'this' instance
func (s *Service) Add(title string) (Task, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return Task{}, ErrEmptyTitle
	}

	// solves reader-writer problem
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.store.Load()
	if err != nil {
		return Task{}, fmt.Errorf("load tasks: %w", err)
	}

	t := Task{
		ID:        nextID(tasks),
		Title:     title,
		Done:      false,
		CreatedAt: s.now().UTC(),
	}

	tasks = append(tasks, t)

	if err := s.store.Save(tasks); err != nil {
		return Task{}, fmt.Errorf("save tasks: %w", err)
	}

	return t, nil
}

func (s *Service) MarkDone(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.store.Load()
	if err != nil {
		return Task{}, fmt.Errorf("load tasks: %w", err)
	}

	index := indexByID(tasks, id)
	if index == -1 {
		return Task{}, fmt.Errorf("%w: %d", ErrTaskNotFound, id)
	}

	tasks[index].Done = true

	if err := s.store.Save(tasks); err != nil {
		return Task{}, fmt.Errorf("save tasks: %w", err)
	}

	return tasks[index], nil
}

func (s *Service) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.store.Load()
	if err != nil {
		return fmt.Errorf("load tasks: %w", err)
	}

	index := indexByID(tasks, id)
	if index == -1 {
		return fmt.Errorf("%w: %d", ErrTaskNotFound, id)
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	if err := s.store.Save(tasks); err != nil {
		return fmt.Errorf("save tasks: %w", err)
	}

	return nil
}

func (s *Service) List() ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.store.Load()
	if err != nil {
		return nil, fmt.Errorf("load tasks: %w", err)
	}

	out := make([]Task, len(tasks))
	copy(out, tasks)
	return out, nil
}
