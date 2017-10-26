package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var imp = flag.Bool("import", false, "import from csv")
	var str = flag.Bool("strip", false, "strip nulls")
	flag.Parse()

	const addr = "postgresql://luister@localhost:26257/luister?sslmode=disable"
	db, err := gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	migrate(db)
	db.LogMode(true)
	if *imp {
		importcsv(db)
		return
	}
	if *str {
		stripnulls(db)
		return
	}

	s := Server{DB: db}
	r := mux.NewRouter()
	r.HandleFunc("/", s.IndexHandler)
	r.HandleFunc("/s/{song}/", s.SongHandler).Methods("GET")
	r.HandleFunc("/s/{song}/", s.EditSongHandler).Methods("POST")
	r.HandleFunc("/p/{song}/", s.PlayHandler)
	r.HandleFunc("/r/{song}/", s.RatingHandler).Methods("POST")
	r.HandleFunc("/al/{album}/", s.AlbumHandler).Methods("GET")
	r.HandleFunc("/al/{album}/", s.EditAlbumHandler).Methods("POST")
	r.HandleFunc("/ar/{artist}/", s.ArtistHandler).Methods("GET")
	r.HandleFunc("/ar/{artist}/", s.EditArtistHandler).Methods("POST")
	r.HandleFunc("/ar/", s.ArtistsHandler)
	r.HandleFunc("/search/", s.SearchHandler)
	r.HandleFunc("/random/", s.RandomHandler)

	r.HandleFunc("/api/random/", s.SingleRandomHandler)

	log.Fatal(http.ListenAndServe(":8009", r))
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Artist{})
	db.AutoMigrate(&Year{})
	db.AutoMigrate(&Album{})
	db.AutoMigrate(&Song{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Play{})
	db.AutoMigrate(&Rating{})
	db.AutoMigrate(&Tag{})

	fmt.Println("migrated")
}

func stripnulls(db *gorm.DB) {
	// song titles
	var songs []Song
	db.Find(&songs)
	for idx, song := range songs {
		trimmed := strings.Trim(song.Title, "\x00")
		if trimmed != song.Title {
			fmt.Printf("fixing song [%d]: %s (%X)\n", idx, song.Title, song.Title)
			song.UpdateTitle(db, trimmed)

		}
	}
	// albums
	var albums []Album
	db.Find(&albums)
	for idx, album := range albums {
		trimmed := strings.Trim(album.Name, "\x00")
		if trimmed != album.Name {
			fmt.Printf("fixing album [%d]: %s (%X)\n", idx, album.Name, album.Name)
			album.UpdateName(db, trimmed)
		}
	}

	// artists
	var artists []Artist
	db.Find(&artists)
	for idx, artist := range artists {
		trimmed := strings.Trim(artist.Name, "\x00")
		if trimmed != artist.Name {
			fmt.Printf("fixing artist [%d]: %s (%X)\n", idx, artist.Name, artist.Name)
			artist.UpdateName(db, trimmed)
		}
	}

}

func importcsv(db *gorm.DB) {
	// if an "import.csv" file exists, import from it
	if _, err := os.Stat("import.csv"); os.IsNotExist(err) {
		fmt.Println("no file to import")
	}
	fmt.Println("importing from import.csv...")
	file, err := os.Open("import.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		filename := record[0]
		sha1 := record[1]
		size := record[2]

		// don't care about CAP
		importFile(db, filename, sha1, size)
	}

	fmt.Println("done importing from import.csv")
}

func importFile(db *gorm.DB, filename, sha1 string, size string) {
	// skip if we already have this file
	var cnt int
	db.Model(&File{}).Where("filename = ?", filename).Or("hash = ?", sha1).Count(&cnt)
	if cnt > 0 {
		fmt.Println("already have it")
		return
	}
	ext := filepath.Ext(filename)
	// skip some obvious/common non music files
	if ext == ".nfo" || ext == ".jpg" || ext == ".jpeg" || ext == ".sfv" || ext == ".m3u" || ext == ".txt" || ext == ".ini" {
		return
	}
	// get the file from hakmes
	url := hakmesURL(sha1, ext)
	data, err := getFromHakmes(url)
	if err != nil {
		fmt.Println("failed to retrieve file from hakmes", err)
		return
	}

	// parse out id3 tags
	dreader := bytes.NewReader(data)
	fmt.Println(filename, url)
	m, err := tag.ReadFrom(dreader)
	if err != nil {
		fmt.Println("failed to read tags", err)
		return
	}
	fmt.Println("----------------")
	fmt.Println("Format", m.Format())
	fmt.Println("FileType", m.FileType())
	fmt.Println("Title", m.Title())
	fmt.Println("Album", m.Album())
	fmt.Println("Artist", m.Artist())
	fmt.Println("Year", m.Year())
	track, _ := m.Track()
	disc, _ := m.Disc()
	fmt.Println("Track", track)
	fmt.Println("Disc", disc)

	// insert record(s) into DB

	var year Year
	db.FirstOrCreate(&year, Year{Year: fmt.Sprintf("%d", m.Year())})

	var artist Artist
	db.FirstOrCreate(&artist, Artist{Name: m.Artist()})

	var album Album
	db.FirstOrCreate(&album, Album{Name: m.Album(), YearID: year.ID, ArtistID: artist.ID})

	var song Song
	db.FirstOrCreate(&song, Song{
		Title:    m.Title(),
		ArtistID: artist.ID,
		AlbumID:  album.ID,
		YearID:   year.ID,
		Track:    track,
	})

	i, err := strconv.Atoi(size)
	if err != nil {
		i = 0
	}
	var file File
	db.FirstOrCreate(&file, File{
		SongID:   song.ID,
		Filename: filename,
		Format:   string(m.Format()),
		Filetype: string(m.FileType()),
		Hash:     sha1,
		Filesize: i,
	})
}

func hakmesURL(sha1, ext string) string {
	return "http://localhost:9300/file/" + sha1 + "/file" + ext
}

func getFromHakmes(url string) ([]byte, error) {
	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.Status != "200 OK" {
		return nil, errors.New("404, probably")
	}
	b, _ := ioutil.ReadAll(resp.Body)
	return b, nil
}
