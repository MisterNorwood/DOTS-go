package parsers

import "fmt"

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

func PrintTarget(target Target) {
	fmt.Print("{")
	printMapKey("Aliases", target.Aliases)
	fmt.Print("; ")
	printMapKey("Mails", target.Mails)
	fmt.Print("; ")
	printMapKey("Commits", target.Commits)
	fmt.Print("}")
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
