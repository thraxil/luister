package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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
	UnratedCnt   int
	TotalArtists int
	RecentSongs  []Song
	RecentPlays  []Play
}

func (s Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var cnt int
	s.DB.Model(&Song{}).Count(&cnt)

	var unratedCnt int
	s.DB.Model(&Song{}).Where("rating = 0").Count(&unratedCnt)

	var artistCnt int
	s.DB.Model(&Artist{}).Count(&artistCnt)

	var songs []Song
	s.DB.Limit(10).Order("created_at desc").Preload(
		"Artist").Preload("Album").Find(&songs)

	var plays []Play
	s.DB.Limit(25).Order("created_at desc").Preload(
		"Song").Preload("Song.Artist").Preload("Song.Album").Find(&plays)

	p := indexPage{
		Title:        "Luister",
		TotalSongs:   cnt,
		UnratedCnt:   unratedCnt,
		TotalArtists: artistCnt,
		RecentSongs:  songs,
		RecentPlays:  plays,
	}
	t := getTemplate("index.html")
	t.Execute(w, p)
}

func (s Server) RecentlyPlayedAPIHandler(w http.ResponseWriter, r *http.Request) {
	var plays []Play
	s.DB.Limit(25).Order("created_at desc").Preload(
		"Song").Preload("Song.Artist").Preload("Song.Album").Find(&plays)

	p := struct{ Plays []Play }{
		Plays: plays,
	}
	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s Server) VueHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/dist/index.html")
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
	s.DB.Preload("Artist").Preload("Album").Preload("Year").Preload("Tags").First(&song, songID)

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

func (s Server) EditSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song"]

	var song Song
	s.DB.First(&song, songID)

	newName := r.FormValue("title")

	song = song.UpdateTitle(s.DB, newName)
	http.Redirect(w, r, song.URL(), 302)
}

func (s Server) TagSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID := vars["song"]

	var song Song
	s.DB.First(&song, songID).Preload("Tags")

	tags := strings.Split(r.FormValue("tags"), ",")

	// delete existing
	for _, t := range song.Tags {
		s.DB.Model(&song).Association("Tags").Delete(&t)
	}

	for _, t := range tags {
		t = strings.Trim(t, " \n\r\t,\"'!?")
		var tag Tag
		s.DB.FirstOrCreate(&tag, Tag{Name: t})
		s.DB.Model(&song).Association("Tags").Append(tag)
	}

	http.Redirect(w, r, song.URL(), 302)
}

func (s Server) PlayHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID, _ := strconv.Atoi(vars["song"])

	play := Play{SongID: uint(songID)}
	s.DB.Create(&play)
	fmt.Fprintf(w, "ok")
}

func (s Server) RatingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songID, _ := strconv.Atoi(vars["song"])

	var song Song
	s.DB.First(&song, uint(songID))

	submitted, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		return
	}
	song.Rating = submitted
	s.DB.Save(&song)
	fmt.Fprintf(w, "ok")
}

type tagPage struct {
	Title string
	Tag   Tag
	Songs []Song
}

func (s Server) TagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagName := vars["tag"]

	var tag Tag
	s.DB.Model(&Tag{}).Where("name = ?", tagName).First(&tag)

	var songs []Song
	s.DB.Table("songs").
		Joins("JOIN song_tags on song_tags.tag_id=? AND song_tags.song_id=songs.id", tag.ID).
		Preload("Artist").
		Preload("Album").
		Preload("Files").
		Find(&songs)

	p := tagPage{
		Title: tag.Name,
		Tag:   tag,
		Songs: songs,
	}
	t := getTemplate("tag.html")
	t.Execute(w, p)
}

func (s Server) TagAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagName := vars["tag"]

	var tag Tag
	s.DB.Model(&Tag{}).Where("name = ?", tagName).First(&tag)

	var songs []Song
	s.DB.Table("songs").
		Joins("JOIN song_tags on song_tags.tag_id=? AND song_tags.song_id=songs.id", tag.ID).
		Preload("Artist").
		Preload("Album").
		Preload("Files").
		Find(&songs)

	p := struct {
		Tag   Tag
		Songs []Song
	}{
		Tag:   tag,
		Songs: songs,
	}
	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

type tagsPage struct {
	Title string
	Tags  []Tag
}

func (s Server) TagsHandler(w http.ResponseWriter, r *http.Request) {
	var tags []Tag
	s.DB.Order("upper(name) asc").Find(&tags)

	p := tagsPage{
		Title: "Tags",
		Tags:  tags,
	}
	t := getTemplate("tags.html")
	t.Execute(w, p)
}

func (s Server) TagsAPIHandler(w http.ResponseWriter, r *http.Request) {
	var tags []Tag
	s.DB.Order("upper(name) asc").Find(&tags)

	p := struct{ Tags []Tag }{
		Tags: tags,
	}
	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
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
	s.DB.Preload("Artist").First(&album, albumID)

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

func (s Server) EditAlbumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumID := vars["album"]

	var album Album
	s.DB.Preload("Artist").First(&album, albumID)

	newName := r.FormValue("name")

	album = album.UpdateName(s.DB, newName)
	http.Redirect(w, r, album.URL(), 302)
}

type artistPage struct {
	Title  string
	Artist Artist
	Albums []Album
	Songs  []Song
}

func (s Server) ArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistID := vars["artist"]

	var artist Artist
	s.DB.First(&artist, artistID)

	var albums []Album
	s.DB.Model(&artist).Preload("Year").Order("upper(name) asc").Related(&albums)

	var songs []Song
	s.DB.Where("artist_id = ?", artistID).Order("rating desc, album_id asc, track asc").Preload("Album").Find(&songs)

	p := artistPage{
		Title:  artist.Name,
		Artist: artist,
		Albums: albums,
		Songs:  songs,
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

	artist = artist.UpdateName(s.DB, newName)
	http.Redirect(w, r, artist.URL(), 302)
}

type artistsPage struct {
	Title   string
	Artists []Artist
	Page    paginate.Page
}

type paginatedArtists struct {
	db  *gorm.DB
	cnt int
}

func NewPaginatedArtists(db *gorm.DB, cnt int) paginatedArtists {
	return paginatedArtists{db: db, cnt: cnt}
}

func (p paginatedArtists) TotalItems() int {
	return p.cnt
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
	var cnt int
	s.DB.Model(&Artist{}).Count(&cnt)
	index := NewPaginatedArtists(s.DB, cnt)
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
	n := 25

	var songs []Song
	s.DB.Model(&Song{}).Order("random()").Limit(n).Preload("Files").Preload("Artist").Preload("Album").Find(&songs)

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
	Rating    int
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
		Rating:    song.Rating,
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
