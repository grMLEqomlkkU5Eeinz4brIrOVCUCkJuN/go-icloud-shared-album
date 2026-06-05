package types

type Derivative struct {
	Checksum string  `json:"checksum"`
	FileSize int     `json:"fileSize"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	URL      *string `json:"url,omitempty"`
}
