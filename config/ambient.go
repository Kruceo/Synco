package config

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

var SyncoDir string
var BlobPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	SyncoDir = os.Getenv("SYNCO_DIR")
	if SyncoDir == "" {
		SyncoDir = path.Join(home, ".synco")
	}
	SyncoDir, err = filepath.Abs(SyncoDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of SYNCO_DIR: %s\n%v", os.Getenv("SYNCO_DIR"), err)
		os.Exit(1)
	}

	BlobPath = path.Join(SyncoDir, "blob")
	if _, err := os.Stat(BlobPath); os.IsNotExist(err) {
		err := os.MkdirAll(BlobPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	initLog()
}
