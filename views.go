package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
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
	File  File
}

func (s Server) SongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song"]

	var song Song
	s.DB.Preload("Artist").Preload("Album").Preload("Year").First(&song, songID)

	var file File
	s.DB.Where("song_id = ?", songID).First(&file)

	p := songPage{
		Title: song.Title,
		Song:  song,
		File:  file,
	}
	t := getTemplate("song.html")
	t.Execute(w, p)
}

func (s Server) PlayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID, _ := strconv.Atoi(vars["song"])

	play := Play{SongID: uint(songID)}
	s.DB.Create(&play)
	fmt.Fprintf(w, "ok")
}

type albumPage struct {
	Title string
	Album Album
	Songs []Song
}

func (s Server) AlbumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumID := vars["album"]

	var album Album
	s.DB.Preload("Artist").Preload("Year").First(&album, albumID)

	var songs []Song
	s.DB.Model(&album).Order("track asc").Preload("Files").Related(&songs)

	p := albumPage{
		Title: album.Name,
		Album: album,
		Songs: songs,
	}
	t := getTemplate("album.html")
	t.Execute(w, p)
}

type artistPage struct {
	Title  string
	Artist Artist
	Albums []Album
}

func (s Server) ArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistID := vars["artist"]

	var artist Artist
	s.DB.First(&artist, artistID)

	var albums []Album
	s.DB.Model(&artist).Preload("Year").Related(&albums)

	p := artistPage{
		Title:  artist.Name,
		Artist: artist,
		Albums: albums,
	}
	t := getTemplate("artist.html")
	t.Execute(w, p)
}

type searchPage struct {
	Title   string
	Query   string
	Artists []Artist
	Albums  []Album
	Songs   []Song
}

func (s Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	// TODO: strip/normalize

	if query == "" {
		http.Redirect(w, r, "/", 302)
		return
	}

	var artists []Artist
	s.DB.Where("name ILIKE ?", "%"+query+"%").Order("upper(name) asc").Find(&artists)

	var albums []Album
	s.DB.Where("name ILIKE ?", "%"+query+"%").Order("upper(name) asc").Preload("Artist").Preload("Year").Find(&albums)

	var songs []Song
	s.DB.Where("title ILIKE ?", "%"+query+"%").Order("upper(title) asc").Preload("Artist").Preload("Album").Find(&songs)

	p := searchPage{
		Title:   "search results for '" + query + "'",
		Query:   query,
		Artists: artists,
		Albums:  albums,
		Songs:   songs,
	}
	t := getTemplate("search.html")
	t.Execute(w, p)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(templateDir, "base.html"),
		filepath.Join(templateDir, filename),
	))
}
