package cmd

type SourceMethod int
type ExportFormat int

const (
	SourceFile SourceMethod = iota
	SourceLink
	SourceRepo
)

const (
	TXT ExportFormat = iota
	CSV
	XLS
	XML
	JSON
	ALL
)

type ArgContext struct {
	sourceMethod  SourceMethod
	exportFormats []ExportFormat
}
