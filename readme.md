# Ftree
`ftree` is a Go package for building and traversing file trees. It provides a simple interface for building a tree of directories and files, and for walking the tree and performing actions on each file.


## Usage
### FileTree
#### Building a FileTree
To build a FileTree, use the Build function with any valid path:

```go
ft, err := ftree.Build("/path/to/root")
if err != nil {
    // Handle error
}
```
####  Traversing a FileTree
To traverse the FileTree and perform actions on each directory, use the Traverse method of the Dir type:

```go
ft.Traverse(func(d *ftree.Dir) {
    // Do something with the directory
})
```
#### Finding a Directory or File
To find a directory in the FileTree, use the Find method of FileTree :

```go
dir := ft.Find("/path/to/dir")
```
To find a file in the FileTree, use the FindFile method of FileTree:

```go
file := ft.FindFile("/path/to/file.txt")
```



### Walker
#### Building a Walker
To build a Walker, use the BuildWalker function with any valid path and one or more Stepper interfaces:

```go
walker, err := ftree.BuildWalker("/path/to/root", myStepper, myStepper2)
if err != nil {
	// Handle error
}
```


#### Making a Stepper
To make a Stepper, define a struct that implements the Stepper interface:

```go
type MyStepper struct{}

func (s *MyStepper) Wants(ext string) bool {
    return ext == ".txt"
}

func (s *MyStepper) Walk(e ftree.Entry, r io.Reader) error {
    // Do something with the file
    return nil
}
```

In this example, we define a custom Stepper called `MyStepper`. This Stepper only wants to process files with the `.txt` extension, as indicated by the `Wants` method. The `Walk` method of the Stepper performs some action on the file, such as reading its contents.

#### Walking a FileTree with Steppers
To walk a FileTree and perform actions on each file, use the Walk method of the Walker:

```go
walker, err := ftree.BuildWalker("/path/to/root",
	myStepper,
	myStepper2
)
if err != nil {
    // Handle error
}

walker.AddStepper(myStepper3)

err = walker.Walk()
if err != nil {
    // Handle error
}
```

The Walker can be built with one or more Stepper interfaces as parameters to BuildWalker and additional Steppers can be added with the AddStepper method. When the Walk method is called, the Walker will process each file with the Steppers that it contains, only reading files which are wanted. Each file is read only once, even if it is wanted by multiple Steppers.


#### Walking a FileTree
To walk a FileTree and perform actions on each file, use the Walker type:
```go
walker, err := ftree.BuildWalker("/path/to/root",
	myStepper1,
	myStepper2,
)
if err != nil {
    // Handle error
}

err = walker.Walk()
if err != nil {
    // Handle error
}
```

The BuildWalker function takes one or more Stepper interfaces as arguments. These Stepper interfaces define the actions to be taken on each file. Only files with the desired extensions are processed by the Stepper.

```go
type MyStepper struct{}

func (s *MyStepper) Wants(ext string) bool {
    return ext == ".txt"
}

func (s *MyStepper) Walk(e ftree.Entry, r io.Reader) error {
    // Do something with the file
    return nil
}

walker, err := ftree.BuildWalker("/path/to/root", &MyStepper{})
if err != nil {
    // Handle error
}

err = walker.Walk()
if err != nil {
    // Handle error
}
```