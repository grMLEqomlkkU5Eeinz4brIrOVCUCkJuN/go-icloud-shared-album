package types

type Metadata struct {
	StreamName    string        `json:"streamName"`
	UserFirstName string        `json:"userFirstName"`
	UserLastName  string        `json:"userLastName"`
	StreamCtag    string        `json:"streamCtag"`
	ItemsReturned int           `json:"itemsReturned"`
	Locations     interface{} `json:"locations"` // Changed from []interface{} to interface{}
}
