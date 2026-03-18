package projectdetails

import (
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/mapper"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/summary"
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
	return project.ProjectItem{
		BaseDTO:      mapper.BaseToDTO(s.BaseEntity),
		ProjectID:    s.ProjectID,
		ItemType:     project.ItemTypeSummary,
		Content:      s.SummaryMarkdown.String,
		EntriesCount: s.EntriesCount,
		EndEntryID:   s.EndEntryID,
	}
}
