package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"

	"github.com/navidrome/navidrome/plugins/pdk/go/metadata"
	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
)

// Configuration keys (must match manifest.json)
const (
	configAPIUrl              = "apiUrl"
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

// audioMuseResponse represents a single track from AudioMuse-AI API
type audioMuseResponse struct {
	ItemID   string  `json:"item_id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Album    string  `json:"album"`
	Distance float64 `json:"distance"`
}

const pluginID = "audiomuseai"

type audioMusePlugin struct{}

func init() {
	metadata.Register(&audioMusePlugin{})
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

func (p *audioMusePlugin) GetSimilarSongsByTrack(input metadata.SimilarSongsByTrackRequest) (*metadata.SimilarSongsResponse, error) {
	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] GetSimilarSongsByTrack called for track ID: %s, Name: %s, Artist: %s", input.ID, input.Name, input.Artist))

	// Read configuration
	apiBaseURL := getConfigString(configAPIUrl, defaultAPIUrl)
	eliminateDuplicates := getConfigBool(configEliminateDuplicates, defaultEliminateDuplicates)
	radiusSimilarity := getConfigBool(configRadiusSimilarity, defaultRadiusSimilarity)

	pdk.Log(pdk.LogDebug, fmt.Sprintf("[AudioMuse] Config - API URL: %s, TrackCount: %d, EliminateDuplicates: %v, RadiusSimilarity: %v",
		apiBaseURL, input.Count, eliminateDuplicates, radiusSimilarity))

	// Build the API URL with query parameters
	params := url.Values{}
	params.Set("item_id", input.ID)
	params.Set("n", strconv.Itoa(int(input.Count)))
	params.Set("eliminate_duplicates", strconv.FormatBool(eliminateDuplicates))
	params.Set("radius_similarity", strconv.FormatBool(radiusSimilarity))

	apiURL := fmt.Sprintf("%s/api/similar_tracks?%s", apiBaseURL, params.Encode())

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Calling API: %s", apiURL))

	// Make HTTP GET request to AudioMuse-AI using PDK
	req := pdk.NewHTTPRequest(pdk.MethodGet, apiURL)
	resp := req.Send()

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] API response status: %d", resp.Status()))

	if resp.Status() != 200 {
		errMsg := fmt.Sprintf("[AudioMuse] ERROR: AudioMuse-AI returned status %d", resp.Status())
		pdk.Log(pdk.LogError, errMsg)
		return nil, fmt.Errorf("AudioMuse-AI returned status %d", resp.Status())
	}

	// Parse JSON response
	var tracks []audioMuseResponse
	body := resp.Body()
	pdk.Log(pdk.LogDebug, fmt.Sprintf("[AudioMuse] Response body length: %d bytes", len(body)))

	if err := json.Unmarshal(body, &tracks); err != nil {
		errMsg := fmt.Sprintf("[AudioMuse] ERROR: Failed to parse response: %v", err)
		pdk.Log(pdk.LogError, errMsg)
		return nil, fmt.Errorf("failed to parse AudioMuse-AI response: %w", err)
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Successfully parsed %d similar tracks", len(tracks)))

	// Sort tracks by distance ascending (smaller distance = more similar)
	//sort.Slice(tracks, func(i, j int) bool { return tracks[i].Distance < tracks[j].Distance })

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

	return &metadata.SimilarSongsResponse{
		Songs: songs,
	}, nil
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
