package utils

import (
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func MakeCacheDir() string {

	dname, err := os.MkdirTemp("", "repos")
	check(err)
	return dname

}
