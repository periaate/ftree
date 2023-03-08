package ftree

import (
	"bytes"
	"io"
	"os"
)

// Walker is used to walk a FileTree and run Steppers on each file.
// This makes it easier to do things like read files which multiple
// Steppers may want to read.
type Walker struct {
	ft       *FileTree
	steppers []Stepper
}

// Stepper is an interface used to act in place of fs.WalkDirFunc.
// To determine if a Stepper wants to read a file, the extension of a file
// is passed onto the want method.
type Stepper interface {
	// Walk is ran on each file in the FileTree.
	Walk(e Entry, r io.Reader) error
	// Given an extension, returns true if the Stepper wants to read the file.
	Wants(ext string) bool
}

// NewWalker returns a new Walker with the given FileTree.
func NewWalker(ft *FileTree) *Walker {
	return &Walker{
		ft: ft,
	}
}

// BuildWalker returns a new Walker with the given root and steppers.
func BuildWalker(root string, steppers ...Stepper) (*Walker, error) {
	ft, err := Build(root)
	if err != nil {
		return nil, err
	}
	w := NewWalker(ft)
	for _, s := range steppers {
		w.AddStepper(s)
	}
	return w, nil
}

// AddStepper adds a Stepper to the Walker.
func (w *Walker) AddStepper(s Stepper) {
	w.steppers = append(w.steppers, s)
}

// Walk walks the FileTree and runs the Steppers on each file.
// Only files which are wanted are ever read.
// Walk is ran by traversing the filetree directory by directory.
func (w *Walker) Walk() error {
	var (
		r   = &bytes.Buffer{}
		err error
		b   []byte
	)

	w.ft.Root.Traverse(func(d *Dir) {
		for _, f := range d.Files {
			opened := false
			for _, s := range w.steppers {
				if s.Wants(f.Ext()) {
					if !opened {
						b, err = os.ReadFile(f.Abs())
						if err != nil {
							continue
						}
						opened = true
					}
					r.Write(b)
					s.Walk(f, r)
					r.Reset()
				}
			}
		}
	})
	return nil
}
