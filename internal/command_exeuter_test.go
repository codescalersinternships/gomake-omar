package internal

import "testing"

func TestExecute(t *testing.T) {
	t.Run("valid commands", func(t *testing.T) {
		gomake := NewGomake()
		err := gomake.addActionLine("target", "echo 'action'")
		assertErr(t, err, nil)
		err = gomake.addActionLine("target", "@echo 'action'")
		assertErr(t, err, nil)

		err = gomake.executer.execute([]target{"target", "target"})
		assertErr(t, err, nil)
	})
	t.Run("invalid commands", func(t *testing.T) {
		gomake := NewGomake()
		err := gomake.addActionLine("target", "echo_o")
		assertErr(t, err, nil)
		err = gomake.addActionLine("target", "@echo -f 'action'")
		assertErr(t, err, nil)

		err = gomake.executer.execute([]target{"target", "target"})
		assertErr(t, err, ErrCouldntExecuteCommand)
	})
	t.Run("commands not exist", func(t *testing.T) {
		gomake := NewGomake()

		err := gomake.executer.execute([]target{"target"})
		assertErr(t, err, ErrDependencyNotFound)
	})
}
