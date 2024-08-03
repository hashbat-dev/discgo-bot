package imgur

var (
	BaseUrl = "https://api.imgur.com/3/image"
)

type UploadResponse struct {
	Status  int  `json:"status"`
	Success bool `json:"success"`
	Data    struct {
		ID          string      `json:"id"`
		DeleteHash  string      `json:"deletehash"`
		AccountID   interface{} `json:"account_id"`
		AccountURL  interface{} `json:"account_url"`
		AdType      interface{} `json:"ad_type"`
		AdURL       interface{} `json:"ad_url"`
		Title       string      `json:"title"`
		Description interface{} `json:"description"`
		Name        string      `json:"name"`
		Type        string      `json:"type"`
		Width       int         `json:"width"`
		Height      int         `json:"height"`
		Size        int         `json:"size"`
		Views       int         `json:"views"`
		Section     interface{} `json:"section"`
		Vote        interface{} `json:"vote"`
		Bandwidth   int         `json:"bandwidth"`
		Animated    bool        `json:"animated"`
		Favorite    bool        `json:"favorite"`
		InGallery   bool        `json:"in_gallery"`
		InMostViral bool        `json:"in_most_viral"`
		HasSound    bool        `json:"has_sound"`
		IsAd        bool        `json:"is_ad"`
		NSFW        interface{} `json:"nsfw"`
		Link        string      `json:"link"`
		Tags        []string    `json:"tags"`
		Datetime    int         `json:"datetime"`
		Mp4         string      `json:"mp4"`
		Hls         string      `json:"hls"`
	} `json:"data"`
}
