package notes

import (
	"os"
	"path/filepath"
	"strings"
)

type FileRepository struct {
	storageDir string
}

func NewFileRepository(dir string) (*FileRepository, error) {

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &FileRepository{storageDir: dir}, nil
}

func (r *FileRepository) Save(filename string, content []byte) error {
	if !strings.HasSuffix(filename, ".md") {
		filename = filename + ".md"
	}
	targetPath := filepath.Join(r.storageDir, filename)
	return os.WriteFile(targetPath, content, 0644)
}

func (r *FileRepository) List() ([]NoteMetadata, error) {
	entries, err := os.ReadDir(r.storageDir)
	if err != nil {
		return nil, err
	}

	var list []NoteMetadata
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			list = append(list, NoteMetadata{
				Filename:  entry.Name(),
				Size:      info.Size(),
				UpdatedAt: info.ModTime(),
			})
		}
	}
	return list, nil
}

func (r *FileRepository) ReadContent(filename string) ([]byte, error) {
	targetPath := filepath.Join(r.storageDir, filename)
	return os.ReadFile(targetPath)
}
