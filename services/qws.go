package services

import (
	"log"
	"reflect"
	"sort"
)

type QuestionWithStatus struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Title    string `json:"title"`
	Points   int    `json:"points"`
	Solved   bool   `json:"solved"`
}

func (us *UserService) GetAllQuestionsWithStatus(userID int) ([]QuestionWithStatus, error) {
	query := `SELECT q.id, q.question, q.answer, q.title, q.points,
           CASE WHEN tcq.team_id IS NOT NULL THEN 1 ELSE 0 END as solved
    FROM questions q
    LEFT JOIN team_completed_questions tcq ON q.id = tcq.question_id AND tcq.team_id = ?
    ORDER BY q.id
    `

	rows, err := us.UserStore.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error querying questions with status: %v", err)
		return nil, err
	}
	defer rows.Close()

	var questions []QuestionWithStatus
	for rows.Next() {
		var q QuestionWithStatus
		var solved int
		err := rows.Scan(&q.ID, &q.Question, &q.Answer, &q.Title, &q.Points, &solved)
		if err != nil {
			log.Printf("Error scanning question row: %v", err)
			return nil, err
		}
		q.Solved = solved == 1
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating question rows: %v", err)
		return nil, err
	}

	return questions, nil
}

func (us *UserService) HasCompletedAllQuestions(userID int) (bool, error) {
	// Get all question IDs
	var allQuestionIDs []int
	query := `SELECT id FROM questions`
	rows, err := us.UserStore.DB.Query(query)
	if err != nil {
		log.Printf("Error getting all question IDs: %v", err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return false, err
		}
		allQuestionIDs = append(allQuestionIDs, id)
	}

	// Get completed question IDs for the user
	completedQuestions, err := us.GetCompletedQuestions(userID)
	if err != nil {
		return false, err
	}

	// Compare the slices
	if len(allQuestionIDs) != len(completedQuestions) {
		return false, nil
	}
	sort.Ints(allQuestionIDs)
	sort.Ints(completedQuestions)
	return reflect.DeepEqual(allQuestionIDs, completedQuestions), nil
}

func (us *UserService) MarkQuestionAsCompleted(userID, questionID int) error {
	query := `INSERT OR IGNORE INTO user_completed_questions (user_id, question_id) VALUES (?, ?)`
	_, err := us.UserStore.DB.Exec(query, userID, questionID)
	if err != nil {
		log.Printf("Error marking question %d as completed for user %d: %v", questionID, userID, err)
		return err
	}
	return nil
}

func (us *UserService) GetCompletedQuestions(userID int) ([]int, error) {
	query := `SELECT question_id FROM user_completed_questions WHERE user_id = ?`
	rows, err := us.UserStore.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error getting completed questions for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var completedQuestions []int
	for rows.Next() {
		var questionID int
		if err := rows.Scan(&questionID); err != nil {
			log.Printf("Error scanning completed question: %v", err)
			return nil, err
		}
		completedQuestions = append(completedQuestions, questionID)
	}

	return completedQuestions, nil
}
