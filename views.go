package main

import (
	"net/http"
	"text/template"

	"github.com/jinzhu/gorm"
)

type Server struct {
	DB *gorm.DB
}

type indexPage struct {
	Title       string
	TotalSongs  int
	RecentSongs []Song
}

func (s Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var cnt int
	s.DB.Model(&Song{}).Count(&cnt)

	var songs []Song
	s.DB.Limit(10).Order("created_at desc").Preload(
		"Artist").Preload("Album").Find(&songs)

	p := indexPage{
		Title:       "Luister",
		TotalSongs:  cnt,
		RecentSongs: songs,
	}
	t, _ := template.New("index").Parse(indexTemplate)
	t.Execute(w, p)
}
