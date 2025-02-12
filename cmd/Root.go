package cmd

import (
	"context"
	"fmt"
	. "github.com/MisterNorwood/DOTS-go/pkg/exporters"
	"github.com/MisterNorwood/DOTS-go/pkg/utils"
	"github.com/urfave/cli/v3"
	"os"
	"strings"
)

func Execute() {
	app := &cli.Command{
		Name:    "dots-go",
		Usage:   "OSINT tool for scraping github repositories",
		Version: "1.0.0",
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
				Action: func(ctx context.Context, cli *cli.Command, file string) error {
					_, err := pathExists(file)
					if err != nil {
						return err
					}

					return nil
				},
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
				Action: func(ctx context.Context, cli *cli.Command, dir string) error {
					_, err := pathExists(dir)
					if err != nil {
						return err
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    "exportForm",
				Value:   "TXT",
				Usage:   "Forms to export data to (STDOUT, CSV, XLS, TXT, JSON, XML, ALL)",
				Aliases: []string{"e"},
			},
			&cli.BoolFlag{
				Name:    "stripNoreply",
				Value:   true,
				Usage:   "Stripping default noreply anonymous mails",
				Aliases: []string{"n"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			splash()
			utils.MakeCacheDir()

			var sourceFlags []any
			sourceFlags = append(sourceFlags, cmd.String("file"))
			sourceFlags = append(sourceFlags, cmd.StringSlice("links"))
			sourceFlags = append(sourceFlags, cmd.String("repoDir"))

			var argContext ArgContext

			eMethod := verifySources(sourceFlags, &argContext.sourceMethod)
			if eMethod != nil {
				return eMethod
			}
			fmt.Println("Method type: ", SourceMethod(argContext.sourceMethod))

			eExportForm := verifyExports(cmd.String("exportForm"), &argContext.exportFormats)
			if eExportForm != nil {
				return eExportForm
			}

			processSource(cmd, argContext)
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func splash() {
	fmt.Print("Splash goes here later\n")
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err

}

func verifyExports(e string, exports *[]ExportFormat) error {
	formatMap := map[string]ExportFormat{
		"STDOUT": STDOUT,
		"TXT":    TXT,
		"CSV":    CSV,
		"XLS":    XLS,
		"XML":    XML,
		"JSON":   JSON,
		"ALL":    ALL,
	}
	var unknown []string
	var result []ExportFormat

	exportFormats := strings.Split(e, ",")
	for _, format := range exportFormats {
		if export, ok := formatMap[strings.ToUpper(format)]; ok {
			if export == ALL { //Early ALL catch
				*exports = []ExportFormat{ALL}
				return nil
			}
			result = append(result, export)
		} else {
			unknown = append(unknown, format)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("Unknown formats: %v", unknown)
	}

	if len(result) == 0 {
		result = append(result, TXT)
	}

	*exports = result
	return nil
}

func verifySources[T any](sourceFlags []T, method *SourceMethod) error {
	selected := 0

	for i, flag := range sourceFlags {
		switch flagType := any(flag).(type) {
		case string:
			if flagType != "" {
				selected++
				//Due to how things are sorted, loop 0 (1) will always be file, and loop 2 (3) will always be repo
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
