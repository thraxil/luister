package main

import (
	"fmt"
	"path/filepath"

	"github.com/jinzhu/gorm"
)

type Artist struct {
	gorm.Model
	Name string `gorm:"unique_index"`

	Albums []Album
	Songs  []Song
}

func (a Artist) URL() string {
	return fmt.Sprintf("/ar/%d/", a.ID)
}

func (a Artist) DisplayName() string {
	if a.Name == "" {
		return "[missing]"
	}
	return a.Name
}

func (a Artist) UpdateName(db *gorm.DB, newName string) Artist {
	// if there's already one with that name, merge this one into it
	var nartists []Artist
	db.Model(&Artist{}).Where("name = ?", newName).Find(&nartists)
	if len(nartists) > 0 {
		nartist := nartists[0]

		// update albums
		var albums []Album
		db.Model(&a).Related(&albums)
		for _, album := range albums {
			album.Artist = nartist
			db.Save(&album)
		}

		// update songs
		var songs []Song
		db.Model(&a).Related(&songs)
		for _, song := range songs {
			song.Artist = nartist
			db.Save(&song)
		}

		// delete
		db.Delete(&a)

		return nartist
	}
	// otherwise, just do a simple edit and save
	a.Name = newName
	db.Save(&a)
	return a
}

type Album struct {
	gorm.Model
	Name     string
	ArtistID uint `gorm:"index"`
	YearID   int  `gorm:"index"`
	Artist   Artist
	Year     Year

	Songs []Song
}

func (a Album) URL() string {
	return fmt.Sprintf("/al/%d/", a.ID)
}

func (a Album) DisplayName() string {
	if a.Name == "" {
		return "[missing]"
	}
	return a.Name
}

func (a Album) UpdateName(db *gorm.DB, newName string) Album {
	// if there's already one with that name, merge this one into it
	var nalbums []Album
	db.Model(&Album{}).Where("name = ?", newName).Where("artist_id = ?", a.Artist.ID).Find(&nalbums)
	if len(nalbums) > 0 {
		nalbum := nalbums[0]

		// update songs
		var songs []Song
		db.Model(&a).Related(&songs)
		for _, song := range songs {
			song.Album = nalbum
			db.Save(&song)
		}

		// delete
		db.Delete(&a)

		return nalbum
	}
	// otherwise, just do a simple edit and save
	a.Name = newName
	db.Save(&a)
	return a
}

type Year struct {
	ID   int
	Year string

	Albums []Album
	Songs  []Song
}

type Song struct {
	gorm.Model
	Title    string
	ArtistID uint `gorm:"index"`
	AlbumID  uint `gorm:"index"`
	YearID   int  `gorm:"index"`
	Track    int

	Plays []Play
	Tags  []Tag `gorm:"many2many:song_tags"`

	Artist Artist
	Album  Album
	Year   Year
	Files  []File
}

func (s Song) HakmesURL() string {
	// assumes that Files has been preloaded
	return s.Files[0].HakmesURL()
}

func (s Song) URL() string {
	return fmt.Sprintf("/s/%d/", s.ID)
}

func (s Song) PlayURL() string {
	return fmt.Sprintf("/p/%d/", s.ID)
}

func (s Song) DisplayTitle() string {
	if s.Title == "" {
		return "[missing]"
	}
	return s.Title
}

type File struct {
	gorm.Model
	SongID   uint   `gorm:"index"`
	Filename string `gorm:"index"`
	Format   string `gorm:"index"`
	Filetype string
	Hash     string
	Filesize int

	Song Song
}

func (f File) HakmesURL() string {
	ext := filepath.Ext(f.Filename)
	return "http://localhost:9300/file/" + f.Hash + "/file" + ext
}

type Play struct {
	gorm.Model
	SongID uint `gorm:"index"`
	Song   Song
}

type Rating struct {
	gorm.Model
	Rating int
	SongID uint `gorm:"index"`
	Song   Song
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique_index"`
}
