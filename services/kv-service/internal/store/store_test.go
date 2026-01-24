package store

import (
	"testing"

	"github.com/alexey-y-a/go-txlog-microservices/libs/txlog"
	"github.com/stretchr/testify/require"
)


type fakeLog struct {
    events []txlog.Event
}

func (f *fakeLog) Append(e txlog.Event) error {
    f.events = append(f.events, e)
    return nil
}

func TestStore_SetAndGet(t *testing.T) {
    t.Helper()

    flog := &fakeLog{}
    s := NewStore(flog)

    err := s.Set("user42", "Alice")
    require.NoError(t, err, "Set should not return error")

    value, ok := s.Get("user42")
    require.True(t, ok, "Get should return correct value")
    require.Equal(t, "Alice", value, "Get should return correct value")

    require.Len(t, flog.events, 1, "fakeLog should contain one event")
    require.Equal(t, "set", flog.events[0].Op, "event Op should be 'set'")
    require.Equal(t, "user42", flog.events[0].Key, "event Key should match")
    require.Equal(t, "Alice", flog.events[0].Value, "event Value should match")
}

func TestStore_Delete(t *testing.T) {
    t.Helper()

    flog := &fakeLog{}
    s := NewStore(flog)

    err := s.Set("user1", "Bob")
    require.NoError(t, err, "Set should not return error")

    value, ok := s.Get("user1")
    require.True(t, ok, "Get should report that key exist before delete")
    require.Equal(t, "Bob", value, "Get should return correct value before delete")

    err = s.Delete("user1")
    require.NoError(t, err, "Delete should not return error")

    _, ok = s.Get("user1")
    require.False(t, ok, "Get should report that key does not exist after delete")

    require.Len(t, flog.events, 2, "fakeLog should contain two events")

    deleteEvent := flog.events[1]
    require.Equal(t, "delete", deleteEvent.Op, "second event Op should be 'delete'")
    require.Equal(t, "user1", deleteEvent.Key, "second event Key should match")
}