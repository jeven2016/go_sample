package _struct

import (
	"reflect"
	"testing"
)

type NameInterface interface {
	GetName() string
}

//===============parent================
type Parent struct {
	Name string
	Desc string
}

func (p *Parent) GetName() string {
	return p.Name
}

//===============Child================
type Child struct {
	Parent
	Type string
}

func (c *Child) GetName() string {
	return c.Name
}

func TestExtend(t *testing.T) {
	var p NameInterface = createChild()
	if reflect.TypeOf(p) == reflect.TypeOf(&Child{}) {
		child := p.(*Child)
		print(child.Name)
	}
}

func createChild() NameInterface {
	var chd = &Child{Parent: Parent{Name: "p", Desc: "p desc"}, Type: "child"}
	return chd
}
