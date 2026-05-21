package notes

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
)

type Service struct {
	repo *FileRepository
}

func NewService(repo *FileRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RenderToHTML(filename string) (string, error) {
	rawBytes, err := s.repo.ReadContent(filename)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := goldmark.Convert(rawBytes, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Service) CheckGrammar(text string) GrammarReport {
	var suggestions []string
	lowerText := strings.ToLower(text)

	if strings.Contains(lowerText, "i am variables") {
		suggestions = append(suggestions, "Change 'i am variables' to 'variables are'")
	}
	if strings.Contains(lowerText, "dont") {
		suggestions = append(suggestions, "Replace 'dont' with 'proper apostrophe use: 'don't'")
	}
	return GrammarReport{
		OriginalText: text,
		IssueCount:   len(suggestions),
		Suggestions:  suggestions,
	}
}

func (s *Service) SaveNote(filename string, content []byte) error {
	return s.repo.Save(filename, content)
}

func (s *Service) ListNotes() ([]NoteMetadata, error) {
	return s.repo.List()
}
