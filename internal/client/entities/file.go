package entities

type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
}
type Files []File

func NewFile(name, path string, size int64) *File {
	return &File{
		Name: name,
		Size: size,
		Path: path,
	}
}
