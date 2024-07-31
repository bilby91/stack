package worker

type Definition struct {
	Func any
	Name string
}

type DefinitionSet []Definition

func NewDefinitionSet() DefinitionSet {
	return DefinitionSet{}
}

func (d DefinitionSet) Append(definition Definition) DefinitionSet {
	d = append(d, definition)

	return d
}
