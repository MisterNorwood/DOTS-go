package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/MisterNorwood/DOTS-go/pkg/executors"
	"github.com/urfave/cli/v3"
)

func processSource(cmd *cli.Command, method SourceMethod) {

	var extractionDirs []string
	threads := cmd.Int("threads")

	switch method {
	case 0:
		// Case for input file
		// Worker Spawning
		dirChannel := make(chan string)
		linkChannel := make(chan string, threads)
		var wg sync.WaitGroup

		for i := 0; i < int(threads); i++ {
			wg.Add(1)
			go worker(i, linkChannel, &wg, dirChannel)

		}
		//File Opening
		file, err := os.Open(cmd.String("file"))
		if err != nil {
			fmt.Print(err)
			panic(err)
		}
		defer file.Close()
		//Scanner
		scanner := bufio.NewScanner(file)
		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				linkChannel <- line
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
			}
			close(linkChannel) // Close channel after sending all lines
		}()

		// Collect results from workers
		go func() {
			wg.Wait()
			close(dirChannel)
		}()

		for dir := range dirChannel {
			extractionDirs = append(extractionDirs, dir)
			fmt.Println(dir)
		}

	case 1:
	//  Case for array of links
	//Straightforward to add, just put it onto queue and read off of it
	case 2:
		// Case for local repo directory

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
