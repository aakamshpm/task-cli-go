package task

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrEmptyTitle = errors.New("title cannot be empty")

// custom time type named 'Clock'
// this is handy in test cases where we would create a fixed time and pass it to instance of service during task creation or any other operation
// we are using a type function becaues the field 'now' in 'Service' struct uses a function call to get current time.
// so in real usage, we pass real current time. but in test, we create a fixed time and use this function to pass it.
// in short, this is just a line of code to eliminate writing 'func() time.Time' everytime in function params
type Clock func() time.Time

type Service struct {
	mu      sync.RWMutex
	tasks   []Task
	nextInt int
	now     Clock
}

// in GO, we don't have class, new keyword, and typical constructor creation and usage with "this" like other languages
// hence we need to create a custom constuctor function starting with 'New' keyword for convention followed by the type name. Eg: NewService
func NewService(now Clock) *Service {
	if now == nil {
		now = time.Now // we assign the Now() funtion to the field 'now', this will be called each time on task creation. NOTE: Service.now is a function
	}

	return &Service{
		tasks:   make([]Task, 0), // we intentionaly create a slice thats empty and ready to grow
		nextInt: 1,
		now:     now,
	}
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

	t := Task{
		ID:        s.nextInt,
		Title:     title,
		Done:      false,
		CreatedAt: s.now().UTC(),
	}

	s.tasks = append(s.tasks, t)
	s.nextInt++

	return t, nil
}

func (s *Service) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}
