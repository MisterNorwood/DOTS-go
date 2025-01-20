package cmd

import (
	"github.com/MisterNorwood/DOTS-go/pkg/exporters"
)

type SourceMethod int

const (
	SourceFile SourceMethod = iota
	SourceLink
	SourceRepo
)

type ArgContext struct {
	sourceMethod  SourceMethod
	exportFormats []exporters.ExportFormat
}
