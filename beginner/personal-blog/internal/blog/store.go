package blog

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Article struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published string `json:"published"`
}

const dataDir = "data"

func SaveArticle(a Article) error {
	os.MkdirAll(dataDir, 0755)
	filePath := filepath.Join(dataDir, a.Slug+".json")
	data, _ := json.MarshalIndent(a, "", "")
	return os.WriteFile(filePath, data, 0644)
}

func GetAllArticles() ([]Article, error) {
	files, _ := os.ReadDir(dataDir)
	var articles []Article
	for _, f := range files {
		data, _ := os.ReadFile(filepath.Join(dataDir, f.Name()))
		var a Article
		json.Unmarshal(data, &a)
		articles = append(articles, a)
	}

	return articles, nil
}

func DeleteArticle(slug string) error {
	return os.Remove(filepath.Join("data", slug+".json"))
}
