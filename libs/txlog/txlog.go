package txlog

import (
	"fmt"
	"os"
)

type Event struct {
    Key   string
    Value string
    Op    string
}


type Log interface {
    Append(e Event) error
    Sync() error
    Close() error
}

type FileLog struct {
    file *os.File
}

func NewFileLog(path string) (*FileLog, error) {
    file, err := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0o644)
    if err != nil {
        return nil, fmt.Errorf("txlog: open file %q: %w", path, err)
    }

    log :=  &FileLog {
        file: file,
    }

    return log, nil
}

func (l *FileLog) Append(e Event) error {
    line := fmt.Sprintf("%s %s %s\n", e.Op, e.Key, e.Value)

    _, err := l.file.Write([]byte(line))
    if err != nil {
        return fmt.Errorf("txlog: append event: %w", err)
    }

    return nil
}

func (l *FileLog) Sync() error {
    err := l.file.Sync()
    if err != nil {
        return fmt.Errorf("txlog: sync file: %w", err)
    }
    return nil
}

func (l *FileLog) Close() error {
    err := l.Sync()
    if err != nil {
        return err
    }

    err = l.file.Close()
    if err != nil {
        return fmt.Errorf("txlog: close file: %w", err)
    }
    return nil
}