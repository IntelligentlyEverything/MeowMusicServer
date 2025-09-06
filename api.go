package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Api Response.
type Response struct {
	Code          int         `json:"code"`
	Msg           string      `json:"msg"`
	Data          interface{} `json:"data"`
	Tips          string      `json:"tips"`
	Ip            string      `json:"ip"`
	Cache         string      `json:"cache"`
	CacheUpdating bool        `json:"cache_updating"`
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
	Txt string `json:"txt"`
}

// apiHandler is the handler function for API requests.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicServer")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	queryParams := r.URL.Query()
	//key := queryParams.Get("key")
	msg := queryParams.Get("msg")
	singer := queryParams.Get("singer")
	numStr := queryParams.Get("num")

	ip, err := IPhandler(r)
	if err != nil {
		ip = "0.0.0.0"
	}

	if msg == "" {
		response := Response{
			Code:  1,
			Msg:   "API Operation failed: 'msg' parameter is required.",
			Data:  []interface{}{},
			Tips:  "Provide by " + os.Getenv("WEBSITE_NAME"),
			Ip:    ip,
			Cache: "no-cache",
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
				Code:  2,
				Msg:   "Invalid 'num' parameter provided.",
				Data:  []interface{}{},
				Tips:  "Provide by " + os.Getenv("WEBSITE_NAME"),
				Ip:    ip,
				Cache: "no-cache",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Construct a complete file path for cache file
	cacheFilePath := filepath.Join("./cache", msg+".json")
	cacheDir := "./cache"

	// Delete all expired cache files in the background
	cacheTime, err := strconv.Atoi(os.Getenv("API_CACHE_TIME"))
	if err != nil {
		fmt.Println("Error getting API_CACHE_TIME from environment variable: ", err)
		return
	}
	fmt.Println("Deleting expired cache files in the background...")
	go deleteExpiredCacheFiles(cacheDir, cacheTime)

	// Check if cache file exists and is not expired
	if isCacheValid(cacheFilePath) {
		// Read from cache file
		songs, timestamp := readCacheFile(cacheFilePath)

		// Filter songs based on num if provided
		var filteredSongs []Song
		if numStr != "" {
			for _, song := range songs {
				if song.Num == num {
					filteredSongs = append(filteredSongs, song)
				}
			}
		} else {
			filteredSongs = songs
		}

		// Filter songs based on singer if provided
		if singer != "" {
			for _, song := range songs {
				if strings.Contains(song.Singer, singer) {
					filteredSongs = append(filteredSongs, song)
				}
			}
		} else {
			filteredSongs = songs
		}

		// If no songs found, return an empty array
		if len(filteredSongs) == 0 {
			response := Response{
				Code:  3,
				Msg:   "No songs found for the given query.",
				Data:  []interface{}{},
				Tips:  "Provide by " + os.Getenv("WEBSITE_NAME"),
				Ip:    ip,
				Cache: timestamp,
			}
			json.NewEncoder(w).Encode(response)
			fmt.Println("Starting update cache file: ", cacheFilePath)
			go func() {
				newSongs := apiSongHandlerOnMetadata(msg)
				newTimestamp := time.Now().Format(time.RFC3339)
				writeCacheFile(cacheFilePath, newSongs, newTimestamp)
				fmt.Println("Updated cache file: ", cacheFilePath)
			}()
			return
		}

		// Prepare the response with cache timestamp
		response := Response{
			Code:  0,
			Msg:   "API Operation successful.",
			Data:  filteredSongs,
			Tips:  "Provide by " + os.Getenv("WEBSITE_NAME"),
			Ip:    ip,
			Cache: timestamp,
		}

		// Update cache in the background
		response.CacheUpdating = true
		// Encode and send the response
		json.NewEncoder(w).Encode(response)
		fmt.Println("Starting update cache file: ", cacheFilePath)
		go func() {
			newSongs := apiSongHandlerOnMetadata(msg)
			newTimestamp := time.Now().Format(time.RFC3339)
			writeCacheFile(cacheFilePath, newSongs, newTimestamp)
			fmt.Println("Updated cache file: ", cacheFilePath)
		}()
		return
	}

	// If cache file does not exist or is expired, get songs based on msg
	fmt.Println(msg + " cache file not found or expired.")
	var songs []Song
	// get local songs
	localSongs := getLocalSongs(msg)
	songs = append(songs, localSongs...)

	// Filter songs based on num if provided
	var filteredSongs []Song
	if numStr != "" {
		for _, song := range songs {
			if song.Num == num {
				filteredSongs = append(filteredSongs, song)
			}
		}
	} else {
		filteredSongs = songs
	}

	// Filter songs based on singer if provided
	if singer != "" {
		for _, song := range songs {
			if strings.Contains(song.Singer, singer) {
				filteredSongs = append(filteredSongs, song)
			}
		}
	} else {
		filteredSongs = songs
	}

	// If no songs found, return an empty array
	if len(songs) == 0 {
		response := Response{
			Code:          3,
			Msg:           "No songs found for the given query.",
			Data:          []interface{}{},
			Tips:          "Provide by " + os.Getenv("WEBSITE_NAME"),
			Ip:            ip,
			Cache:         "no-cache",
			CacheUpdating: numStr == "",
		}
		json.NewEncoder(w).Encode(response)
		fmt.Println("Starting update cache file: ", cacheFilePath)
		go func() {
			newSongs := apiSongHandlerOnMetadata(msg)
			newTimestamp := time.Now().Format(time.RFC3339)
			writeCacheFile(cacheFilePath, newSongs, newTimestamp)
			fmt.Println("Updated cache file: ", cacheFilePath)
		}()
		return
	}

	// Prepare the response
	response := Response{
		Code:          0,
		Msg:           "API Operation successful.",
		Data:          filteredSongs,
		Tips:          "Provide by " + os.Getenv("WEBSITE_NAME"),
		Ip:            ip,
		Cache:         "no-cache",
		CacheUpdating: numStr == "",
	}

	// Then get songs from external APIs in the background
	response.CacheUpdating = true
	json.NewEncoder(w).Encode(response)
	fmt.Println("Starting update cache file: ", cacheFilePath)
	go func() {
		newSongs := apiSongHandlerOnMetadata(msg)
		newTimestamp := time.Now().Format(time.RFC3339)
		writeCacheFile(cacheFilePath, newSongs, newTimestamp)
		fmt.Println("Updated cache file: ", cacheFilePath)
	}()
}

type API struct {
	APIURL  string `json:"api_url"`
	APIType string `json:"api_type"`
	APIKey  string `json:"api_key"`
	Sources string `json:"sources"`
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

// Local Song handler.
func getLocalSongs(msg string) []Song {
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
		song := Song{
			Num:      0, // Initialize Num to 0
			Song:     artistSongFolder.songName,
			Singer:   artistSongFolder.artistName,
			Album:    artistSongFolder.albumName,
			Cover:    getCoverURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
			MusicURL: getMusicURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
			Lyric:    getLyricURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
		}

		// Check if the song matches the msg
		if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
			songCounter++
			song.Num = songCounter
			filteredSongs = append(filteredSongs, song)
		}
	}

	// Convert OtherSong to Song
	for _, otherSong := range otherSongs {
		song := Song{
			Num:      0, // Initialize Num to 0
			Song:     otherSong.SongName,
			Singer:   otherSong.Singer,
			Album:    otherSong.Album,
			Cover:    otherSong.Cover,
			MusicURL: otherSong.MusicURL,
			Lyric:    otherSong.LyricURL,
		}

		// Check if the song matches the msg
		if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
			songCounter++
			song.Num = songCounter
			filteredSongs = append(filteredSongs, song)
		}
	}

	// Return the filtered song array
	return filteredSongs
}

// API Song handler.
func apiSongHandlerOnMetadata(msg string) []Song {
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
		song := Song{
			Num:      0, // Initialize Num to 0
			Song:     artistSongFolder.songName,
			Singer:   artistSongFolder.artistName,
			Album:    artistSongFolder.albumName,
			Cover:    getCoverURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
			MusicURL: getMusicURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
			Lyric:    getLyricURL(artistSongFolder.artistName, artistSongFolder.songName, artistSongFolder.albumName),
		}

		// Check if the song matches the msg
		if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
			songCounter++
			song.Num = songCounter
			filteredSongs = append(filteredSongs, song)
		}
	}

	// Convert OtherSong to Song
	for _, otherSong := range otherSongs {
		song := Song{
			Num:      0, // Initialize Num to 0
			Song:     otherSong.SongName,
			Singer:   otherSong.Singer,
			Album:    otherSong.Album,
			Cover:    otherSong.Cover,
			MusicURL: otherSong.MusicURL,
			Lyric:    otherSong.LyricURL,
		}

		// Check if the song matches the msg
		if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
			songCounter++
			song.Num = songCounter
			filteredSongs = append(filteredSongs, song)
		}
	}

	for _, api := range metadata.API {
		// Handling API requests from the same system
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

		// Handling API requests from 枫雨API
		if api.APIType == "api.yuanfeng.cn" {
			internalSongs := YuafengAPIResponseHandler(api.APIKey, msg, api.Sources)
			for _, internalSong := range internalSongs {
				song := Song{
					Num:      0, // Initialize Num to 0
					Song:     internalSong.SongName,
					Singer:   internalSong.Singer,
					Album:    internalSong.Album,
					Cover:    internalSong.Cover,
					MusicURL: internalSong.MusicURL,
					Lyric:    internalSong.LyricURL,
				}
				// Check if the song matches the msg
				if strings.Contains(song.Song, msg) || strings.Contains(song.Singer, msg) || strings.Contains(song.Album, msg) {
					songCounter++
					song.Num = songCounter
					filteredSongs = append(filteredSongs, song)
				}
			}
		}
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
				songAndAlbum := strings.Join(parts[1:], "-")
				songAlbumParts := strings.Split(songAndAlbum, "@")
				songName := strings.TrimSpace(songAlbumParts[0])
				var albumName string
				if len(songAlbumParts) > 1 {
					albumName = strings.TrimSpace(songAlbumParts[1])
				}
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
	PORT := os.Getenv("PORT")
	coverPath := fmt.Sprintf("/file/%s-%s@%s/cover.png", url.PathEscape(artistName), url.PathEscape(songName), url.PathEscape(albumName))
	fullCoverURL := fmt.Sprintf("%s%s", homeURL, coverPath)
	testCoverURL := fmt.Sprintf("%s:%s%s", "http://127.0.0.1", PORT, coverPath)
	if checkURL(testCoverURL) == "" {
		return ""
	}
	return fullCoverURL
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
	PORT := os.Getenv("PORT")
	mrcPath := fmt.Sprintf("/file/%s-%s@%s/lyric.mrc", url.PathEscape(artistName), url.PathEscape(songName), url.PathEscape(albumName))
	lrcPath := fmt.Sprintf("/file/%s-%s@%s/lyric.lrc", url.PathEscape(artistName), url.PathEscape(songName), url.PathEscape(albumName))
	txtPath := fmt.Sprintf("/file/%s-%s@%s/lyric.txt", url.PathEscape(artistName), url.PathEscape(songName), url.PathEscape(albumName))
	mrcURL := fmt.Sprintf("%s%s", homeURL, mrcPath)
	lrcURL := fmt.Sprintf("%s%s", homeURL, lrcPath)
	txtURL := fmt.Sprintf("%s%s", homeURL, txtPath)
	testMrcURL := fmt.Sprintf("%s:%s%s", "http://127.0.0.1", PORT, mrcPath)
	testLrcURL := fmt.Sprintf("%s:%s%s", "http://127.0.0.1", PORT, lrcPath)
	testTxtURL := fmt.Sprintf("%s:%s%s", "http://127.0.0.1", PORT, txtPath)
	mrcAvailable := checkURL(testMrcURL) != ""
	lrcAvailable := checkURL(testLrcURL) != ""
	txtAvailable := checkURL(testTxtURL) != ""
	return Lyric{
		Mrc: conditionalURL(mrcAvailable, mrcURL),
		Lrc: conditionalURL(lrcAvailable, lrcURL),
		Txt: conditionalURL(txtAvailable, txtURL),
	}
}

// Helper function to get the music file URL based on quality and format
func getMusicFileURL(artistName, songName, albumName, quality, format string) string {
	homeURL := os.Getenv("HOME_URL")
	PORT := os.Getenv("PORT")
	musicPath := fmt.Sprintf("/file/%s-%s@%s/%s%s", url.PathEscape(artistName), url.PathEscape(songName), url.PathEscape(albumName), url.PathEscape(quality), url.PathEscape(format))
	fullMusicURL := fmt.Sprintf("%s%s", homeURL, musicPath)
	testMusicURL := fmt.Sprintf("%s:%s%s", "http://127.0.0.1", PORT, musicPath)
	if checkURL(testMusicURL) == "" {
		return ""
	}
	return fullMusicURL
}

// Helper function to return URL if available, otherwise return empty string
func conditionalURL(isAvailable bool, url string) string {
	if isAvailable {
		return url
	}
	return ""
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

type YuafengAPIFreeResponse struct {
	Code int `json:"code"`
	Data []struct {
		Num       int    `json:"num"`
		Song      string `json:"song"`
		Singer    string `json:"singer"`
		Cover     string `json:"cover"`
		AlbumName string `json:"album_name"`
	} `json:"data"`
}

type YuafengAPIFreeSingleResponse struct {
	Code int `json:"code"`
	Data struct {
		Song      string `json:"song"`
		Singer    string `json:"singer"`
		Cover     string `json:"cover"`
		AlbumName string `json:"album_name"`
		Music     string `json:"music"`
		Lyric     string `json:"lyric"`
	} `json:"data"`
}

type YuafengAPIQSResponse struct {
	Code int `json:"code"`
	Data []struct {
		Song     string `json:"song"`
		Singer   string `json:"singer"`
		Cover    string `json:"cover"`
		Audition string `json:"music_low"`
		Standard string `json:"music_high"`
		Lyric    string `json:"lyric"`
	} `json:"data"`
}

// 枫雨API response handler.
func YuafengAPIResponseHandler(key string, msg string, sources string) []OtherSong {
	var songs []OtherSong
	if key == "" {
		var url string
		switch sources {
		case "kuwo":
			url = "https://api.yuafeng.cn/API/ly/kwmusic.php"
		case "netease":
			url = "https://api.yuafeng.cn/API/ly/wymusic.php"
		case "migu":
			url = "https://api.yuafeng.cn/API/ly/mgmusic.php"
		case "baidu":
			url = "https://api.yuafeng.cn/API/ly/bdmusic.php"
		default:
			return []OtherSong{}
		}
		resp, err := http.Get(url + "?msg=" + msg)
		if err != nil {
			fmt.Println("Error fetching the data form Yuafeng free API:", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading the response body from Yuafeng free API:", err)
		}
		var response YuafengAPIFreeResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error unmarshalling the data from Yuafeng free API:", err)
		}
		maxNum := 0
		for _, item := range response.Data {
			if item.Num > maxNum {
				maxNum = item.Num
			}
		}
		for i := 1; i <= maxNum; i++ {
			var musicURL = MusicURL{}
			// Loop through different formats
			for _, format := range []string{"LQ", "PQ", "HQ", "SQ"} {
				singleUrl := url + "?msg=" + msg + "&n=" + strconv.Itoa(i) + "&format=" + format
				resp, err := http.Get(singleUrl)
				if err != nil {
					fmt.Println("Error fetching the data form Yuafeng free API:", err)
					continue
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error reading the response body from Yuafeng free API:", err)
					continue
				}
				var singleResponse YuafengAPIFreeSingleResponse
				err = json.Unmarshal(body, &singleResponse)
				if err != nil {
					fmt.Println("Error unmarshalling the data form Yuafeng free API:", err)
					continue
				}
				switch format {
				case "LQ":
					musicURL.Audition = singleResponse.Data.Music
				case "PQ":
					musicURL.Standard = singleResponse.Data.Music
				case "HQ":
					musicURL.Highquality = singleResponse.Data.Music
				case "SQ":
					musicURL.Superquality = singleResponse.Data.Music
				}
				// Check if the song data is valid before appending to songs
				if response.Data[i-1].Song != "" && response.Data[i-1].Singer != "" && response.Data[i-1].AlbumName != "" && response.Data[i-1].Cover != "" && (musicURL.Audition != "" || musicURL.Standard != "" || musicURL.Highquality != "" || musicURL.Superquality != "") {
					song := OtherSong{
						SongName: response.Data[i-1].Song,
						Singer:   response.Data[i-1].Singer,
						Album:    response.Data[i-1].AlbumName,
						Cover:    response.Data[i-1].Cover,
						MusicURL: musicURL,
					}
					songs = append(songs, song)
				}
			}
		}
	} //else {
	//url := "https://api-v2.yuafeng.cn/API/"
	//if sources == "qsmusic" {
	//resp, err := http.Get(url + "?key=" + key + "&msg=" + msg)
	//} else {
	// No other resource support has been provided temporarily.
	//return []OtherSong{}
	//}
	return songs
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

// Helper function to check if cache file is valid
func isCacheValid(filePath string) bool {
	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("Cache file does not exist:", filePath)
		return false
	}

	// Read the cache file metadata
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Error reading cache file metadata: ", err)
		return false
	}

	// Get the cache time from environment variable
	var cacheTime int
	if cacheTimeEnv := os.Getenv("API_CACHE_TIME"); cacheTimeEnv != "" {
		cacheTime, err = strconv.Atoi(cacheTimeEnv)
		if err != nil {
			fmt.Println("Error converting API_CACHE_TIME to int: ", err)
			return false
		}
	} else {
		// Default cache time if not set
		cacheTime = 1 // 1 hour
	}

	// Compare file modification time with current time
	return time.Since(fileInfo.ModTime()).Hours() < float64(cacheTime)
}

// Helper function to delete expired cache files
func deleteExpiredCacheFiles(cacheDir string, cacheTime int) {
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		fmt.Println("Error reading cache directory: ", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(cacheDir, file.Name())
		if file.IsDir() || file.Name() == "embedded.json" {
			continue
		}

		// Read the cache file
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading cache file: ", filePath, err)
			continue
		}

		var cacheData struct {
			Songs     []Song `json:"songs"`
			Timestamp string `json:"timestamp"`
		}
		if err := json.Unmarshal(fileContent, &cacheData); err != nil {
			fmt.Println("Error decoding cache file: ", filePath, err)
			continue
		}

		// Check if the cache is expired
		cacheAt, err := time.Parse(time.RFC3339, cacheData.Timestamp)
		if err != nil {
			fmt.Println("Error parsing cache time: ", filePath, err)
			continue
		}

		cacheDuration := time.Since(cacheAt)
		if cacheDuration > time.Duration(cacheTime)*time.Hour {
			// Delete the expired cache file
			if err := os.Remove(filePath); err != nil {
				fmt.Println("Error deleting expired cache file: ", filePath, err)
			} else {
				fmt.Println("Deleted expired cache file: ", filePath)
			}
		}
	}
}

// Helper function to read from cache file
func readCacheFile(filePath string) ([]Song, string) {
	var cacheFile struct {
		Songs     []Song `json:"songs"`
		Timestamp string `json:"timestamp"`
	}

	// Open the cache file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening cache file: ", err)
		return []Song{}, ""
	}
	defer file.Close()

	// Decode JSON response
	err = json.NewDecoder(file).Decode(&cacheFile)
	if err != nil {
		fmt.Println("Error decoding cache file: ", err)
		return []Song{}, ""
	}

	return cacheFile.Songs, cacheFile.Timestamp
}

// Helper function to write to cache file
func writeCacheFile(filePath string, songs []Song, timestamp string) {
	cacheFile := struct {
		Songs     []Song `json:"songs"`
		Timestamp string `json:"timestamp"`
	}{
		Songs:     songs,
		Timestamp: timestamp,
	}

	// Create the cache directory if it does not exist
	cacheDir := filepath.Dir(filePath)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating cache directory: ", err)
			return
		}
	}

	// Create or update the cache file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating cache file: ", err)
		return
	}
	defer file.Close()

	// Encode JSON response
	err = json.NewEncoder(file).Encode(cacheFile)
	if err != nil {
		fmt.Println("Error encoding cache file: ", err)
	}
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
