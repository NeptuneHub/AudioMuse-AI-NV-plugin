package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/navidrome/navidrome/plugins/pdk/go/host"
	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
)

type SimilarArtistsResponse struct {
	Artist           string `json:"artist"`
	ArtistID         string `json:"artist_id"`
	ComponentMatches []struct {
		Artist1RepresentativeSongs []struct {
			ItemID string `json:"item_id"`
			Title  string `json:"title"`
		} `json:"artist1_representative_songs"`
		Artist2RepresentativeSongs []struct {
			ItemID string `json:"item_id"`
			Title  string `json:"title"`
		} `json:"artist2_representative_songs"`
	} `json:"component_matches"`
}

func getSimilarArtists(id string, includeComponentMatches bool) ([]SimilarArtistsResponse, error) {
	apiBaseURL := getConfigString(configAPIUrl, defaultAPIUrl)
	artistCount := getConfigInt("artistSimilarCount", defaultArtistSimilarCount)

	params := url.Values{}
	params.Set("artist", id)
	params.Set("n", strconv.Itoa(artistCount))
	params.Set("include_component_matches", strconv.FormatBool(includeComponentMatches))

	apiURL := fmt.Sprintf("%s/api/similar_artists?%s", apiBaseURL, params.Encode())
	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Calling GetSimilarArtists API for artist ID %s: %s", id, apiURL))

	// Make HTTP GET request to AudioMuse-AI using the host HTTP service.
	// This uses host.HTTPSend as recommended by Navidrome upstream (migrated from pdk.NewHTTPRequest).
	resp, err := host.HTTPSend(host.HTTPRequest{
		Method: "GET",
		URL:    apiURL,
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

	var artists []SimilarArtistsResponse

	body := resp.Body
	if err := json.Unmarshal(body, &artists); err != nil {
		errMsg := fmt.Sprintf("[AudioMuse] ERROR: Failed to parse artist response: %v", err)
		pdk.Log(pdk.LogError, errMsg)
		return nil, fmt.Errorf("failed to parse AudioMuse-AI artist response: %w", err)
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("[AudioMuse] Successfully parsed %d similar artists", len(artists)))

	return artists, nil
}
