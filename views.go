package main

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var templateDir = "templates"

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
	t := getTemplate("index.html")
	t.Execute(w, p)
}

type songPage struct {
	Title string
	Song  Song
}

func (s Server) SongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song"]

	var song Song
	s.DB.Preload("Artist").Preload("Album").Preload("Year").First(&song, songID)

	p := songPage{
		Title: song.Title,
		Song:  song,
	}
	t := getTemplate("song.html")
	t.Execute(w, p)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(templateDir, "base.html"),
		filepath.Join(templateDir, filename),
	))
}
