package parsers

import (
	"bytes"
	"fmt"
)

// Target struct
type Target struct {
	Aliases map[string]struct{}
	Mails   map[string]struct{}
	Commits map[string]struct{}
}

func NewTarget(alias, mail, commit string) *Target {
	return &Target{
		Aliases: map[string]struct{}{alias: {}},
		Mails:   map[string]struct{}{mail: {}},
		Commits: map[string]struct{}{commit: {}},
	}
}

func TargetContains(set map[string]struct{}, value string) bool {
	_, exists := set[value]
	return exists
}

func TargetAdd(set map[string]struct{}, value string) {
	set[value] = struct{}{}
}

func (target Target) PrintFancy() {
	fmt.Print("{")
	printMapKey("Aliases", target.Aliases)
	fmt.Print("; ")
	printMapKey("Mails", target.Mails)
	fmt.Print("; ")
	printMapKey("Commits", target.Commits)
	fmt.Print("}")
}

func (target Target) ToCsv() string {
	csvFormatter := func(items map[string]struct{}, buffer *bytes.Buffer) {
		first := true
		for item := range items {
			buffer.WriteString(item)
			if !first {
				buffer.WriteString(", ")
			}
			first = false
		}
	}
	var line bytes.Buffer
	csvFormatter(target.Aliases, &line)
	line.WriteString(";")
	csvFormatter(target.Mails, &line)
	line.WriteString(";")
	csvFormatter(target.Commits, &line)
	return line.String()
}

func (target Target) ToSlice() []string {
	formatter := func(items map[string]struct{}) string {
		var line bytes.Buffer
		count := 0
		target := len(items) - 1
		for item := range items {
			line.WriteString(item)
			if len(items) <= 1 || count == target {
				break
			}
			line.WriteString(", ")
			count++
		}
		return line.String()
	}
	var slice []string
	slice = append(slice, formatter(target.Aliases))
	slice = append(slice, formatter(target.Mails))
	slice = append(slice, formatter(target.Commits))
	return slice
}

func (target Target) ToMapSlice() map[string][]string {
	formatter := func(items map[string]struct{}) []string {
		var slice []string
		for item := range items {
			slice = append(slice, item)
		}
		return slice
	}
	mapSlice := make(map[string][]string)

	mapSlice["Aliases"] = formatter(target.Aliases)
	mapSlice["Mails"] = formatter(target.Mails)
	mapSlice["Commits"] = formatter(target.Commits)
	return mapSlice
}

func printMapKey(label string, items map[string]struct{}) {
	fmt.Print(label + ": ")
	first := true
	for item := range items {
		if !first {
			fmt.Print(", ")
		}
		fmt.Print(item)
		first = false
	}
}

func (target Target) AliasesAsSlice() []string {
	var aliases []string
	for alias := range target.Aliases {
		aliases = append(aliases, alias)
	}
	return aliases
}

func (target Target) MailsAsSlice() []string {
	var mails []string
	for mail := range target.Mails {
		mails = append(mails, mail)
	}
	return mails
}

func (target Target) CommitsAsSlice() []string {
	var commits []string
	for commit := range target.Commits {
		commits = append(commits, commit)
	}
	return commits
}
