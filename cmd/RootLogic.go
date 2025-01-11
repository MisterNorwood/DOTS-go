package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/MisterNorwood/DOTS-go/pkg/executors"
	"github.com/urfave/cli/v3"
)

func processSource(cmd *cli.Command, method SourceMethod) {
	var extractionDirs []string
	threads := cmd.Int("threads")

	runWorkers := func(inputProvider func(chan string), dirCollector func(chan string)) {
		dirChannel := make(chan string)
		linkChannel := make(chan string, threads)
		var wg sync.WaitGroup

		for i := 0; i < int(threads); i++ {
			wg.Add(1)
			go worker(i, linkChannel, &wg, dirChannel)
		}

		go func() {
			inputProvider(linkChannel)
			close(linkChannel)
		}()

		go func() {
			wg.Wait()
			close(dirChannel)
		}()

		dirCollector(dirChannel)
	}

	// Input and directory collector functions
	appendDirs := func(dirChannel chan string) {
		for dir := range dirChannel {
			extractionDirs = append(extractionDirs, dir)
		}
	}

	switch method {
	case 0: // Case for input file
		runWorkers(
			func(linkChannel chan string) { // Input provider for file
				file, err := os.Open(cmd.String("file"))
				if err != nil {
					panic(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					linkChannel <- scanner.Text()
				}
				if err := scanner.Err(); err != nil {
					fmt.Println("Error reading file:", err)
				}
			},
			appendDirs,
		)

	case 1: // Case for array of links
		runWorkers(
			func(linkChannel chan string) {
				links := cmd.StringSlice("links")
				for _, link := range links {
					linkChannel <- link
				}
			},
			appendDirs,
		)

	case 2: // Case for local repo directory
		reposDir := cmd.String("repoDir")
		repoDirEntries, err := os.ReadDir(reposDir)
		if err != nil {
			panic(err)
		}

		for _, entry := range repoDirEntries {
			if entry.IsDir() {
				if exists, err := pathExists(filepath.Join(entry.Name(), ".git")); exists && err == nil {
					extractionDirs = append(extractionDirs, entry.Name())
				}
			}
		}
	}
}

func worker(id int, linkChannel <-chan string, wg *sync.WaitGroup, dirChannel chan<- string) {
	defer wg.Done()
	for link := range linkChannel {
		repoDir := executors.RetriveRepositories(link)
		if repoDir != "" {
			dirChannel <- repoDir
			fmt.Printf("Worker %d: Repo %q clone succeded\n", id, link)
		} else {
			fmt.Printf("Worker %d: Repo %q clone failed, Skipping...\n", id, link)
		}
	}

}
