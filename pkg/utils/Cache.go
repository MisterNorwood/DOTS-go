package utils

import (
	"errors"
	"fmt"
	"os"
)

var cache string

const cache_dir string = "DOTS_CACHE"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetCacheDir() string {
	if cache != "" {
		return cache
	} else {
		MakeCacheDir()
		return cache
	}

}

func MakeCacheDir() {
	dir, exists := os.LookupEnv(cache_dir)
	if exists {
		err := os.Mkdir(dir, 0755)
		if err == nil || errors.Is(err, os.ErrExist) {
			cache = dir
			return
		} else {
			fmt.Printf("Warning: %q failed to be created or accessed. Reverting to default cache directory\n Error: %s ", dir, err)
		}
	} else {
		dname, err := os.MkdirTemp("", "repos")
		check(err)
		cache = dname
		return
	}

}
