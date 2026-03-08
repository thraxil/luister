package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockStore struct {
	CountSongsFunc        func() (int64, error)
	CountUnratedSongsFunc func() (int64, error)
	CountArtistsFunc      func() (int64, error)
	GetRecentSongsFunc    func(limit int) ([]Song, error)
	GetRecentPlaysFunc    func(limit int) ([]Play, error)
}

func (m *MockStore) CountSongs() (int64, error) {
	if m.CountSongsFunc != nil {
		return m.CountSongsFunc()
	}
	return 0, nil
}

func (m *MockStore) CountUnratedSongs() (int64, error) {
	if m.CountUnratedSongsFunc != nil {
		return m.CountUnratedSongsFunc()
	}
	return 0, nil
}

func (m *MockStore) CountArtists() (int64, error) {
	if m.CountArtistsFunc != nil {
		return m.CountArtistsFunc()
	}
	return 0, nil
}

func (m *MockStore) GetRecentSongs(limit int) ([]Song, error) {
	if m.GetRecentSongsFunc != nil {
		return m.GetRecentSongsFunc(limit)
	}
	return []Song{}, nil
}

func (m *MockStore) GetRecentPlays(limit int) ([]Play, error) {
	if m.GetRecentPlaysFunc != nil {
		return m.GetRecentPlaysFunc(limit)
	}
	return []Play{}, nil
}

func TestIndexHandler(t *testing.T) {
	mockStore := &MockStore{
		CountSongsFunc: func() (int64, error) {
			return 100, nil
		},
		CountUnratedSongsFunc: func() (int64, error) {
			return 10, nil
		},
		CountArtistsFunc: func() (int64, error) {
			return 5, nil
		},
		GetRecentSongsFunc: func(limit int) ([]Song, error) {
			return []Song{{Title: "Test Song"}}, nil
		},
		GetRecentPlaysFunc: func(limit int) ([]Play, error) {
			return []Play{{Song: Song{Title: "Played Song"}}}, nil
		},
	}

	server := Server{Store: mockStore}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.IndexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Luister"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want substring %v",
			rr.Body.String(), expected)
	}

	if !strings.Contains(rr.Body.String(), "Test Song") {
		t.Errorf("handler returned unexpected body: missing recent song")
	}
}

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&Song{}, &Artist{}, &Album{}, &Tag{}, &Play{}, &File{}, &Year{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func TestSongHandler(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/song/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"song": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.SongHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "Test Song") {
		t.Errorf("handler returned unexpected body: missing song title")
	}
}

func TestAlbumHandler(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/album/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"album": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.AlbumHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "Test Album") {
		t.Errorf("handler returned unexpected body: missing album name")
	}
	if !strings.Contains(rr.Body.String(), "Test Song") {
		t.Errorf("handler returned unexpected body: missing song title")
	}
}

func TestArtistHandler(t *testing.T) {
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

	req, err := http.NewRequest("GET", "/artist/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"artist": "1"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.ArtistHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "Test Artist") {
		t.Errorf("handler returned unexpected body: missing artist name")
	}
	if !strings.Contains(rr.Body.String(), "Test Song") {
		t.Errorf("handler returned unexpected body: missing song title")
	}
}

func TestTagHandler(t *testing.T) {
	db := SetupTestDB(t)
	tag := Tag{Name: "rock"}
	db.Create(&tag)
	song := Song{Title: "Test Song"}
	db.Create(&song)
	file := File{Filename: "test.mp3", SongID: song.ID}
	db.Create(&file)
	if err := db.Model(&song).Association("Tags").Append(&tag); err != nil {
		t.Fatalf("failed to append tag: %v", err)
	}

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/tag/rock", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"tag": "rock"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.TagHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "rock") {
		t.Errorf("handler returned unexpected body: missing tag name")
	}
	if !strings.Contains(rr.Body.String(), "Test Song") {
		t.Errorf("handler returned unexpected body: missing song title")
	}
}

func TestTagsHandler(t *testing.T) {
	db := SetupTestDB(t)
	tag := Tag{Name: "rock"}
	db.Create(&tag)
	tag2 := Tag{Name: "pop"}
	db.Create(&tag2)

	server := Server{DB: db}

	req, err := http.NewRequest("GET", "/tags", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.TagsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "rock") {
		t.Errorf("handler returned unexpected body: missing tag name")
	}
	if !strings.Contains(rr.Body.String(), "pop") {
		t.Errorf("handler returned unexpected body: missing tag name")
	}
}
