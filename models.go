package main

import "github.com/jinzhu/gorm"

type Artist struct {
	gorm.Model
	Name string `gorm:"unique_index"`

	Albums []Album
	Songs  []Song
}

type Album struct {
	gorm.Model
	Name     string
	ArtistID uint `gorm:"index"`
	YearID   int  `gorm:"index"`
	Artist   Artist
	Year     Year
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
