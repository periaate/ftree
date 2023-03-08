package ftree_test

import (
	"encoding/json"
	"ftree"
	"io"
	"testing"
)

type stepper struct {
	name   string
	failed bool
	read   int
	t      *testing.T
}

func makeStepper(name string, t *testing.T) *stepper {
	return &stepper{
		name: name,
		t:    t,
	}
}

func (s *stepper) Walk(e ftree.Entry, r io.Reader) error {
	currentRead := s.read
	var v map[string]interface{}
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		s.failed = true
		s.t.Errorf("%s failed to decode from reader at %s succesfully", s.name, e.Rel())
		return nil
	}
	for k := range v {
		if !(k == "test" || k == "hello") {
			s.failed = true
			s.t.Errorf("expected key to be test or hello, got %s", k)
			return nil
		}
		s.read++
	}
	if currentRead == s.read {
		s.failed = true
		s.t.Errorf("%s failed to read from reader at %s succesfully", s.name, e.Rel())
	}
	return nil
}

func (s *stepper) Wants(ext string) bool {
	return ext == ".json"
}

func TestWalker(t *testing.T) {
	ft, err := ftree.Build("./test")
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if ft == nil {
		t.Errorf("filetree is nil")
	}

	w := ftree.NewWalker(ft)
	w.AddStepper(makeStepper("stepper1", t))
	w.AddStepper(makeStepper("stepper2", t))
	w.Walk()
}

func TestBuildWalker(t *testing.T) {
	w, err := ftree.BuildWalker("./test",
		makeStepper("stepper1", t),
		makeStepper("stepper2", t),
	)
	if err != nil {
		t.Errorf("error building walker (filetree): %s", err)
	}
	w.Walk()
}
