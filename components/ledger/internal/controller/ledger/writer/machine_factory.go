package writer

import "github.com/formancehq/ledger/internal/machine/vm"

//go:generate mockgen -source machine_factory.go -destination machine_factory_generated.go -package writer . MachineFactory
type MachineFactory interface {
	Make(script string) (Machine, error)
}

type DefaultMachineFactory struct {
	store    vm.Store
	compiler Compiler
}

func (d *DefaultMachineFactory) Make(script string) (Machine, error) {
	ret, err := d.compiler.Compile(script)
	if err != nil {
		return nil, err
	}
	return NewDefaultMachine(*ret, d.store), nil
}

func NewDefaultMachineFactory(compiler Compiler, store vm.Store) *DefaultMachineFactory {
	return &DefaultMachineFactory{
		compiler: compiler,
		store:    store,
	}
}

var _ MachineFactory = (*DefaultMachineFactory)(nil)
