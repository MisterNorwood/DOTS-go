package executors

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/MisterNorwood/DOTS-go/pkg/utils"
)

func RetriveRepositories(link string) string {
	repoSlice := strings.Split(link, "/")
	repoDir := repoSlice[len(repoSlice)-1] + ".git/"

	if _, lookErr := os.Stat(utils.GetCacheDir() + "/" + repoDir); os.IsNotExist(lookErr) {
		gitClone := exec.Command("git", "clone", "--filter=blob:none", "--bare", link)
		gitClone.Dir = utils.GetCacheDir()

		if execError := gitClone.Run(); execError != nil {
			fmt.Printf("Error on clone of %q \n %s", link, execError)
			return ""
		}
	}
	return repoDir
}

func RetriveLogStream(repoDir string) string {
	fullPath := utils.GetCacheDir() + repoDir
	logCmd := exec.Command("git", "log", "--pretty=\"%an;%ae;%h\"")
	logCmd.Dir = fullPath
	var ioStream bytes.Buffer
	logCmd.Stdout = &ioStream
	if execErr := logCmd.Run(); execErr != nil {
		fmt.Printf("Error on log extraction of %q \n %s", logCmd.Dir, execErr)
		return ""
	}
	return ioStream.String()
}
