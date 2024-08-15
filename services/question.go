package services

type Question struct {
	ID         int    `json:"id"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Title      string `json:"title"`
	Points     int    `json:"points"`
	Difficulty string `json:"difficulty"`
}
