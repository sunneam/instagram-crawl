package model

type R struct {
	Media []struct{
		User      string  `json:"user"`
		Caption   string  `json:"caption"`
		Code      string  `json:"code"`
		Date      string  `json:"date"`
		MediaType float64 `json:"media_type"`
		Owner     struct {
			FullName      string `json:"full_name"`
			Pk            string `json:"pk"`
			ProfilePicUrl string `json:"profile_pic_url"`
			Username      string `json:"username"`
		} `json:"owner"`
		ThumbnailSrc  string `json:"thumbnail_src"`
	} `json:"media"`
	Next  float64 `json:"next"`
}

//类型1,共三种类型
type Tone struct {
	Nodes []struct{
		ThumbnailSrc string 	`json:"thumbnail_src"`
		TakenAt interface{} `json:"taken_at"`
		Code string 	`json:"code"`
		Caption string `json:"caption"`
		DisplaySrc string `json:"display_src"`
		MediaType float64 `json:"media_type"`
	} `json:"nodes"`

	PageInfo struct{
		HasNextPage bool `json:"has_next_page"`
		EndCursor string `json:"end_cursor"`
	}`json:"page_info"`
}