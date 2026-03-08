package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestRecentlyPlayedAPIHandler(t *testing.T) {
	db := SetupTestDB(t)
	artist := Artist{Name: "Test Artist"}
	db.Create(&artist)
	album := Album{Name: "Test Album", ArtistID: artist.ID}
	db.Create(&album)
	song := Song{Title: "Test Song", ArtistID: artist.ID, AlbumID: album.ID}
	db.Create(&song)
	file := File{Filename: "test.mp3", SongID: song.ID}
	db.Create(&file)
	play := Play{SongID: song.ID}
	db.Create(&play)

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/api/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RecentlyPlayedAPIHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Plays []randomSong
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response.Plays) != 1 {
		t.Errorf("expected 1 play, got %d", len(response.Plays))
	}
	if response.Plays[0].Title != "Test Song" {
		t.Errorf("expected title 'Test Song', got '%s'", response.Plays[0].Title)
	}
}

func TestPlayHandler(t *testing.T) {
	db := SetupTestDB(t)
	song := Song{Title: "Test Song"}
	db.Create(&song)

	server := Server{DB: db}

	req, err := http.NewRequest("POST", "/api/play/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"song": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.PlayHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var count int64
	db.Model(&Play{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 play created, got %d", count)
	}
}

func TestRatingHandler(t *testing.T) {
	db := SetupTestDB(t)
	song := Song{Title: "Test Song", Rating: 0}
	db.Create(&song)

	server := Server{DB: db}

	req, err := http.NewRequest("POST", "/api/rate/1", strings.NewReader("rating=5"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = mux.SetURLVars(req, map[string]string{"song": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RatingHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedSong Song
	db.First(&updatedSong, song.ID)
	if updatedSong.Rating != 5 {
		t.Errorf("expected rating 5, got %d", updatedSong.Rating)
	}
}

func TestTagAPIHandler(t *testing.T) {
	db := SetupTestDB(t)
	tag := Tag{Name: "rock"}
	db.Create(&tag)
	artist := Artist{Name: "Test Artist"}
	db.Create(&artist)
	album := Album{Name: "Test Album", ArtistID: artist.ID}
	db.Create(&album)
	song := Song{Title: "Test Song", ArtistID: artist.ID, AlbumID: album.ID}
	db.Create(&song)
	file := File{Filename: "test.mp3", SongID: song.ID}
	db.Create(&file)
	if err := db.Model(&song).Association("Tags").Append(&tag); err != nil {
		t.Fatalf("failed to append tag: %v", err)
	}

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/api/tag/rock", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"tag": "rock"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.TagAPIHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Tag   Tag
		Songs []randomSong
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Tag.Name != "rock" {
		t.Errorf("expected tag name 'rock', got '%s'", response.Tag.Name)
	}
	if len(response.Songs) != 1 {
		t.Errorf("expected 1 song, got %d", len(response.Songs))
	}
}

func TestTagsAPIHandler(t *testing.T) {
	db := SetupTestDB(t)
	tag := Tag{Name: "rock"}
	db.Create(&tag)
	tag2 := Tag{Name: "pop"}
	db.Create(&tag2)

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/api/tags", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.TagsAPIHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Tags []Tag
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(response.Tags))
	}
}

func TestSingleRandomHandler(t *testing.T) {
	db := SetupTestDB(t)
	artist := Artist{Name: "Test Artist"}
	db.Create(&artist)
	album := Album{Name: "Test Album", ArtistID: artist.ID}
	db.Create(&album)
	song := Song{Title: "Test Song", ArtistID: artist.ID, AlbumID: album.ID}
	db.Create(&song)
	file := File{Filename: "test.mp3", SongID: song.ID}
	db.Create(&file)

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/api/random", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.SingleRandomHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response randomSong
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Title != "Test Song" {
		t.Errorf("expected title 'Test Song', got '%s'", response.Title)
	}
}

func TestRandomPlaylistHandler(t *testing.T) {
	db := SetupTestDB(t)
	artist := Artist{Name: "Test Artist"}
	db.Create(&artist)
	album := Album{Name: "Test Album", ArtistID: artist.ID}
	db.Create(&album)
	song := Song{Title: "Test Song", ArtistID: artist.ID, AlbumID: album.ID}
	db.Create(&song)
	file := File{Filename: "test.mp3", SongID: song.ID}
	db.Create(&file)

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/api/random/playlist", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RandomPlaylistHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response struct {
		Songs []randomSong
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response.Songs) != 1 {
		t.Errorf("expected 1 song, got %d", len(response.Songs))
	}
	if response.Songs[0].Title != "Test Song" {
		t.Errorf("expected title 'Test Song', got '%s'", response.Songs[0].Title)
	}
}
