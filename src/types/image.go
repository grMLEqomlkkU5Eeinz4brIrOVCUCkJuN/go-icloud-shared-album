package types

import "time"

type Image struct {
	BatchGUID            string                `json:"batchGuid"`
	Derivatives          map[string]Derivative `json:"derivatives"`
	ContributorLastName  string                `json:"contributorLastName"`
	BatchDateCreated     time.Time             `json:"batchDateCreated"`
	ContributorFirstName string                `json:"contributorFirstName"`
	ContributorFullName  string                `json:"contributorFullName"`
	Caption              string                `json:"caption"`
	Height               int                   `json:"height"`
	Width                int                   `json:"width"`
	MediaAssetType       *string               `json:"mediaAssetType,omitempty"`
	PhotoGuid            string                `json:"photoGuid"`   // Added based on TypeScript
	DateCreated          time.Time             `json:"dateCreated"` // Added based on TypeScript
}
