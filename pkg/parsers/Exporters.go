package parsers

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

// TODO: Buffer and multithread this
func ExportCSV(targetDB []Target, filepath string) error {
	csvFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := bufio.NewWriter(csvFile)
	defer writer.Flush()

	csvWriter := csv.NewWriter(writer)
	for _, target := range targetDB {
		if err := csvWriter.Write(target.toSlice()); err != nil {
			return err
		}
	}
	return nil
}

func ExportTXT(targetDB []Target, filepath string) error {
	txtFile, err := os.Create(filepath)
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

// TODO: make it fancier
func ExportXLS(targetDB []Target, filepath string) error {
	file := excelize.NewFile()
	rows := [][]string{}
	for _, target := range targetDB {
		rows = append(rows, target.toSlice())
	}

	for rowIndex, row := range rows {
		for colIndex, cellVal := range row {
			alphaIndex, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				return err
			}
			cell := alphaIndex + fmt.Sprintf("%d", rowIndex+1)
			file.SetCellStr("Sheet1", cell, cellVal)
		}
	}
	return file.SaveAs(filepath)
}

func ExportJSON(targetDB []Target, filepath string) error {
	jsonFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	writer := bufio.NewWriter(jsonFile)
	defer writer.Flush()

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // For prettyPrint

	if err := encoder.Encode(targetDB); err != nil {
		return err
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
	xmlFile, err := os.Create(filepath)
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

	// Create an XML encoder
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "  ") // Indent for readabilty

	encoder.Encode(xml.Header)
	if err := encoder.Encode(xmlTargets); err != nil {
		return err
	}

	return nil
}
