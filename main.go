package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MisterNorwood/DOTS-go/pkg/utils"

	"github.com/urfave/cli/v3"
)

type SourceMethod int

const (
	SourceFile SourceMethod = iota
	SourceLink
	SourceRepo
)

func main() {
	app := &cli.Command{
		Name:  "dotsCli",
		Usage: "Scrape repository data and export it in various formats",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "threads",
				Value:   4,
				Usage:   "Number of threads for multithreaded workloads",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:    "file",
				Usage:   "Plain text file with repository links",
				Aliases: []string{"f"},
			},
			&cli.StringSliceFlag{
				Name:    "links",
				Usage:   "Links of repositories to be scraped (multiple allowed)",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "repoDir",
				Usage:   "Directory containing repositories to scrape",
				Aliases: []string{"r"},
			},
			&cli.StringFlag{
				Name:    "exportForm",
				Value:   "TXT",
				Usage:   "Forms to export data to (CSV, XLS, TXT, JSON)",
				Aliases: []string{"e"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			splash()
			CACHE := utils.MakeCacheDir()
			println("\n", CACHE)

			var sourceFlags []any
			sourceFlags = append(sourceFlags, cmd.String("file"))
			sourceFlags = append(sourceFlags, cmd.StringSlice("links"))
			sourceFlags = append(sourceFlags, cmd.String("repoDir"))

			var method SourceMethod

			e := verifySources(sourceFlags, &method)
			if e != nil {
				return e
			}
			fmt.Println("Method type: ", SourceMethod(method))

			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func splash() {
	fmt.Print("Splash goes here later\n")
}

// TODO: This is retatded, restrict it to only strings and slices.
func verifySources[T any](sourceFlags []T, method *SourceMethod) error {
	selected := 0

	for i, flag := range sourceFlags {
		switch flagType := any(flag).(type) {
		case string:
			if flagType != "" {
				selected++
				if i == 0 {
					*method = SourceMethod(0)
				} else {
					*method = SourceMethod(2)
				}
			}
		case []string:
			if flagType != nil && len(flagType) != 0 {
				selected++
				*method = SourceMethod(1)
			}
		default:
			e := fmt.Errorf("Invalid input type!")
			return e
		}
	}
	fmt.Printf("Current selected: %d\n", selected)

	if selected > 1 {
		fmt.Println("Flags; file, links and repoDir are mutually exclusive!")
		e := fmt.Errorf("Too many sources")
		return e
	} else if selected == 0 {
		fmt.Println("No sources provided!")
		e := fmt.Errorf("No sources")
		return e
	}
	return nil
}
