package projectdetails

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/mapper"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/summary"
	"github.com/reflect-homini/stora/internal/domain/timeframe"
)

func entryToItem(e entry.Entry) project.ProjectItem {
	return project.ProjectItem{
		BaseDTO:   mapper.BaseToDTO(e.BaseEntity),
		ProjectID: e.ProjectID,
		ItemType:  project.ItemTypeEntry,
		Content:   e.Content,
	}
}

func summaryToItem(s summary.ProjectSummary) project.ProjectItem {
	now := time.Now()
	sentence := timeframe.ClassifyRelativeTimeframeSentence(s.PeriodStart, s.PeriodEnd, s.EntriesCount, now)
	normalizedSummary := normalizeSummary(s.SummaryText.String)

	content := fmt.Sprintf("%s %s", sentence, normalizedSummary)

	return project.ProjectItem{
		BaseDTO:           mapper.BaseToDTO(s.BaseEntity),
		ProjectID:         s.ProjectID,
		ItemType:          project.ItemTypeSummary,
		Content:           content,
		AdditionalContent: s.SummaryMarkdown.String,
		EntriesCount:      s.EntriesCount,
		EndEntryID:        s.EndEntryID,
	}
}

func normalizeSummary(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimLeft(s, ",.")
	s = strings.TrimSpace(s)

	if len(s) < 2 {
		return s
	}

	runes := []rune(s)
	// ONLY IF: First char is uppercase, Second char is lowercase
	if !unicode.IsUpper(runes[0]) || !unicode.IsLower(runes[1]) {
		return s
	}

	// Heuristic: If the next word is also capitalized, it's likely a named entity.
	words := strings.Fields(s)
	if len(words) > 1 {
		nextWord := []rune(words[1])
		if len(nextWord) > 0 && unicode.IsUpper(nextWord[0]) {
			return s
		}
	}

	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
