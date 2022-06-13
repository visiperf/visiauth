package visiauth

type Machine struct {
	id string
}

func NewMachine(id string) *Machine {
	return &Machine{id}
}

func (m Machine) ID() string {
	return m.id
}
