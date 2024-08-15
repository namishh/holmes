package services

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Title    string `json:"title"`
	Points   int    `json:"points"`
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

func (us *UserService) CreateQuestion(q Question, images []string, audios []string, videos []string) error {
	// Create a question and get its ID
	stmt := `INSERT INTO questions (question, answer, title, points) VALUES (?, ?, ?, ?) RETURNING id`
	ans, err := bcrypt.GenerateFromPassword([]byte(q.Answer), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing answer: %v", err)
		return err
	}
	err = us.UserStore.DB.QueryRow(stmt, q.Question, string(ans), q.Title, q.Points).Scan(&q.ID)
	if err != nil {
		log.Printf("Error inserting question: %v", err)
		return err
	}
	log.Printf("Created question with ID: %d", q.ID)

	// Create images
	for _, img := range images {
		stmt = `INSERT INTO images (path, parent_question_id) VALUES (?, ?)`
		_, err = us.UserStore.DB.Exec(stmt, img, q.ID)
		if err != nil {
			log.Printf("Error inserting image: %v", err)
			return err
		}
	}

	// Create audios
	for _, audio := range audios {
		stmt = `INSERT INTO audios (path, parent_question_id) VALUES (?, ?)`
		_, err = us.UserStore.DB.Exec(stmt, audio, q.ID)
		if err != nil {
			log.Printf("Error inserting audio: %v", err)
			return err
		}
	}

	// Create videos
	for _, video := range videos {
		stmt = `INSERT INTO videos (path, parent_question_id) VALUES (?, ?)`
		_, err = us.UserStore.DB.Exec(stmt, video, q.ID)
		if err != nil {
			log.Printf("Error inserting video: %v", err)
			return err
		}
	}

	return nil
}

// Function to retrieve all questions
func (us *UserService) GetAllQuestions() ([]Question, error) {

	query := `SELECT id, title, points FROM questions`
	questions := make([]Question, 0)

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return questions, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return questions, err
	}

	for rows.Next() {
		var u Question
		err := rows.Scan(&u.ID, &u.Title, &u.Points)
		if err != nil {
			return questions, err
		}

		questions = append(questions, u)
	}

	return questions, nil
}
