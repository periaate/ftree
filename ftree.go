package ftree

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileTree is a tree of directories and files.
type FileTree struct {
	Root *Dir `json:"head"`
}

// Find finds a directory in the filetree by its path.
// The path requested must define the full path relative to the root directory.
func (ft *FileTree) Find(path string) *Dir {
	path = filepath.ToSlash(filepath.Clean(path))
	path = strings.TrimPrefix(path, "/")
	return ft.Root.find(strings.Split(path, "/"))
}

// FindFile finds a file in the filetree by its path.
func (ft *FileTree) FindFile(path string) Entry {
	fp := filepath.Dir(path)
	d := ft.Find(fp)
	if d == nil {
		return nil
	}
	return d.Files[filepath.Base(path)]
}

// Build builds a filetree with the given path as root directory.
func Build(root string) (*FileTree, error) {
	root, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return nil, err
	}
	vfs := os.DirFS(root)
	ft := &FileTree{}

	d := &Dir{
		Fp: Fpath{
			abs: root,
			rel: ".",
			dir: filepath.Base(root),
		},
		Children: make(map[string]*Dir),
		Files:    map[string]Entry{},
	}
	ft.Root = d

	d, err = d.build(vfs, ft.Root.Fp)
	if err != nil {
		return nil, err
	}
	ft.Root = d
	return ft, nil
}

// Dir is a directory in the filetree.
// It contains a map of its children directories and a map of its files.
type Dir struct {
	// Fp contains the absolute path, relative path and current directory name
	Fp       Fpath            `json:"fpath,omitempty"`
	Children map[string]*Dir  `json:"children,omitempty"`
	Files    map[string]Entry `json:"files,omitempty"`
}

type Entry interface {
	// Ext returns the file extension in lowercase
	Ext() string

	// Abs returns the absolute path to the file
	Abs() string

	// Rel returns the relative path from the root of where it was scanned
	// to the file. This is the path that is used to find the file in the
	// FileTree.
	Rel() string

	// Dir returns the directory name of the file
	Dir() string

	fs.DirEntry
}

type file struct {
	Fpath
	fs.DirEntry
}

// Ext returns the file extension in lowercase
func (f *file) Ext() string { return strings.ToLower(filepath.Ext(f.Name())) }

type Fpath struct {
	abs string
	rel string
	dir string
}

// Abs returns the absolute path to the file
func (f *Fpath) Abs() string { return f.abs }

// Rel returns the relative path from the root of where it was scanned
// to the file. This is the path that is used to find the file in the
// FileTree.
func (f *Fpath) Rel() string { return f.rel }

// Dir returns the directory name of the file
func (f *Fpath) Dir() string { return f.dir }

// build recursively builds the file tree
func (n *Dir) build(vfs fs.FS, f Fpath) (*Dir, error) {
	r, err := fs.ReadDir(vfs, ".")
	if err != nil {
		return nil, err
	}

	for _, v := range r {
		if v.IsDir() {
			fp := Fpath{
				abs: filepath.Join(f.Abs(), v.Name()),
				rel: filepath.Join(f.Rel(), v.Name()),
				dir: v.Name(),
			}
			t := &Dir{
				Fp:       fp,
				Children: make(map[string]*Dir),
				Files:    map[string]Entry{},
			}
			sfs, err := fs.Sub(vfs, v.Name())
			if err != nil {
				return nil, err
			}

			c, err := t.build(sfs, fp)
			if err != nil || c == nil {
				continue
			}
			n.Children[v.Name()] = c
		}
		n.Files[v.Name()] = &file{
			Fpath: Fpath{
				abs: filepath.Join(f.Abs(), v.Name()),
				rel: filepath.Join(f.Rel(), v.Name()),
				dir: f.Dir(),
			},
			DirEntry: v,
		}
	}
	return n, nil
}

// find recursively finds a directory in the file tree
func (n *Dir) find(path []string) *Dir {
	if len(path) == 0 {
		return n
	}

	if _, ok := n.Children[path[0]]; ok {
		return n.Children[path[0]].find(path[1:])
	}

	return nil
}

// NewFpath creates a new Fpath struct from a path
func NewFpath(root string) (*Fpath, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	abs = filepath.ToSlash(abs)
	rel := filepath.ToSlash(filepath.Clean(root))
	cur := filepath.Base(root)
	return &Fpath{
		abs: abs,
		rel: rel,
		dir: cur,
	}, nil
}

// Traverse traverses the file tree recursively, calling the function f on each
// directory.
func (d *Dir) Traverse(f func(*Dir)) {
	f(d)
	for _, v := range d.Children {
		v.Traverse(f)
	}
}

// Traverse calls traverse on the root directory of the file tree.
func (ft *FileTree) Traverse(f func(*Dir)) {
	ft.Root.Traverse(f)
}
