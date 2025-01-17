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

func TargetPrint(target Target) {
	fmt.Print("[")
	fmt.Print("Aliases:")
	for alias := range target.Aliases {
		fmt.Print(alias + ",")
	}
	fmt.Print("Mails:")
	for mail := range target.Mails {
		fmt.Print(mail + ",")
	}
	fmt.Print("Commits:")
	for commit := range target.Commits {
		fmt.Print(commit + ",")
	}
	fmt.Print("]")
}
