package exporters

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	. "github.com/MisterNorwood/DOTS-go/pkg/parsers"
	"github.com/xuri/excelize/v2"
)

type ExportFormat int

const (
	TXT ExportFormat = iota
	STDOUT
	CSV
	XLS
	XML
	JSON
	ALL
)

// Lazy init to avoid init self reference error
var ActionMap map[ExportFormat]func([]Target, string) error

func Init() {
	ActionMap = map[ExportFormat]func([]Target, string) error{
		STDOUT: func(targetDB []Target, filepath string) error {
			for _, target := range targetDB {
				target.PrintFancy()
				fmt.Print("\n")
			}
			return nil
		},
		TXT:  ExportTXT,
		CSV:  ExportCSV,
		XLS:  ExportXLS,
		XML:  ExportXML,
		JSON: ExportJSON,
		ALL: func(targetDB []Target, filepath string) error {
			for format, action := range ActionMap {
				if format != ALL && action != nil { // Avoid infinite recursion with ALL
					if err := action(targetDB, filepath); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
}

func ExportCSV(targetDB []Target, filepath string) error {
	fmt.Println("Exporting to " + filepath + ".csv")
	csvFile, err := os.Create(filepath + ".csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := bufio.NewWriter(csvFile)
	defer writer.Flush()

	csvWriter := csv.NewWriter(writer)
	for _, target := range targetDB {
		if err := csvWriter.Write(target.ToSlice()); err != nil {
			return err
		}
	}
	return nil
}

func ExportTXT(targetDB []Target, filepath string) error {
	fmt.Println("Exporting to " + filepath + ".txt")
	txtFile, err := os.Create(filepath + ".txt")
	if err != nil {
		return err
	}
	defer txtFile.Close()
	writer := bufio.NewWriter(txtFile)
	defer writer.Flush()

	for _, target := range targetDB {
		if _, err := writer.WriteString(target.ToCsv() + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// FIXME: broken formatting
func ExportXLS(targetDB []Target, filepath string) error {
	fmt.Println("Exporting to " + filepath + ".xlsx")
	file := excelize.NewFile()
	rows := [][]string{}
	for _, target := range targetDB {
		rows = append(rows, target.ToSlice())
	}

	for rowIndex, row := range rows {
		for colIndex, cellVal := range row {
			alphaIndex, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return err
			}
			cell := alphaIndex + strconv.Itoa(rowIndex+1)
			file.SetCellStr("Sheet1", cell, cellVal)
		}
	}
	return file.SaveAs(filepath + ".xlsx")
}

func ExportJSON(targetDB []Target, filepath string) error {
	fmt.Println("Exporting to " + filepath + ".json")
	jsonFile, err := os.Create(filepath + ".json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	writer := bufio.NewWriter(jsonFile)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // For prettyPrint

	for _, target := range targetDB {
		if err := encoder.Encode(target.ToMapSlice()); err != nil {
			return err
		}

	}

	return nil
}

type TargetXML struct {
	XMLName xml.Name `xml:"Target"`
	Aliases []string `xml:"Aliases>Alias"`
	Mails   []string `xml:"Mails>Mail"`
	Commits []string `xml:"Commits>Commit"`
}

func ExportXML(targetDB []Target, filepath string) error {
	fmt.Println("Exporting to " + filepath + ".xml")
	xmlFile, err := os.Create(filepath + ".xml")
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	writer := bufio.NewWriter(xmlFile)
	defer writer.Flush()

	var xmlTargets []TargetXML
	for _, target := range targetDB {
		xmlTargets = append(xmlTargets, TargetXML{
			Aliases: target.AliasesAsSlice(),
			Mails:   target.MailsAsSlice(),
			Commits: target.CommitsAsSlice(),
		})
	}

	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "  ") // Indent for readabilty

	encoder.Encode(xml.Header)
	if err := encoder.Encode(xmlTargets); err != nil {
		return err
	}

	return nil
}
