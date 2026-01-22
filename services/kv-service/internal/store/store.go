package store

import (
	"fmt"
	"sync"

	"github.com/alexey-y-a/go-txlog-microservices/libs/txlog"
)

type Store struct{
    mu sync.RWMutex
    data map[string]string
    log txlog.Log
}

func NewStore(log txlog.Log) *Store {
    return &Store {
        data: make(map[string]string),
        log: log,
    }
}

func (s *Store) Set(key, value string) error {
    event := txlog.Event {
        Key: key,
        Value: value,
        Op: "set",
    }

    err := s.log.Append(event)
    if err != nil {
        return fmt.Errorf("store: append set event: %w", err)
    }

    s.mu.Lock()
    s.data[key] = value
    s.mu.Unlock()

    return nil
}

func (s *Store) Get(key string) (string, bool) {
    s.mu.RLock()
    value, ok :=  s.data[key]
    s.mu.RUnlock()

    return value, ok
}

func (s *Store) Delete(key string) error {
    event := txlog.Event {
        Key: key,
        Value: "",
        Op: "delete",
    }

    err := s.log.Append(event)
    if err != nil {
        return fmt.Errorf("store: append delete event: %w", err)
    }

    s.mu.Lock()
    delete(s.data, key)
    s.mu.Unlock()

    return nil
}

