package ftree_test

import (
	"fmt"
	"ftree"
	"path/filepath"
	"testing"
)

var fullpath = "dir1/dir2/dir3/dir4/dir5/dir6/dir7/dir8/dir9/dir10"

func TestFpath(t *testing.T) {
	root := "./test/dir1"
	fp, err := ftree.NewFpath(root)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	abs = filepath.ToSlash(abs)
	rel := "test/dir1"
	cur := "dir1"
	if fp.Abs() != abs {
		t.Errorf("fp.abs : expected %s, got %s", abs, fp.Abs())
	}

	if fp.Rel() != rel {
		t.Errorf("fp.rel : expected %s, got %s", rel, fp.Rel())
	}

	if fp.Dir() != cur {
		t.Errorf("fp.cur : expected %s, got %s", cur, fp.Dir())
	}
}

func TestFtree(t *testing.T) {
	ft, err := ftree.Build("./test")
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if ft == nil {
		t.Errorf("filetree is nil")
	}

	d := ft.Find(fullpath)
	if d == nil {
		t.Errorf("filetree was not built correctly")
	}

	f := ft.FindFile(fmt.Sprintf("%s/test2.json", fullpath))
	if f == nil {
		t.Errorf("filetree was not built correctly")
	}
}
