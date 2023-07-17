package internal

import (
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("valid command, suppressed", func(t *testing.T) {
		c := command{cmdName: "echo", cmdArgs: []string{"action"}, suppressed: true}
		err := c.execute()
		assertErr(t, err, nil)
	})

	t.Run("valid command, not suppressed", func(t *testing.T) {
		c := command{cmdName: "echo", cmdArgs: []string{"action"}, suppressed: false}
		err := c.execute()
		assertErr(t, err, nil)
	})

	t.Run("invalid command", func(t *testing.T) {
		c := command{cmdName: "invalid", cmdArgs: []string{"action"}, suppressed: true}
		err := c.execute()
		assertErr(t, err, ErrCouldntExecuteCommand)
	})
}
