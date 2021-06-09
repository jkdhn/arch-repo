package pkg

import "arch-repo/pkg/desc"

type Package struct {
	description *desc.Description
}

func (p *Package) Description() *desc.Description {
	return p.description
}
