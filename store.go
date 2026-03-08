package main

import "gorm.io/gorm"

type Store interface {
	CountSongs() (int64, error)
	CountUnratedSongs() (int64, error)
	CountArtists() (int64, error)
	GetRecentSongs(limit int) ([]Song, error)
	GetRecentPlays(limit int) ([]Play, error)
}

type DBStore struct {
	DB *gorm.DB
}

func (s *DBStore) CountSongs() (int64, error) {
	var cnt int64
	err := s.DB.Model(&Song{}).Count(&cnt).Error
	return cnt, err
}

func (s *DBStore) CountUnratedSongs() (int64, error) {
	var cnt int64
	err := s.DB.Model(&Song{}).Where("rating = 0").Count(&cnt).Error
	return cnt, err
}

func (s *DBStore) CountArtists() (int64, error) {
	var cnt int64
	err := s.DB.Model(&Artist{}).Count(&cnt).Error
	return cnt, err
}

func (s *DBStore) GetRecentSongs(limit int) ([]Song, error) {
	var songs []Song
	err := s.DB.Limit(limit).Order("created_at desc").Preload(
		"Artist").Preload("Album").Find(&songs).Error
	return songs, err
}

func (s *DBStore) GetRecentPlays(limit int) ([]Play, error) {
	var plays []Play
	err := s.DB.Limit(limit).Order("created_at desc").Preload(
		"Song").Preload("Song.Artist").Preload("Song.Album").Find(&plays).Error
	return plays, err
}
