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
	var rawLogs []string
	threads := cmd.Int("threads")

	runWorkers := func(inputProvider func(chan string), dirCollector func(chan string, *[]string)) {
		dirChannel := make(chan string, threads)
		linkChannel := make(chan string, threads)
		var wg sync.WaitGroup

		for i := 0; i < int(threads); i++ {
			wg.Add(1)
			go cloneWorker(i, linkChannel, &wg, dirChannel)
		}

		go func() {
			inputProvider(linkChannel)
			close(linkChannel)
		}()

		go func() {
			wg.Wait()
			close(dirChannel)
		}()

		dirCollector(dirChannel, &rawLogs)

		for _, rawLog := range rawLogs {
			fmt.Printf("Print Raw Log: %s\n", rawLog)
		}
	}

	// Input and directory collector functions
	retriveLogs := func(dirChannel chan string, rawLog *[]string) {
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 0; i < int(threads); i++ {
			wg.Add(1)
			go extractWorker(i, dirChannel, &wg, &rawLogs, &mu)
		}

		wg.Wait()

	}

	switch method {
	case 0: // Case for input file
		runWorkers(
			func(linkChannel chan string) {
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
			retriveLogs,
		)

	case 1: // Case for array of links
		runWorkers(
			func(linkChannel chan string) {
				links := cmd.StringSlice("links")
				for _, link := range links {
					linkChannel <- link
				}
			},
			retriveLogs,
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

func cloneWorker(id int, linkChannel <-chan string, wg *sync.WaitGroup, dirChannel chan<- string) {
	defer wg.Done()
	fmt.Printf("Clone worker %d started\n", id)
	for link := range linkChannel {
		repoDir := executors.RetriveRepositories(link)
		if repoDir != "" {
			dirChannel <- repoDir
			fmt.Printf("Clone Worker %d: Repo %q clone succeded\n", id, link)
		} else {
			fmt.Printf("Clone Worker %d: Repo %q clone failed, Skipping...\n", id, link)
		}
	}

}

func extractWorker(id int, dirChannel <-chan string, wg *sync.WaitGroup, rawLogs *[]string, mu *sync.Mutex) {
	defer wg.Done()
	fmt.Printf("Extract worker %d started\n", id)
	for repo := range dirChannel {
		rawLog := executors.RetriveLogStream(repo)
		if rawLog != "" {
			mu.Lock()
			*rawLogs = append(*rawLogs, rawLog)
			mu.Unlock()
			fmt.Printf("Extract Worker %d: Log for %q retrived\n", id, repo)
		} else {
			fmt.Printf("Extract Worker %d: Log for %q retrival failed, Skipping...\n", id, repo)
		}
	}

}
