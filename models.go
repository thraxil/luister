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
