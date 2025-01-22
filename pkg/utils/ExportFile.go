package utils

import (
	"errors"
	"fmt"
	"os"
)

var dbDir string

const dbDirVar string = "DOTS_EXPORT"

func GetDBDir() string {
	if dbDir != "" {
		return dbDir
	} else {
		MakeDBDir()
		return dbDir
	}

}

func MakeDBDir() {
	dir, exists := os.LookupEnv(dbDir)
	if exists {
		err := os.Mkdir(dir, 0755)
		if err == nil || errors.Is(err, os.ErrExist) {
			dbDir = dir
			return
		} else {
			fmt.Errorf("Warning: %q failed to be created or accessed. Reverting to default dbDir directory\n Error: %s ", dir, err)
		}
	} else {
		if err := os.Mkdir("DotsExports", 0775); err != nil && !errors.Is(err, os.ErrExist) {
			panic(err)
		}

		dbDir = "DotsExports"
		return
	}

}
