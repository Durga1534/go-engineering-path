package notes

import "time"

type NoteMetadata struct {
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GrammarReport struct {
	OriginalText string   `json:"original_text"`
	IssueCount   int      `json:"issue_count"`
	Suggestions  []string `json:"suggestions"`
}
