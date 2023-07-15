package internal

import "testing"

func TestExecute(t *testing.T) {
	t.Run("valid command with @", func(t *testing.T) {
		gomake := NewGomake()
		err := gomake.setCommand(target{name: "target"}, command{cmd: "echo 'action'", suppressed: true})
		assertErr(t, err, nil)

		err = gomake.targets["target"].commands[0].execute()
		assertErr(t, err, nil)
	})

	t.Run("valid command without @", func(t *testing.T) {
		gomake := NewGomake()
		err := gomake.setCommand(target{name: "target"}, command{cmd: "echo 'action'", suppressed: false})

		assertErr(t, err, nil)

		err = gomake.targets["target"].commands[0].execute()
		assertErr(t, err, nil)
	})

	t.Run("invalid commands", func(t *testing.T) {
		gomake := NewGomake()
		err := gomake.setCommand(target{name: "target"}, command{cmd: "invalid", suppressed: true})

		assertErr(t, err, nil)

		err = gomake.targets["target"].commands[0].execute()
		assertErr(t, err, ErrCouldntExecuteCommand)
	})

	t.Run("dependency not exist", func(t *testing.T) {
		gomake := NewGomake()

		err := gomake.Run("target")
		assertErr(t, err, ErrDependencyNotFound)
	})
}
