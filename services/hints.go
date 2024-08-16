package services

type Hint struct {
	ID               int    `json:"id"`
	Hint             string `json:"hint"`
	Worth            int    `json:"worth"`
	ParentQuestionID int    `json:"parent_question_id"`
}
