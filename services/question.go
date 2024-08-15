package services

type Question struct {
	ID         int    `json:"id"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Title      string `json:"title"`
	Points     int    `json:"points"`
	Difficulty string `json:"difficulty"`
}

type Image struct {
	ID               int    `json:"id"`
	Path             string `json:"path"`
	ParentQuestionID int    `json:"parent_question_id"`
}

type Video struct {
	ID               int    `json:"id"`
	Path             string `json:"path"`
	ParentQuestionID int    `json:"parent_question_id"`
}

type Audio struct {
	ID               int    `json:"id"`
	Path             string `json:"path"`
	ParentQuestionID int    `json:"parent_question_id"`
}
