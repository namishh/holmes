package services

import (
	"log"
	"time"
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
    ORDER BY q.points ASC
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

func (us *UserService) HasCompletedAllQuestions(teamID int) (bool, error) {
	// Get total number of questions
	var totalQuestions int
	queryTotal := `SELECT COUNT(*) FROM questions`
	err := us.UserStore.DB.QueryRow(queryTotal).Scan(&totalQuestions)
	if err != nil {
		log.Printf("Error getting total question count: %v", err)
		return false, err
	}

	// Get number of completed questions for the team
	var completedCount int
	queryCompleted := `SELECT COUNT(*) FROM team_completed_questions WHERE team_id = ?`
	err = us.UserStore.DB.QueryRow(queryCompleted, teamID).Scan(&completedCount)
	if err != nil {
		log.Printf("Error getting completed question count for team %d: %v", teamID, err)
		return false, err
	}

	// Compare counts
	if totalQuestions == 0 {
		return false, nil
	}
	return completedCount >= totalQuestions, nil
}

func (us *UserService) MarkQuestionAsCompleted(userID, questionID int) error {
	query := `INSERT OR IGNORE INTO team_completed_questions (team_id, question_id) VALUES (?, ?)`
	_, err := us.UserStore.DB.Exec(query, userID, questionID)
	if err != nil {
		log.Printf("Error marking question %d as completed for user %d: %v", questionID, userID, err)
		return err
	}
	return nil
}

func (us *UserService) GetCompletedQuestions(userID int) ([]int, error) {
	query := `SELECT question_id FROM team_completed_questions WHERE team_id = ?`
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

// Take a question ID, team ID and check if the question is solved by the team
// Return true if solved, false otherwise

func (us *UserService) IsQuestionSolvedByTeam(teamID, questionID int) (bool, error) {
	query := `SELECT COUNT(*) FROM team_completed_questions WHERE team_id = ? AND question_id = ?`
	var count int
	err := us.UserStore.DB.QueryRow(query, teamID, questionID).Scan(&count)
	if err != nil {
		log.Printf("Error checking if question %d is solved by team %d: %v", questionID, teamID, err)
		return false, err
	}
	return count > 0, nil
}

func (us *UserService) UpdateTeamLastAnsweredQuestion(teamID int) error {
	query := `
    UPDATE teams
    SET last_answered_question = ?
    WHERE id = ?
    `

	// Get current timestamp
	currentTime := time.Now()

	// Execute the update
	_, err := us.UserStore.DB.Exec(query, currentTime, teamID)
	if err != nil {
		log.Printf("Error updating last answered question for team %d: %v", teamID, err)
		return err
	}

	log.Printf("Update operation completed for team %d at %v", teamID, currentTime)
	return nil
}

type LeaderBoardUser struct {
	Username string
	Points   int
}

func (us *UserService) GetLeaderbaord() ([]LeaderBoardUser, error) {
	stmt := `SELECT name, points FROM teams ORDER BY points DESC, last_answered_question ASC;`
	rows, err := us.UserStore.DB.Query(stmt)
	if err != nil {
		log.Printf("Eror fetching leaderboard")
		return nil, err
	}
	defer rows.Close()

	var users []LeaderBoardUser

	for rows.Next() {
		var user LeaderBoardUser
		if err := rows.Scan(&user.Username, &user.Points); err != nil {
			log.Printf("Error scanning completed question: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
