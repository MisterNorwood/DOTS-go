package parsers

import (
	"fmt"
	"sort"
	"strings"
)

func ParseLog(rawLog string, targetDB *[]Target) {
	fmt.Print(rawLog)
	lines := splitLines(rawLog)
	sort.Strings(lines)
	fmt.Println(lines)
	fmt.Println(len(lines))

	for i, line := range lines {
		dataSlice := strings.Split(line, ";")
		if len(dataSlice) != 3 {
			fmt.Printf("Error: Invalid data slice %d: %s! Skipping... (%s)\n", i, dataSlice, line)
			continue
		}

		alias := dataSlice[0]
		mail := dataSlice[1]
		commit := dataSlice[2]

		// Attempt to find existing targets
		targetWithAlias, targetWithMail := findTargets(*targetDB, alias, mail)

		// Add data to appropriate targets or create a new one
		if targetWithAlias != nil && targetWithMail != nil {
			TargetAdd(targetWithMail.Commits, commit)
		} else if targetWithAlias != nil && targetWithMail == nil {
			TargetAdd(targetWithAlias.Mails, mail)
			TargetAdd(targetWithAlias.Commits, commit)
		} else if targetWithAlias == nil && targetWithMail != nil {
			TargetAdd(targetWithMail.Aliases, alias)
			TargetAdd(targetWithMail.Commits, commit)
		} else {
			newTarget := NewTarget(alias, mail, commit)
			*targetDB = append(*targetDB, *newTarget)
		}
	}
}

func findTargets(targetDB []Target, alias, mail string) (*Target, *Target) {
	var targetWithAlias, targetWithMail *Target
	for i := range targetDB {
		target := &targetDB[i]
		if TargetContains(target.Aliases, alias) {
			targetWithAlias = target
		}
		if TargetContains(target.Mails, mail) {
			targetWithMail = target
		}
	}
	return targetWithAlias, targetWithMail
}

func StripNoreply(targetDB *[]Target, keepStripped bool) {
	var strippedTargetDB []Target
	for _, target := range *targetDB {
		filteredMails := make(map[string]struct{})
		for mail := range target.Mails {
			if !strings.Contains(mail, "users.noreply.github.com") {
				filteredMails[mail] = struct{}{}
			}
		}
		target.Mails = filteredMails
		if len(target.Mails) > 0 || keepStripped {
			strippedTargetDB = append(strippedTargetDB, target)
		}
	}
	*targetDB = strippedTargetDB

}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}
