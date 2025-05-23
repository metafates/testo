package myplugin

import (
	"testman"
)

var _ testman.Plugin = (*Plugin)(nil)

type Plugin interface {
	Label(name string)
}

type myPlugin struct{}

func New() Plugin {
	return myPlugin{}
}

func (myPlugin) Label(name string) {
	println("label: " + name)
}
