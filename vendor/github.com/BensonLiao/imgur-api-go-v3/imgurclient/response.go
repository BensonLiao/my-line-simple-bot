package imgurclient

// Response type
type Response struct {
	Data    ResponseData `json:"data"`
	Success bool         `json:"success"`
	Status  int          `json:"status"`
}

// ResponseData type
type ResponseData struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Datetime    int         `json:"datetime"`
	Type        string      `json:"type"`
	Animated    bool        `json:"animated"`
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	Size        int         `json:"size"`
	Views       int         `json:"views"`
	Bandwidth   int         `json:"bandwidth"`
	Vote        interface{} `json:"vote"`
	Favorite    bool        `json:"favorite"`
	NFSW        interface{} `json:"nsfw"`
	Section     interface{} `json:"section"`
	AccountURL  string      `json:"account_url"`
	APIBaseID   int         `json:"account_id"`
	IsAD        bool        `json:"is_ad"`
	InMostViral bool        `json:"in_most_viral"`
	Tags        []string    `json:"tags"`
	ADType      int         `json:"ad_type"`
	ADURL       string      `json:"ad_url"`
	InGellery   bool        `json:"in_gallery"`
	DeleteHash  string      `json:"deletehash"`
	Name        string      `json:"name"`
	Link        string      `json:"link"`
}

// DeleteResponse type
type DeleteResponse struct {
	Data    bool `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}
