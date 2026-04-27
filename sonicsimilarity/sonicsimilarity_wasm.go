//go:build wasip1

package sonicsimilarity

import "github.com/navidrome/navidrome/plugins/pdk/go/pdk"

const NotImplementedCode int32 = -2

//go:wasmexport nd_get_sonic_similar_tracks
func _NdGetSonicSimilarTracks() int32 {
	if getSonicSimilarTracksImpl == nil {
		return NotImplementedCode
	}

	var input GetSonicSimilarTracksRequest
	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return -1
	}

	output, err := getSonicSimilarTracksImpl(input)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	if err := pdk.OutputJSON(output); err != nil {
		pdk.SetError(err)
		return -1
	}

	return 0
}

//go:wasmexport nd_find_sonic_path
func _NdFindSonicPath() int32 {
	if findSonicPathImpl == nil {
		return NotImplementedCode
	}

	var input FindSonicPathRequest
	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return -1
	}

	output, err := findSonicPathImpl(input)
	if err != nil {
		pdk.SetError(err)
		return -1
	}

	if err := pdk.OutputJSON(output); err != nil {
		pdk.SetError(err)
		return -1
	}

	return 0
}
