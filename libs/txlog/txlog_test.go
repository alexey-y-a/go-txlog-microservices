package txlog

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestFileLog_Append(t *testing.T) {
    t.Helper()
    tempDir := t.TempDir()
    logPath := tempDir + "/test.log"

    logFile, err := NewFileLog(logPath)
    require.NoError(t, err, "NewFileLog should not return error")

    defer func() {
        err := logFile.Close()
        require.NoError(t, err, "Close should not return error")
    }()

    event1 := Event {
        Key: "user1",
        Value: "Alice",
        Op: "set",
    }

    err = logFile.Append(event1)
    require.NoError(t, err, "Append for event1 should not return error")

    event2 := Event {
        Key: "user2",
        Value: "",
        Op: "delete",
    }

    err = logFile.Append(event2)
    require.NoError(t, err, "Append for event2 should not return error")

    data, err := os.ReadFile(logPath)
    require.NoError(t, err, "ReadFile should not return error")

    content := string(data)

    require.Contains(t, content, "set 5 5 user1Alice", "log should contain encoded first event")
    require.Contains(t, content, "delete 5 0 user2", "log should contain encoded second event")
}

func TestFileLog_AppendTooLarge(t *testing.T) {
    t.Helper()

    tempDir := t.TempDir()
    logPath := tempDir + "/test.log"

    logFile, err := NewFileLog(logPath)
    require.NoError(t, err)
    defer func() {
        err := logFile.Close()
        require.NoError(t, err)
    }()

    tooLongKey := strings.Repeat("a", MaxKeySize+1)

    err = logFile.Append(Event{
        Key: tooLongKey,
        Value: "x",
        Op: "set",
    })
    require.ErrorIs(t, err, ErrKeyTooLarge)

    tooLongValue := strings.Repeat("b", MaxValueSize+1)

    err = logFile.Append(Event{
        Key: "key",
        Value: tooLongValue,
        Op: "set",
    })
    require.ErrorIs(t, err, ErrValueTooLarge)
}