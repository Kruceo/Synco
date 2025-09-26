package utils

import (
	"crypto/sha256"
	"io"
	"os"
)

func FastCopy(from string, to string) {
	file, err := os.OpenFile(from, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	file, err = os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	file.Truncate(0)
	defer file.Close()
	file.Write(data)
}

func Sha256(filePath string) [32]byte {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	return sha256.Sum256(data)
}
