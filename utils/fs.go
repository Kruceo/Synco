package utils

import (
	"crypto/sha256"
	"io"
	"os"
)

func FastCopy(from string, to string) {
	file, err := os.OpenFile(from, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	file, err = os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	file.Truncate(0)
	defer file.Close()
	file.Write(data)
}

func Sha256(filePath string) [32]byte {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	return sha256.Sum256(data)
}

func AppendFilesToBuffer(filePaths []string) ([]byte, error) {
	buffer := []byte{}

	for _, filePath := range filePaths {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(f)
		if err != nil {
			f.Close()
			return nil, err
		}

		f.Close()

		// O buffer recebe o caminho do arquivo e o conteúdo
		buffer = append(buffer, []byte(filePath)...)
		buffer = append(buffer, data...)
	}

	return buffer, nil
}
