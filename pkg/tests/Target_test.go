package parsers_test

import (
	. "github.com/MisterNorwood/DOTS-go/pkg/parsers"
	"reflect"
	"testing"
)

var testTarget = Target{
	Aliases: map[string]struct{}{
		"alias1": {},
	},
	Mails: map[string]struct{}{
		"test@example.com":  {},
		"dev@domain.com":    {},
		"admin@example.com": {},
	},
	Commits: map[string]struct{}{
		"commit1": {},
		"commit2": {},
		"commit3": {},
		"commit4": {},
	},
}

func TestTargetToSlice(t *testing.T) {
	expected := []string{"alias1", "admin@example.com, dev@domain.com, test@example.com", "commit1, commit2, commit3, commit4"}
	result := testTarget.ToSlice()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf(`expected "%s", got result "%s"`, expected, result)
	}
}

func TestTargetToMapSlice(t *testing.T) {
	expected := map[string][]string{
		"Aliases": {"alias1"},
		"Commits": {"commit1", "commit2", "commit3", "commit4"},
		"Mails":   {"admin@example.com", "dev@domain.com", "test@example.com"},
	}
	result := testTarget.ToMapSlice()
	if !reflect.DeepEqual(expected, result) {
		t.Errorf(`expected:  "%s", got result:  "%s"`, expected, result)
	}
}

func TestTargetToCsv(t *testing.T) {
	expected := `alias1; admin@example.com, dev@domain.com, test@example.com; commit1, commit2, commit3, commit4`
	result := testTarget.ToCsv()
	if expected != result {
		t.Errorf(`expected:  "%s", got result:  "%s"`, expected, result)
	}
}

func TestStripNoreply(t *testing.T) {
	targets := []Target{
		{
			Mails: map[string]struct{}{
				"legit@example.com":             {},
				"user@users.noreply.github.com": {},
				"user@noreply.codeberg.org":     {},
				"user@no-reply.gitlab.com":      {},
			},
		},
	}

	StripNoreply(&targets, false)

	if len(targets[0].Mails) != 1 {
		t.Errorf("Expected 1 mail, got %d: %v", len(targets[0].Mails), targets[0].Mails)
	}
	if _, ok := targets[0].Mails["legit@example.com"]; !ok {
		t.Error("legit@example.com should still be present")
	}
}
