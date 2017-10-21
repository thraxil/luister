package main

import "github.com/jinzhu/gorm"

type Artist struct {
	gorm.Model
	Name string `gorm:"unique_index"`

	Albums []Album `gorm:"many2many:album_artists;"`
	Songs  []Song
}

type Album struct {
	gorm.Model
	Name string

	Artists []Artist `gorm:"many2many:album_artists;"`
	YearID  int      `gorm:"index"`
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
	ArtistID int `gorm:"index"`
	AlbumID  int `gorm:"index"`
	YearID   int `gorm:"index"`
	Track    int

	Plays []Play
	Tags  []Tag `gorm:"many2many:song_tags"`
}

type File struct {
	gorm.Model
	SongID   int `gorm:"index"`
	Filename string
	Format   string
	Filetype string
	Hash     string
	Bitrate  int
	Filesize int
}

type Play struct {
	gorm.Model
	SongID int `gorm:"index"`
}

type Rating struct {
	gorm.Model
	Rating int
	SongID int `gorm:"index"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique_index"`
}
