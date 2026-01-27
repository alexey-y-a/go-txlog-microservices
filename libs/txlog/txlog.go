package txlog

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

const (
    MaxKeySize = 1024
    MaxValueSize = 65536
)

var (
    ErrKeyTooLarge = errors.New("txlog: key size exceeds MaxKeySize")
    ErrValueTooLarge = errors.New("txlog: value size exceeds MaxValueSize")
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
    keyBytes := []byte(e.Key)
    valBytes := []byte(e.Value)

    if len(keyBytes) > MaxKeySize {
        return ErrKeyTooLarge
    }

    if len(valBytes) > MaxValueSize {
        return ErrValueTooLarge
    }

    prefix := fmt.Sprintf("%s %d %d ", e.Op, len(keyBytes), len(valBytes))

    var buf bytes.Buffer

    _, err := buf.WriteString(prefix)
    if err != nil {
        return fmt.Errorf("txlog: write prefix: %w", err)
    }

    _, err = buf.Write(keyBytes)
    if err != nil {
        return fmt.Errorf("txlog: write key: %w", err)
    }

    _, err = buf.Write(valBytes)
    if err != nil {
        return fmt.Errorf("txlog: write value: %w", err)
    }

    err = buf.WriteByte('\n')
    if err != nil {
        return fmt.Errorf("txlog: write newline: %w", err)
    }

    _, err = l.file.Write(buf.Bytes())
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