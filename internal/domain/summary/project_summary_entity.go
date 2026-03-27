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
	crud.BaseEntity
	ProjectID       uuid.UUID
	SummaryText     sql.NullString
	SummaryMarkdown sql.NullString
	InsightsJSON    sql.NullString
	SummaryLevel    Level
	StartEntryID    uuid.UUID
	EndEntryID      uuid.UUID
	EntriesCount    int
	PeriodStart     time.Time
	PeriodEnd       time.Time
	GeneratedAt     time.Time
}
