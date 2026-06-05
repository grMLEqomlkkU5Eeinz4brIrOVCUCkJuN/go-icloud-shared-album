package types

type ApiResponse struct {
	Photos     map[string]Image `json:"photos"`
	PhotoGuids []string         `json:"photoGuids"`
	Metadata   Metadata         `json:"metadata"`
}
