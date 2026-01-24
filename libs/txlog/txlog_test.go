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

    require.True(t, strings.Contains(content, "set user1 Alice\n"), "log should contain first event line")
    require.True(t, strings.Contains(content, "delete user2 \n"), "log should contain second event line")
}
