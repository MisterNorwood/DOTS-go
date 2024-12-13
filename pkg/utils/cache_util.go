package utils

import (
	"os"
)

var CACHE string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func MakeCacheDir() {

	dname, err := os.MkdirTemp("", "repos")
	check(err)
	CACHE = dname

}
