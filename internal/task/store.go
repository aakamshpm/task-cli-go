package task

// an interface is like a contract.
// here the Store only cares about two things: you must be able to Load() and return tasks; you must be able to Save() tasks
// anything that meets these two conditions can be a store: JSON file, database, s3....
type Store interface {
	Load() ([]Task, error)
	Save([]Task) error
}
