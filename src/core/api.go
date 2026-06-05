package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album/src/types"
)

// Keep these static
var headers = map[string]string{
	"Origin":          "https://www.icloud.com",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Content-Type":    "text/plain",
	"Accept":          "*/*",
	"Referer":         "https://www.icloud.com/sharedalbum/",
	"Connection":      "keep-alive",
}

func parseDate(dateStr string) *time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z0700", dateStr)
		if err != nil {
			return nil
		}
	}
	return &t
}

func GetApiResponse(baseUrl string) (*types.ApiResponse, error) {
	url := baseUrl + "webstream"
	dataString := []byte(`{"streamCtag":null}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataString))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var webstreamResponse types.WebstreamResponse
	err = json.Unmarshal(body, &webstreamResponse)
	if err != nil {
		return nil, err
	}

	photos := make(map[string]types.Image)
	photoGuids := []string{}

	for _, photo := range webstreamResponse.Photos {
		height, _ := strconv.Atoi(photo.Height)
		width, _ := strconv.Atoi(photo.Width)

		derivatives := make(map[string]types.Derivative)
		for key, rawDerivative := range photo.Derivatives {
			fileSize, _ := strconv.Atoi(rawDerivative.FileSize)
			dWidth, _ := strconv.Atoi(rawDerivative.Width)
			dHeight, _ := strconv.Atoi(rawDerivative.Height)

			derivatives[key] = types.Derivative{
				Checksum: rawDerivative.Checksum,
				FileSize: fileSize,
				Width:    dWidth,
				Height:   dHeight,
				URL:      rawDerivative.URL,
			}
		}

		photos[photo.PhotoGuid] = types.Image{
			BatchGUID:            photo.BatchGUID,
			Derivatives:          derivatives,
			ContributorLastName:  photo.ContributorLastName,
			BatchDateCreated:     *parseDate(photo.BatchDateCreated),
			ContributorFirstName: photo.ContributorFirstName,
			ContributorFullName:  photo.ContributorFullName,
			Caption:              photo.Caption,
			Height:               height,
			Width:                width,
			MediaAssetType:       photo.MediaAssetType,
			PhotoGuid:            photo.PhotoGuid,
			DateCreated:          *parseDate(photo.DateCreated),
		}
		photoGuids = append(photoGuids, photo.PhotoGuid)
	}

	itemsReturned, _ := strconv.Atoi(webstreamResponse.ItemsReturned) // Convert string to int
	apiResponse := &types.ApiResponse{
		Photos:     photos,
		PhotoGuids: photoGuids,
		Metadata: types.Metadata{
			StreamName:    webstreamResponse.StreamName,
			UserFirstName: webstreamResponse.UserFirstName,
			UserLastName:  webstreamResponse.UserLastName,
			StreamCtag:    webstreamResponse.StreamCtag,
			ItemsReturned: itemsReturned, // Use the converted int
			Locations:     webstreamResponse.Locations,
		},
	}

	return apiResponse, nil
}

func GetUrls(baseUrl string, photoGuids []string) (map[string]string, error) {
	url := baseUrl + "webasseturls"

	requestBody := map[string][]string{"photoGuids": photoGuids}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // Replaced ioutil.ReadAll
	if err != nil {
		return nil, err
	}

	var response struct {
		Items map[string]struct {
			URLPath     string `json:"url_path"`
			URLLocation string `json:"url_location"`
		} `json:"items"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	items := make(map[string]string)
	for itemId, item := range response.Items {
		items[itemId] = "https://" + item.URLLocation + item.URLPath
	}

	return items, nil
}
