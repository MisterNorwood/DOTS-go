package cmd

import (
	"fmt"
	"github.com/MisterNorwood/DOTS-go/pkg/exporters"
	"github.com/MisterNorwood/DOTS-go/pkg/parsers"
	"github.com/MisterNorwood/DOTS-go/pkg/utils"
	"github.com/urfave/cli/v3"
	"time"
)

// FIXME: CSV, TXT and XML ok. XLS and JSON need fixes
func ExportLogic(cmd *cli.Command, tDB *[]parsers.Target, argContext *ArgContext) error {
	exporters.Init()
	path := utils.GetDBDir() + "/" + time.Now().Format("15-04-02-01-2006")

	for _, format := range argContext.exportFormats {
		if action, exists := exporters.ActionMap[format]; exists {
			action(*tDB, path)
			if format == exporters.ALL {
				break
			}
		} else {
			fmt.Printf("Unknown format: %v\n", format)
		}
	}

	return nil
}
