package main

import (
	"fmt"
	"html/template"
	"net/http"
	"personal-blog/internal/auth"
	"personal-blog/internal/blog"
	"strings"
	"time"
)

var tmpls = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/article/", handleArticleView)

	http.HandleFunc("/admin", auth.BasicAuth(handleDashboard))
	http.HandleFunc("/admin/new", auth.BasicAuth(handleNewArticle))
	http.HandleFunc("/admin/edit/", auth.BasicAuth(handleEditArticle))
	http.HandleFunc("/admin/delete/", auth.BasicAuth(handleDeleteArticle))

	fmt.Println("🚀 Blog server started at http://localhost:8080")
	fmt.Println("Admin credentials: admin / password123")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	articles, _ := blog.GetAllArticles()
	tmpls.ExecuteTemplate(w, "home.html", articles)
}

func handleArticleView(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/article/")
	articles, _ := blog.GetAllArticles()

	var found blog.Article
	for _, a := range articles {
		if a.Slug == slug {
			found = a
			break
		}
	}

	if found.Slug == "" {
		http.NotFound(w, r)
		return
	}
	tmpls.ExecuteTemplate(w, "article.html", found)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	articles, _ := blog.GetAllArticles()
	tmpls.ExecuteTemplate(w, "dashboard.html", articles)
}

func handleNewArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		newArt := blog.Article{
			Title:     r.FormValue("title"),
			Content:   r.FormValue("content"),
			Published: time.Now().Format("2006-01-02"),
		}
		newArt.Slug = strings.ToLower(strings.ReplaceAll(newArt.Title, " ", "-"))

		blog.SaveArticle(newArt)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	tmpls.ExecuteTemplate(w, "add.html", nil)
}

func handleEditArticle(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/admin/edit/")

	if r.Method == http.MethodPost {
		updated := blog.Article{
			Slug:      slug,
			Title:     r.FormValue("title"),
			Content:   r.FormValue("content"),
			Published: r.FormValue("date"),
		}
		blog.SaveArticle(updated)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	articles, _ := blog.GetAllArticles()
	for _, a := range articles {
		if a.Slug == slug {
			tmpls.ExecuteTemplate(w, "edit.html", a)
			return
		}
	}
}

func handleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/admin/delete/")
	err := blog.DeleteArticle(slug)
	if err != nil {
		http.Error(w, "Failed to delete", 500)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
