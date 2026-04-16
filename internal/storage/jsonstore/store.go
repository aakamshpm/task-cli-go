package jsonstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aakamshpm/task-cli-go/internal/task"
)

type Store struct {
	filePath string
}

func NewStore(filepath string) *Store {
	return &Store{filePath: filepath}
}

func (s *Store) Load() ([]task.Task, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		// when we open the file for first time, there wont be any file; hence we handle that gracefully with ErrNotExist
		if errors.Is(err, os.ErrNotExist) {
			return []task.Task{}, nil
		}

		return nil, fmt.Errorf("open file: %w", err)
	}

	// defer should come after checking for error returned from file opening
	defer file.Close()

	var tasks []task.Task

	// creates a json decoder that reads from the file stream
	decoder := json.NewDecoder(file)

	// coverts the json into tasks format and store it in the slice decalared above
	if err := decoder.Decode(&tasks); err != nil {
		// if the decoded tasks have no tasks in it, the error would be an EndOfFile(EOF)
		if errors.Is(err, io.EOF) {
			return []task.Task{}, nil
		}
		return nil, fmt.Errorf("decode json: %w", err)
	}

	if tasks == nil {
		return []task.Task{}, nil
	}

	return tasks, nil
}

func (s *Store) Save(tasks []task.Task) error {
	// we extract the folder name using 'filepath.Dir()' from filePath and create it if it doesn't exist with proper permissions
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	// saving a task is an atomic process; we dont directly edit the json filePath
	// instead we create a .tmp file and copy the contents from original json (if it exists) to it along with what we intend to save
	// if the write was a success to tmp file, we rename the tmp file to original file name to make the new data available instant and atomic
	// even if we loose power in between writes, the original json stays untouched
	tmpPath := s.filePath + ".tmp"
	file, err := os.Create(tmpPath)
	// defer make sures the file is closed after every operation even if error occurs
	defer func() { _ = file.Close() }()

	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	// if anything goes wrong, we remove the tmp file at last
	defer func() { _ = os.Remove(tmpPath) }()

	// creates a json encoder to perform write operations
	encoder := json.NewEncoder(file)
	// we format the json with indentation (pretty format) so that it doesnt get stored in one line
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	// 'file.Sync()' forces the OS to flush everything from memory buffers to actual disk
	// without this, data may still reside in RAM and not get written to disk
	// if power cuts down before sync, data will be lost
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync temp file: %w", err)
	}
	// we need to close the tmp file before reading, because thats the standard practice when doing rename operation
	// in OS like Windows, it wont allow rename if the file is open
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, s.filePath); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
