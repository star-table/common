package copyer

import (
	"testing"
)

type Source struct {
	Name                         string
}

type Destination struct {
	Name                         string
}


func TestCopy(t *testing.T) {
	s:= &Source{
		Name: "nico",
	}
	d := &Destination{}

	ss := &[]*Source{s,}
	dd := &[]*Destination{d,}
	Copy(ss, dd)
	t.Log((*dd)[0].Name)
}