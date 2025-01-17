package parsers

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func ParseLog(rawLog string, targetDB *[]Target) {
	lines := splitLines(rawLog)
	sort.Strings(lines)

	for i, line := range lines {
		dataSlice := strings.Split(line, ";")
		//TODO: Slice 0 seems to be empty every time
		if len(dataSlice) != 3 {
			fmt.Printf("Error: Invalid data slice %d: %s! Skipping...\n", i, dataSlice)
			continue
		}

		alias := dataSlice[0]
		mail := dataSlice[1]
		commit := dataSlice[2]

		targetWithAlias, targetWithMail := findTargets(*targetDB, alias, mail)

		//Important that this is done here, as loops with range are copies of the data, hence now
		//I merely do the check first and cache a pointer to the chosen target.
		//Totally haven't caused a memory leak here before  due to this^
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
			if isMail, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, mail); isMail == true && !strings.Contains(mail, "users.noreply.github.com") {
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
