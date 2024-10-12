package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
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

func (us *UserService) MakeArray(label string, form *multipart.Form, short string) (list []string, err error) {
	bucketName := os.Getenv("BUCKET_NAME")
	files := form.File[label]
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return list, err
		}
		defer src.Close()

		u := uuid.New().String()
		filename := fmt.Sprintf("%s-%s%s", short, u, filepath.Ext(file.Filename))

		_, err = us.MinioClient.PutObject(context.Background(), bucketName, filename, src, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
		if err != nil {
			return list, fmt.Errorf("failed to upload file to MinIO: %v", err)
		}

		fmt.Println(bucketName, filename)

		presignedURL, err := us.MinioClient.PresignedGetObject(context.Background(), bucketName, filename, time.Second*60*60*24*7, nil)
		if err != nil {
			return list, fmt.Errorf("failed to generate presigned URL: %v", err)
		}
		list = append(list, presignedURL.String())
	}
	return list, nil
}

func (us *UserService) CreateMedia(ID int, images []string, videos []string, audios []string) error {

	// Create images
	for _, img := range images {
		stmt := `INSERT INTO images (path, parent_question_id) VALUES (?, ?)`
		_, err := us.UserStore.DB.Exec(stmt, img, ID)
		if err != nil {
			log.Printf("Error inserting image: %v", err)
			return err
		}
	}

	// Create audios
	for _, audio := range audios {
		stmt := `INSERT INTO audios (path, parent_question_id) VALUES (?, ?)`
		_, err := us.UserStore.DB.Exec(stmt, audio, ID)
		if err != nil {
			log.Printf("Error inserting audio: %v", err)
			return err
		}
	}

	// Create videos
	for _, video := range videos {
		stmt := `INSERT INTO videos (path, parent_question_id) VALUES (?, ?)`
		_, err := us.UserStore.DB.Exec(stmt, video, ID)
		if err != nil {
			log.Printf("Error inserting video: %v", err)
			return err
		}
	}

	return nil
}

func (us *UserService) CreateQuestion(q Question, images []string, videos []string, audios []string) error {
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

	us.CreateMedia(q.ID, images, videos, audios)

	return nil
}

// Function to retrieve all questions
func (us *UserService) GetAllQuestions() ([]Question, error) {

	query := `SELECT id, title, points FROM questions ORDER BY points ASC`
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

func (us *UserService) DeleteQuestion(id int) error {
	query := `DELETE FROM questions WHERE id = ?`
	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	stmt.Exec(id)

	c := make([]string, 4)
	c[0] = "images"
	c[1] = "audios"
	c[2] = "videos"
	c[3] = "hints"

	for _, table := range c {
		query = fmt.Sprintf(`DELETE FROM %s  WHERE parent_question_id = ?`, table)
		stmt, err = us.UserStore.DB.Prepare(query)
		if err != nil {
			return err
		}

		defer stmt.Close()

		stmt.Exec(id)

	}

	return nil
}

func (us *UserService) DeleteMedia(id int, table string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, table)
	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	stmt.Exec(id)

	return nil
}

// get id by path
func (us *UserService) GetIdByPath(path string, table string) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE path = ?`, table)
	var id int
	err := us.UserStore.DB.QueryRow(query, path).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (us *UserService) GetQuestionById(id int) (Question, error) {
	var q Question

	query := `SELECT id, question, answer, title, points FROM questions WHERE id = ?`

	err := us.UserStore.DB.QueryRow(query, id).Scan(&q.ID, &q.Question, &q.Answer, &q.Title, &q.Points)

	if err != nil {
		log.Printf("Error querying question with ID %d: %v", id, err)
		return Question{}, err
	}

	log.Printf("Successfully retrieved question with ID: %d", id)
	return q, nil
}

func (us *UserService) GetMedia(query string) ([]string, error) {
	media := make([]string, 0)
	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return media, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return media, err
	}

	for rows.Next() {
		var u string
		err := rows.Scan(&u)
		if err != nil {
			return media, err
		}

		media = append(media, u)
	}

	return media, nil
}

func (us *UserService) UpdateQuestion(id int, title string, question string, points int, answer string) error {
	query := `UPDATE questions
              SET title = ?, question = ?, points = ?, answer = ?
              WHERE id = ?`

	// Execute the update statement
	_, err := us.UserStore.DB.Exec(query, title, question, points, answer, id)
	if err != nil {
		log.Printf("Error updating question with ID %d: %v", id, err)
		return err
	}

	log.Printf("Update operation completed for question with ID: %d", id)
	return nil
}

// make a function that takes questions id and returns all the media associated with it
func (us *UserService) GetMediaByQuestionId(id int) (map[string][]string, error) {
	m := make(map[string][]string)

	stmt := fmt.Sprintf(`SELECT path FROM images WHERE parent_question_id = %d`, id)
	images, err := us.GetMedia(stmt)
	if err != nil {
		return nil, err
	}

	m["images"] = images

	stmt = fmt.Sprintf(`SELECT path FROM videos WHERE parent_question_id = %d`, id)
	videos, err := us.GetMedia(stmt)
	if err != nil {
		return nil, err
	}

	m["videos"] = videos

	stmt = fmt.Sprintf(`SELECT path FROM audios WHERE parent_question_id = %d`, id)
	audios, err := us.GetMedia(stmt)
	if err != nil {
		return nil, err
	}

	m["audios"] = audios

	return m, nil
}

func (us *UserService) AddPointsToTeam(teamID int, points int) error {
	query := `
    UPDATE teams
    SET points = points + ?
    WHERE id = ?
    `

	// Execute the update
	_, err := us.UserStore.DB.Exec(query, points, teamID)
	if err != nil {
		log.Printf("Error adding points to team %d: %v", teamID, err)
		return err
	}

	log.Printf("Successfully added %d points to team %d", points, teamID)
	return nil
}
