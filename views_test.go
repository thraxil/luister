package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
