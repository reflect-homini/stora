package summary

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type Level string

const (
	DailyLevel Level = "daily"
)

type ProjectSummary struct {
	crud.BaseEntity `json:"-"`
	ProjectID       uuid.UUID      `json:"projectId"`
	SummaryMarkdown sql.NullString `json:"summaryMarkdown"`
	InsightsJSON    sql.NullString `json:"insightsJson"`
	SummaryLevel    Level          `json:"summaryLevel"`
	StartEntryID    uuid.UUID      `json:"startEntryId"`
	EndEntryID      uuid.UUID      `json:"endEntryId"`
	EntriesCount    int            `json:"entriesCount"`
	TimeframeLabel  string         `json:"timeframeLabel"`
	PeriodStart     time.Time      `json:"periodStart"`
	PeriodEnd       time.Time      `json:"periodEnd"`
	GeneratedAt     time.Time      `json:"generatedAt"`
}
