package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s Server) RecentlyPlayedAPIHandler(w http.ResponseWriter, r *http.Request) {
	n := 10
	var plays []Play
	s.DB.Limit(n).Order("created_at desc").Preload(
		"Song").Preload("Song.Artist").Preload("Song.Album").Preload("Song.Files").Find(&plays)

	randomSongs := make([]randomSong, n)
	for i, play := range plays {
		randomSongs[i] = randomSong{
			Title:     play.Song.DisplayTitle(),
			Track:     play.Song.Track,
			SongURL:   play.Song.URL(),
			Artist:    play.Song.Artist.DisplayName(),
			ArtistID:  fmt.Sprintf("%d", play.Song.Artist.ID),
			ArtistURL: play.Song.Artist.URL(),
			Album:     play.Song.Album.DisplayName(),
			AlbumID:   fmt.Sprintf("%d", play.Song.Album.ID),
			AlbumURL:  play.Song.Album.URL(),
			URL:       play.Song.HakmesURL(),
			ID:        fmt.Sprintf("%d", play.Song.ID),
			PlayURL:   play.Song.PlayURL(),
			Rating:    play.Song.Rating,
		}
	}

	p := struct{ Plays []randomSong }{
		Plays: randomSongs,
	}

	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
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

type randomSong struct {
	Title     string
	Track     int
	SongURL   string
	Artist    string
	ArtistID  string
	ArtistURL string
	Album     string
	AlbumID   string
	AlbumURL  string
	URL       string
	ID        string
	PlayURL   string
	Rating    int
}

func (s Server) NRandomSongs(n int) []randomSong {
	var songs []Song

	s.DB.Model(&Song{}).Order("random()").Limit(n).Preload(
		"Files").Preload(
		"Artist").Preload("Album").Find(&songs)

	randomSongs := make([]randomSong, n)
	for i, song := range songs {
		randomSongs[i] = randomSong{
			Title:     song.DisplayTitle(),
			Track:     song.Track,
			SongURL:   song.URL(),
			Artist:    song.Artist.DisplayName(),
			ArtistID:  fmt.Sprintf("%d", song.Artist.ID),
			ArtistURL: song.Artist.URL(),
			Album:     song.Album.DisplayName(),
			AlbumID:   fmt.Sprintf("%d", song.Album.ID),
			AlbumURL:  song.Album.URL(),
			URL:       song.HakmesURL(),
			ID:        fmt.Sprintf("%d", song.ID),
			PlayURL:   song.PlayURL(),
			Rating:    song.Rating,
		}
	}
	return randomSongs
}

func (s Server) SingleRandomHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(s.NRandomSongs(1)[0])
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s Server) RandomPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	songs := s.NRandomSongs(5)
	p := struct{ Songs []randomSong }{
		Songs: songs,
	}
	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
