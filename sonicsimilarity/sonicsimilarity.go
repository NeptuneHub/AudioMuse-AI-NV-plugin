package sonicsimilarity

import "github.com/navidrome/navidrome/plugins/pdk/go/metadata"

type GetSonicSimilarTracksRequest struct {
	Song  metadata.SongRef `json:"song"`
	Count int32            `json:"count"`
}

type FindSonicPathRequest struct {
	StartSong metadata.SongRef `json:"startSong"`
	EndSong   metadata.SongRef `json:"endSong"`
	Count     int32            `json:"count"`
}

type SonicSimilarityResponse struct {
	Matches []SonicMatch `json:"matches"`
}

type SonicMatch struct {
	Song       metadata.SongRef `json:"song"`
	Similarity float64         `json:"similarity"`
}

type SonicSimilarity interface {
	GetSonicSimilarTracks(GetSonicSimilarTracksRequest) (SonicSimilarityResponse, error)
	FindSonicPath(FindSonicPathRequest) (SonicSimilarityResponse, error)
}

var getSonicSimilarTracksImpl func(GetSonicSimilarTracksRequest) (SonicSimilarityResponse, error)
var findSonicPathImpl func(FindSonicPathRequest) (SonicSimilarityResponse, error)

func Register(impl SonicSimilarity) {
	getSonicSimilarTracksImpl = impl.GetSonicSimilarTracks
	findSonicPathImpl = impl.FindSonicPath
}
