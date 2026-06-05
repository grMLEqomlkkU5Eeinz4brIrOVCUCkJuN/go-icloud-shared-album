package utils

import (
	"strconv"

	"github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album/src/types"
)

// EnrichImagesWithUrls returns the album's images with every derivative's URL
// populated from the checksum->URL map produced by core.GetUrls. Derivatives
// whose checksum is absent from urls are dropped, and the surviving derivatives
// are re-keyed by their pixel height (with "-N" suffixes for collisions).
func EnrichImagesWithUrls(
	apiResponse types.ApiResponse,
	urls map[string]string,
) []types.Image {

	var result []types.Image

	for _, photo := range apiResponse.Photos {

		newDerivatives := make(map[string]types.Derivative)
		duplicateCount := make(map[string]int)

		for _, d := range photo.Derivatives {

			url, ok := urls[d.Checksum]
			if !ok {
				continue
			}

			key := strconv.Itoa(d.Height)

			// handle duplicate heights
			if _, exists := newDerivatives[key]; exists {
				duplicateCount[key]++
				key = key + "-" + strconv.Itoa(duplicateCount[key])
			} else {
				duplicateCount[key] = 0
			}

			// attach URL
			d.URL = &url

			newDerivatives[key] = d
		}

		photo.Derivatives = newDerivatives
		result = append(result, photo)
	}

	return result
}
