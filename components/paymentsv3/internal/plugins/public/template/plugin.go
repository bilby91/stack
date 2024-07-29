package template

type Plugin struct {
	name string
}

func (p *Plugin) Name() string {
	return p.name
}

// var _ models.Plugin = &Plugin{}
