package types

// RawPhoto represents the structure of a photo object directly from the API response
type RawPhoto struct {
	BatchGUID           string                 `json:"batchGuid"`
	Derivatives         map[string]RawDerivative `json:"derivatives"`
	ContributorLastName string                 `json:"contributorLastName"`
	BatchDateCreated    string                 `json:"batchDateCreated"` // Raw string date
	ContributorFirstName string                `json:"contributorFirstName"`
	ContributorFullName string                 `json:"contributorFullName"`
	Caption             string                 `json:"caption"`
	Height              string                 `json:"height"` // Raw string number
	Width               string                 `json:"width"`  // Raw string number
	MediaAssetType      *string                `json:"mediaAssetType,omitempty"`
	PhotoGuid           string                 `json:"photoGuid"`
	DateCreated         string                 `json:"dateCreated"` // Raw string date
}

// RawDerivative represents the structure of a derivative object directly from the API response
type RawDerivative struct {
	Checksum string  `json:"checksum"`
	FileSize string  `json:"fileSize"` // Raw string number
	Width    string  `json:"width"`    // Raw string number
	Height   string  `json:"height"`   // Raw string number
	URL      *string `json:"url,omitempty"`
}

// WebstreamResponse represents the overall structure of the /webstream API response
type WebstreamResponse struct {
	Photos        []RawPhoto    `json:"photos"`
	StreamName    string        `json:"streamName"`
	UserFirstName string        `json:"userFirstName"`
	UserLastName  string        `json:"userLastName"`
	StreamCtag    string        `json:"streamCtag"`
	ItemsReturned string        `json:"itemsReturned"` // Changed from int to string
	Locations     interface{} `json:"locations"`
}
