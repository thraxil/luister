package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/thraxil/paginate"
)

var templateDir = "templates"

type Server struct {
	DB *gorm.DB
}

type indexPage struct {
	Title        string
	TotalSongs   int
	TotalArtists int
	RecentSongs  []Song
	RecentPlays  []Play
}

func (s Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var cnt int
	s.DB.Model(&Song{}).Count(&cnt)

	var artistCnt int
	s.DB.Model(&Artist{}).Count(&artistCnt)

	var songs []Song
	s.DB.Limit(10).Order("created_at desc").Preload(
		"Artist").Preload("Album").Find(&songs)

	var plays []Play
	s.DB.Limit(10).Order("created_at desc").Preload(
		"Song").Preload("Song.Artist").Preload("Song.Album").Find(&plays)

	p := indexPage{
		Title:        "Luister",
		TotalSongs:   cnt,
		TotalArtists: artistCnt,
		RecentSongs:  songs,
		RecentPlays:  plays,
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
	s.DB.Model(&artist).Preload("Year").Order("upper(name) asc").Related(&albums)

	p := artistPage{
		Title:  artist.Name,
		Artist: artist,
		Albums: albums,
	}
	t := getTemplate("artist.html")
	t.Execute(w, p)
}

func (s Server) EditArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistID := vars["artist"]

	var artist Artist
	s.DB.First(&artist, artistID)

	newName := r.FormValue("name")

	// if there's already one with that name, merge this one into it
	var nartists []Artist
	s.DB.Model(&Artist{}).Where("name = ?", newName).Find(&nartists)
	if len(nartists) > 0 {
		nartist := nartists[0]

		// update albums
		var albums []Album
		s.DB.Model(&artist).Related(&albums)
		for _, album := range albums {
			album.Artist = nartist
			s.DB.Save(&album)
		}

		// update songs
		var songs []Song
		s.DB.Model(&artist).Related(&songs)
		for _, song := range songs {
			song.Artist = nartist
			s.DB.Save(&song)
		}

		// delete
		s.DB.Delete(&artist)

		artist = nartist
	} else {
		// otherwise, just do a simple edit and save
		artist.Name = newName
		s.DB.Save(&artist)
	}

	http.Redirect(w, r, artist.URL(), 302)
}

type artistsPage struct {
	Title   string
	Artists []Artist
	Page    paginate.Page
}

type paginatedArtists struct {
	db *gorm.DB
}

func NewPaginatedArtists(db *gorm.DB) paginatedArtists {
	return paginatedArtists{db: db}
}

func (p paginatedArtists) TotalItems() int {
	var cnt int
	p.db.Model(&Artist{}).Count(&cnt)
	return cnt
}

func (p paginatedArtists) ItemRange(offset, count int) []interface{} {
	var artists []Artist

	p.db.Model(&Artist{}).Order("upper(name) asc").Offset(offset).Limit(count).Find(&artists)

	out := make([]interface{}, len(artists))
	for j, v := range artists {
		out[j] = v
	}
	return out
}

func (s Server) ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	index := NewPaginatedArtists(s.DB)
	var p = paginate.Paginator{ItemList: index, PerPage: 100}
	page := p.GetPage(r)
	iartists := page.Items()
	artists := make([]Artist, len(iartists))
	for i, v := range iartists {
		artists[i] = v.(Artist)
	}

	ap := artistsPage{
		Title:   "all artists",
		Artists: artists,
		Page:    page,
	}
	t := getTemplate("artists.html")
	t.Execute(w, ap)
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

type randomPage struct {
	Title string
	Songs []Song
}

func (s Server) RandomHandler(w http.ResponseWriter, r *http.Request) {
	n := 10
	songs := make([]Song, n)

	for i := 0; i < n; i++ {
		s.DB.Model(&Song{}).Order("random()").Limit(1).Preload("Files").Preload("Artist").Preload("Album").Find(&songs[i])
	}

	p := randomPage{
		Title: "random playlist",
		Songs: songs,
	}
	t := getTemplate("random.html")
	t.Execute(w, p)
}

type randomSong struct {
	Title     string
	Track     int
	SongURL   string
	Artist    string
	ArtistURL string
	Album     string
	AlbumURL  string
	URL       string
	ID        string
	PlayURL   string
}

func (s Server) SingleRandomHandler(w http.ResponseWriter, r *http.Request) {
	var song Song

	s.DB.Model(&Song{}).Order("random()").Limit(1).Preload(
		"Files").Preload(
		"Artist").Preload("Album").Find(&song)

	p := randomSong{
		Title:     song.DisplayTitle(),
		Track:     song.Track,
		SongURL:   song.URL(),
		Artist:    song.Artist.DisplayName(),
		ArtistURL: song.Artist.URL(),
		Album:     song.Album.DisplayName(),
		AlbumURL:  song.Album.URL(),
		URL:       song.HakmesURL(),
		ID:        fmt.Sprintf("%d", song.ID),
		PlayURL:   song.PlayURL(),
	}

	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getTemplate(filename string) *template.Template {
	var t = template.New("base.html")
	return template.Must(t.ParseFiles(
		filepath.Join(templateDir, "base.html"),
		filepath.Join(templateDir, filename),
	))
}
