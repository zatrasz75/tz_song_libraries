package models

type Songs struct {
	ID     int        `json:"ID"`
	Group  string     `json:"group"`
	Song   string     `json:"song"`
	Detail SongDetail `json:"detail"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type LyricResponse struct {
	Chunks     []string `json:"chunks"`
	NextPageID int      `json:"nextPageID"`
}
