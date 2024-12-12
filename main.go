package main

import (
	"context"
	"fmt"
	"log"
	"github.com/MisterNorwood/DOTS-go/pkg/utils"
	"os"

	"github.com/urfave/cli/v3"
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

			var sourceFlags []interface{}
			sourceFlags = append(sourceFlags, cmd.String("file"))
			sourceFlags = append(sourceFlags, cmd.String("repoDir"))
			sourceFlags = append(sourceFlags, cmd.StringSlice("links"))
			selected := 0

			for _, flag := range sourceFlags {
				fmt.Printf("Current selected: %d; Current flag type %T\n", selected, flag)
				switch flagType := flag.(type) {
				case string:
					if flagType != "" {
						selected++
					}
				case []string:
					if flagType != nil && len(flagType) != 0 {
						selected++
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
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func splash() {
	fmt.Print("Splash goes here later\n")
}
verifySources(var sourceFlags)
