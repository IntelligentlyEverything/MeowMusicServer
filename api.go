package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Api Response.
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Tips string      `json:"tips"`
	Ip   string      `json:"ip"`
}

// API Song response.
type Song struct {
	Num      int         `json:"num"`
	Song     string      `json:"song"`
	Singer   string      `json:"singer"`
	Album    string      `json:"album"`
	Cover    string      `json:"cover"`
	MusicURL interface{} `json:"music_url"`
	Lyric    interface{} `json:"lyric"`
}

type MusicURL struct {
	Audition     string `json:"audition"`
	Standard     string `json:"standard"`
	Highquality  string `json:"highquality"`
	Superquality string `json:"superquality"`
	Lossless     string `json:"lossless"`
	Hires        string `json:"hires"`
}

type Lyric struct {
	Mrc string `json:"mrc"`
	Lrc string `json:"lrc"`
}

// apiHandler is the handler function for API requests.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicServer")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	queryParams := r.URL.Query()
	//key := queryParams.Get("key")
	msg := queryParams.Get("msg")
	numStr := queryParams.Get("num")

	ip, err := IPhandler(r)
	if err != nil {
		ip = "0.0.0.0"
	}

	if msg == "" {
		response := Response{
			Code: 1,
			Msg:  "API Operation failed: 'msg' parameter is required.",
			Data: []interface{}{},
			Tips: "Provide by " + os.Getenv("WEBSITE_NAME"),
			Ip:   ip,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert numStr to int if provided
	var num int
	if numStr != "" {
		num, err = strconv.Atoi(numStr)
		if err != nil {
			response := Response{
				Code: 2,
				Msg:  "Invalid 'num' parameter provided.",
				Data: []interface{}{},
				Tips: "Provide by " + os.Getenv("WEBSITE_NAME"),
				Ip:   ip,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Get songs based on msg and num
	songs := apiSongHandlerOnMetadata(msg, num)

	// If no songs found, return an empty array
	if len(songs) == 0 {
		response := Response{
			Code: 3,
			Msg:  "No songs found for the given query.",
			Data: []interface{}{},
			Tips: "Provide by " + os.Getenv("WEBSITE_NAME"),
			Ip:   ip,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Prepare the response
	response := Response{
		Code: 0,
		Msg:  "API Operation successful.",
		Data: songs,
		Tips: "Provide by " + os.Getenv("WEBSITE_NAME"),
		Ip:   ip,
	}

	// Encode and send the response
	json.NewEncoder(w).Encode(response)
}

type API struct {
	APIURL  string `json:"api_url"`
	APIType string `json:"api_type"`
}

type OtherSong struct {
	SongName string   `json:"song_name"`
	Singer   string   `json:"singer"`
	Album    string   `json:"album"`
	Cover    string   `json:"cover"`
	MusicURL MusicURL `json:"music_url"`
	LyricURL Lyric    `json:"lyric_url"`
}

type Metadata struct {
	API   []API       `json:"api"`
	Other []OtherSong `json:"other"`
}

// API Song handler.
func apiSongHandlerOnMetadata(msg string, num int) []Song {
	// Construct a complete file path for metadata.json
	metadataFilePath := filepath.Join("./music-uploads", "metadata.json")

	// Read the metadata file
	metadata, err := readMetadataFile(metadataFilePath)
	if err != nil {
		fmt.Println("Error reading metadata file: ", err)
		return []Song{}
	}

	// Get the song array from metadata
	otherSongs := getSongArray(metadata)

	// Scan the music-uploads directory for artist-song folders
	artistSongFolders := scanArtistSongFolders("./music-uploads")

	// Initialize a counter for songs
	songCounter := 0
	var filteredSongs []Song

	// Convert artist-song folders to Song
	for _, artistSongFolder := range artistSongFolders {
		songCounter++
		song := Song{
			Num:      songCounter,
			Song:     artistSongFolder.songName,
			Singer:   artistSongFolder.artistName,
			Album:    artistSongFolder.albumName,
			Cover:    getCoverURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
			MusicURL: getMusicURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
			Lyric:    getLyricURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
		}

		// Check if the song matches the msg and num
		if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
			if num == 0 || song.Num == num {
				filteredSongs = append(filteredSongs, song)
			}
		}
	}

	// Convert OtherSong to Song
	for _, otherSong := range otherSongs {
		songCounter++
		song := Song{
			Num:      songCounter,
			Song:     otherSong.SongName,
			Singer:   otherSong.Singer,
			Album:    otherSong.Album,
			Cover:    otherSong.Cover,
			MusicURL: otherSong.MusicURL,
			Lyric:    otherSong.LyricURL,
		}

		// Check if the song matches the msg and num
		if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
			if num == 0 || song.Num == num {
				filteredSongs = append(filteredSongs, song)
			}
		}
	}

	// Handling API requests from the same system
	for _, api := range metadata.API {
		if api.APIType == "" {
			// Same system API request
			resp, err := http.Get(api.APIURL)
			if err != nil {
				fmt.Println("Error fetching data from internal API: ", err)
				continue // Continue processing the next API, if there are any errors
			}
			defer resp.Body.Close()

			// Decode JSON response
			var internalSongs []Song
			err = json.NewDecoder(resp.Body).Decode(&internalSongs)
			if err != nil {
				fmt.Println("Error decoding internal API response: ", err)
				continue // Continue processing the next API, if there are any errors
			}

			// Modify the num field for internal songs
			for i := range internalSongs {
				songCounter++
				internalSongs[i].Num = songCounter
			}

			// Append internal songs to filteredSongs
			filteredSongs = append(filteredSongs, internalSongs...)
		}
	}

	// If num is specified but no song matches, return an empty array
	if num != 0 && len(filteredSongs) == 0 {
		return []Song{}
	}

	// Return the filtered song array
	return filteredSongs
}

// Helper function to scan the music-uploads directory for artist-song folders
func scanArtistSongFolders(dirPath string) []struct {
	artistName string
	songName   string
	albumName  string
} {
	var artistSongFolders []struct {
		artistName string
		songName   string
		albumName  string
	}

	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("Error opening directory: ", err)
		return artistSongFolders
	}
	defer dir.Close()

	// Read the contents of the directory
	entries, err := dir.Readdir(0)
	if err != nil {
		fmt.Println("Error reading directory entries: ", err)
		return artistSongFolders
	}

	// Iterate over the entries
	for _, entry := range entries {
		if entry.IsDir() {
			// Split the folder name into artist and song name
			parts := strings.Split(entry.Name(), "-")
			if len(parts) >= 2 {
				artistName := strings.TrimSpace(parts[0])
				songName := strings.TrimSpace(strings.Join(parts[1:], "-"))
				albumName := strings.TrimSpace(strings.Join(parts[2:], "@"))
				artistSongFolders = append(artistSongFolders, struct {
					artistName string
					songName   string
					albumName  string
				}{artistName, songName, albumName})
			}
		}
	}

	return artistSongFolders
}

// Helper function to get the cover URL
func getCoverURL(artistName, songName, albumName string) string {
	homeURL := os.Getenv("HOME_URL")
	coverPath := fmt.Sprintf("/file/%s-%s@%s/cover.png", artistName, songName, albumName)
	fullCoverURL := filepath.Join(homeURL, coverPath)
	return checkURL(fullCoverURL)
}

// Helper function to get the MusicURL
func getMusicURL(artistName, songName, albumName string) MusicURL {
	musicURL := MusicURL{
		Audition:     getMusicFileURL(artistName, songName, albumName, "audition", ".mp3"),
		Standard:     getMusicFileURL(artistName, songName, albumName, "standard", ".mp3"),
		Highquality:  getMusicFileURL(artistName, songName, albumName, "highquality", ".mp3"),
		Superquality: getMusicFileURL(artistName, songName, albumName, "superquality", ".mp3"),
		Lossless:     getMusicFileURL(artistName, songName, albumName, "lossless", ".flac"),
		Hires:        getMusicFileURL(artistName, songName, albumName, "hires", ".flac"),
	}
	return musicURL
}

// Helper function to get the Lyric URL
func getLyricURL(artistName, songName, albumName string) Lyric {
	homeURL := os.Getenv("HOME_URL")
	mrcPath := fmt.Sprintf("/file/%s-%s@%s/lyric.mrc", artistName, songName, albumName)
	lrcPath := fmt.Sprintf("/file/%s-%s@%s/lyric.lrc", artistName, songName, albumName)
	mrcURL := filepath.Join(homeURL, mrcPath)
	lrcURL := filepath.Join(homeURL, lrcPath)
	return Lyric{
		Mrc: checkURL(mrcURL),
		Lrc: checkURL(lrcURL),
	}
}

// Helper function to get the music file URL based on quality and format
func getMusicFileURL(artistName, songName, albumName, quality, format string) string {
	homeURL := os.Getenv("HOME_URL")
	musicPath := fmt.Sprintf("/file/%s-%s@%s/%s%s", artistName, songName, albumName, quality, format)
	fullMusicURL := filepath.Join(homeURL, musicPath)
	return checkURL(fullMusicURL)
}

// Read the metadata.json file and parse it into a metadata structure
func readMetadataFile(filePath string) (*Metadata, error) {
	// Retrieve file content
	fileContent, err := GetFileContent(filePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON content
	var metadata Metadata
	err = json.Unmarshal(fileContent, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// Convert the song information in the Metadata structure into a Song array
func getSongArray(metadata *Metadata) []OtherSong {
	return metadata.Other
}

// Helper function to check if a URL is valid
func checkURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	return url
}

// Processing requests.
func IPhandler(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip, nil
	}
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0]), nil
	}
	ip = r.RemoteAddr
	if ip != "" {
		return strings.Split(ip, ":")[0], nil
	}

	return "", fmt.Errorf("unable to obtain IP address information")
}
