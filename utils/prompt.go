package utils

import (
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func GetOptionsFromDir(di string) []huh.Option[string] {
	files, err := os.ReadDir(di)
	fileOptions := []huh.Option[string]{}
	if err != nil {
		panic(err)
	}
	for _, v := range files {
		if !v.IsDir() {
			realPath, _ := filepath.Abs(path.Join(di, v.Name()))
			fileOptions = append(fileOptions, huh.NewOption(v.Name(), realPath))
		}
	}
	return fileOptions
}
