package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"

	"github.com/navidrome/navidrome/plugins/pdk/go/host"
	"github.com/navidrome/navidrome/plugins/pdk/go/metadata"
	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
	"audiomuse-navidrome-plugin/sonicsimilarity"
)

// Configuration keys (must match manifest.json)
const (
	configAPIUrl              = "apiUrl"
	configAPIToken            = "apiToken"
	configEliminateDuplicates = "eliminateDuplicates"
	configRadiusSimilarity    = "radiusSimilarity"
)

// Default values
const (
	defaultAPIUrl              = "http://192.168.3.203:8000"
	defaultArtistSimilarCount  = 10
	defaultEliminateDuplicates = true
	defaultRadiusSimilarity    = true
)

// Compile-time check that we implement necessary interfaces
var _ metadata.SimilarSongsByArtistProvider = (*audioMusePlugin)(nil)
var _ metadata.SimilarSongsByTrackProvider = (*audioMusePlugin)(nil)
var _ metadata.SimilarArtistsProvider = (*audioMusePlugin)(nil)
var _ sonicsimilarity.SonicSimilarity = (*audioMusePlugin)(nil)

// audioMuseTrackResponse represents a single track from AudioMuse-AI API
// and is used for both similar-track and path responses.
type audioMuseTrackResponse struct {
	ItemID   string  `json:"item_id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Album    string  `json:"album"`
	Distance float64 `json:"distance"`
}

type audioMusePathResponse struct {
	Path []audioMuseTrackResponse `json:"path"`
}

const pluginID = "audiomuseai"

type audioMusePlugin struct{}

func init() {
	metadata.Register(&audioMusePlugin{})
	sonicsimilarity.Register(&audioMusePlugin{})
	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Plugin registered successfully (id: %s)", pluginID))
}

// getConfigString retrieves a string config value with a default fallback
func getConfigString(key, defaultValue string) string {
	if value, ok := pdk.GetConfig(key); ok && value != "" {
		return value
	}
	return defaultValue
}

// getConfigInt retrieves an integer config value with a default fallback
func getConfigInt(key string, defaultValue int) int {
	if value, ok := pdk.GetConfig(key); ok && value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getConfigBool retrieves a boolean config value with a default fallback
func getConfigBool(key string, defaultValue bool) bool {
	if value, ok := pdk.GetConfig(key); ok && value != "" {
		return value == "true"
	}
	return defaultValue
}

// authHeaders returns a headers map with a Bearer token if configured, or nil otherwise.
func authHeaders() map[string]string {
	if token := getConfigString(configAPIToken, ""); token != "" {
		return map[string]string{
			"Authorization": "Bearer " + token,
		}
	}
	return nil
}

func (p *audioMusePlugin) GetSimilarSongsByTrack(input metadata.SimilarSongsByTrackRequest) (*metadata.SimilarSongsResponse, error) {
	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] GetSimilarSongsByTrack called for track ID: %s, Name: %s, Artist: %s", input.ID, input.Name, input.Artist))

	tracks, err := p.getAudioMuseSimilarTracks(input.ID, int(input.Count))
	if err != nil {
		return nil, err
	}

	// Convert to Navidrome SongRef format preserving order
	songs := make([]metadata.SongRef, 0, len(tracks))
	for _, track := range tracks {
		songs = append(songs, metadata.SongRef{
			ID:     track.ItemID,
			Name:   track.Title,
			Artist: track.Author,
			Album:  track.Album,
		})
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Returning %d songs to Navidrome", len(songs)))

	return &metadata.SimilarSongsResponse{Songs: songs}, nil
}

func (p *audioMusePlugin) getAudioMuseSimilarTracks(itemID string, count int) ([]audioMuseTrackResponse, error) {
	apiBaseURL := getConfigString(configAPIUrl, defaultAPIUrl)
	eliminateDuplicates := getConfigBool(configEliminateDuplicates, defaultEliminateDuplicates)
	radiusSimilarity := getConfigBool(configRadiusSimilarity, defaultRadiusSimilarity)

	params := url.Values{}
	params.Set("item_id", itemID)
	params.Set("n", strconv.Itoa(count))
	params.Set("eliminate_duplicates", strconv.FormatBool(eliminateDuplicates))
	params.Set("radius_similarity", strconv.FormatBool(radiusSimilarity))

	apiURL := fmt.Sprintf("%s/api/similar_tracks?%s", apiBaseURL, params.Encode())
	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Calling similar_tracks API: %s", apiURL))

	resp, err := host.HTTPSend(host.HTTPRequest{
		Method:  "GET",
		URL:     apiURL,
		Headers: authHeaders(),
	})
	if err != nil {
		errMsg := fmt.Sprintf("[AudioMuse] ERROR: HTTP request failed: %v", err)
		pdk.Log(pdk.LogError, errMsg)
		return nil, fmt.Errorf("AudioMuse-AI HTTP request failed: %w", err)
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] API response status: %d", resp.StatusCode))
	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("[AudioMuse] ERROR: AudioMuse-AI returned status %d", resp.StatusCode)
		pdk.Log(pdk.LogError, errMsg)
		return nil, fmt.Errorf("AudioMuse-AI returned status %d", resp.StatusCode)
	}

	var tracks []audioMuseTrackResponse
	if err := json.Unmarshal(resp.Body, &tracks); err != nil {
		errMsg := fmt.Sprintf("[AudioMuse] ERROR: Failed to parse similar_tracks response: %v", err)
		pdk.Log(pdk.LogError, errMsg)
		return nil, fmt.Errorf("failed to parse AudioMuse-AI similar tracks response: %w", err)
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Successfully parsed %d similar tracks", len(tracks)))
	return tracks, nil
}

func (p *audioMusePlugin) GetSonicSimilarTracks(input sonicsimilarity.GetSonicSimilarTracksRequest) (sonicsimilarity.SonicSimilarityResponse, error) {
	if input.Song.ID == "" {
		return sonicsimilarity.SonicSimilarityResponse{}, fmt.Errorf("song.id is required")
	}

	count := int(input.Count)
	if count <= 0 {
		count = 10
	}

	tracks, err := p.getAudioMuseSimilarTracks(input.Song.ID, count)
	if err != nil {
		return sonicsimilarity.SonicSimilarityResponse{}, err
	}

	matches := make([]sonicsimilarity.SonicMatch, 0, len(tracks))
	for _, track := range tracks {
		matches = append(matches, sonicsimilarity.SonicMatch{
			Song: metadata.SongRef{
				ID:     track.ItemID,
				Name:   track.Title,
				Artist: track.Author,
				Album:  track.Album,
			},
			Similarity: normalizeSimilarity(track.Distance),
		})
	}

	return sonicsimilarity.SonicSimilarityResponse{Matches: matches}, nil
}

func (p *audioMusePlugin) FindSonicPath(input sonicsimilarity.FindSonicPathRequest) (sonicsimilarity.SonicSimilarityResponse, error) {
	if input.StartSong.ID == "" || input.EndSong.ID == "" {
		return sonicsimilarity.SonicSimilarityResponse{}, fmt.Errorf("startSong.id and endSong.id are required")
	}

	count := int(input.Count)
	if count <= 0 {
		count = 25
	}

	apiBaseURL := getConfigString(configAPIUrl, defaultAPIUrl)
	params := url.Values{}
	params.Set("start_song_id", input.StartSong.ID)
	params.Set("end_song_id", input.EndSong.ID)
	params.Set("count", strconv.Itoa(count))
	params.Set("max_steps", strconv.Itoa(count))
	params.Set("path_fix_size", "false")
	params.Set("mood_pct", "100")

	apiURL := fmt.Sprintf("%s/api/find_path?%s", apiBaseURL, params.Encode())
	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Calling FindSonicPath API from %s to %s: %s", input.StartSong.ID, input.EndSong.ID, apiURL))

	resp, err := host.HTTPSend(host.HTTPRequest{
		Method:  "GET",
		URL:     apiURL,
		Headers: authHeaders(),
	})
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("[AudioMuse] ERROR: HTTP request failed: %v", err))
		return sonicsimilarity.SonicSimilarityResponse{}, fmt.Errorf("AudioMuse-AI HTTP request failed: %w", err)
	}

	if resp.StatusCode != 200 {
		pdk.Log(pdk.LogError, fmt.Sprintf("[AudioMuse] ERROR: AudioMuse-AI returned status %d", resp.StatusCode))
		return sonicsimilarity.SonicSimilarityResponse{}, fmt.Errorf("AudioMuse-AI returned status %d", resp.StatusCode)
	}

	var result audioMusePathResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("[AudioMuse] ERROR: Failed to parse FindSonicPath response: %v", err))
		return sonicsimilarity.SonicSimilarityResponse{}, fmt.Errorf("failed to parse AudioMuse-AI find path response: %w", err)
	}

	matches := make([]sonicsimilarity.SonicMatch, 0, len(result.Path))
	for _, item := range result.Path {
		matches = append(matches, sonicsimilarity.SonicMatch{
			Song: metadata.SongRef{
				ID:     item.ItemID,
				Name:   item.Title,
				Artist: item.Author,
				Album:  item.Album,
			},
			Similarity: -1.0,
		})
	}

	return sonicsimilarity.SonicSimilarityResponse{Matches: matches}, nil
}

func normalizeSimilarity(distance float64) float64 {
	similarity := 1.0 - distance
	if similarity < 0 {
		similarity = 0
	}
	if similarity > 1 {
		similarity = 1
	}
	return similarity
}

func (p *audioMusePlugin) GetSimilarSongsByArtist(input metadata.SimilarSongsByArtistRequest) (*metadata.SimilarSongsResponse, error) {
	artists, err := getSimilarArtists(input.ID, true)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)

	// songSlices contains artist songs in alternating order: [baseArtist, relatedArtist1, baseArtist, relatedArtist2, ...]
	songSlices := [][]metadata.SongRef{}

	for _, a := range artists {
		var artist1Songs, artist2Songs []metadata.SongRef

		for _, cm := range a.ComponentMatches {
			for _, s := range cm.Artist1RepresentativeSongs {

				if s.ItemID == "" {
					continue
				}
				if seen[s.ItemID] {
					continue
				}

				seen[s.ItemID] = true
				artist1Songs = append(artist1Songs, metadata.SongRef{ID: s.ItemID, Name: s.Title})
			}

			for _, s := range cm.Artist2RepresentativeSongs {
				if s.ItemID == "" {
					continue
				}

				if seen[s.ItemID] {
					continue
				}

				seen[s.ItemID] = true
				artist2Songs = append(artist2Songs, metadata.SongRef{ID: s.ItemID, Name: s.Title})
			}
		}

		if len(artist1Songs) > 0 {
			songSlices = append(songSlices, artist1Songs)
		}
		if len(artist2Songs) > 0 {
			songSlices = append(songSlices, artist2Songs)
		}
	}

	songs := make([]metadata.SongRef, 0, input.Count)

	// get songs from our slices until we have enough or we ran out
	artistID := 0
	for len(songs) < int(input.Count) && len(songSlices) > 0 {
		song := songSlices[artistID][0] // take a song
		songs = append(songs, song)

		songSlices[artistID] = songSlices[artistID][1:] // remove it from the pool

		if len(songSlices[artistID]) == 0 {
			// this slice has no more songs, remove it
			songSlices = slices.Delete(songSlices, artistID, artistID+1)
			if len(songSlices) == 0 {
				break
			}
		} else {
			// else, go to the next slice
			artistID++
		}

		artistID = artistID % len(songSlices) // loop around if needed
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Returning %d artist-related songs to Navidrome", len(songs)))

	return &metadata.SimilarSongsResponse{Songs: songs}, nil
}

// GetSimilarArtists implements metadata.SimilarArtistsProvider.
func (p *audioMusePlugin) GetSimilarArtists(input metadata.SimilarArtistsRequest) (*metadata.SimilarArtistsResponse, error) {
	artists, err := getSimilarArtists(input.ID, false)
	if err != nil {
		return nil, err
	}

	res := &metadata.SimilarArtistsResponse{
		Artists: make([]metadata.ArtistRef, 0, len(artists)),
	}

	for _, a := range artists {
		res.Artists = append(res.Artists, metadata.ArtistRef{
			ID:   a.ArtistID,
			Name: a.Artist,
		})
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Returning %d related artists to Navidrome", len(res.Artists)))

	return res, nil
}

func main() {}
