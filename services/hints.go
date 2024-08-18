package services

import "log"

type Hint struct {
	ID               int    `json:"id"`
	Hint             string `json:"hint"`
	Worth            int    `json:"worth"`
	ParentQuestionID int    `json:"parent_question_id"`
}

func (us *UserService) CreateHint(h Hint) error {
	// Create a hint and get its ID
	stmt := `INSERT INTO hints (hint, worth, parent_question_id) VALUES (?, ?, ?) RETURNING id`
	err := us.UserStore.DB.QueryRow(stmt, h.Hint, h.Worth, h.ParentQuestionID).Scan(&h.ID)
	if err != nil {
		log.Printf("Error inserting hint: %v", err)
		return err
	}
	log.Printf("Created hint with ID: %d", h.ID)

	return nil
}

// Get all hints of all questions and sort them by question ID
func (us *UserService) GetHints() ([]Hint, error) {
	// SQL query to select all hints, ordered by parent_question_id
	query := `SELECT id, hint, worth, parent_question_id FROM hints ORDER BY parent_question_id, id`

	// Execute the query
	rows, err := us.UserStore.DB.Query(query)
	if err != nil {
		log.Printf("Error querying hints: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Slice to store the hints
	var hints []Hint

	// Iterate through the rows
	for rows.Next() {
		var h Hint
		err := rows.Scan(&h.ID, &h.Hint, &h.Worth, &h.ParentQuestionID)
		if err != nil {
			log.Printf("Error scanning hint row: %v", err)
			return nil, err
		}
		hints = append(hints, h)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating hint rows: %v", err)
		return nil, err
	}

	return hints, nil
}

func (us *UserService) GetHintsByQuestionID(questionID int) ([]Hint, error) {
	// SQL query to select hints for a specific question ID, ordered by hint ID
	query := `SELECT id, hint, worth, parent_question_id FROM hints WHERE parent_question_id = ? ORDER BY id`

	// Execute the query with the questionID parameter
	rows, err := us.UserStore.DB.Query(query, questionID)
	if err != nil {
		log.Printf("Error querying hints for question ID %d: %v", questionID, err)
		return nil, err
	}
	defer rows.Close()

	// Slice to store the hints
	var hints []Hint

	// Iterate through the rows
	for rows.Next() {
		var h Hint
		err := rows.Scan(&h.ID, &h.Hint, &h.Worth, &h.ParentQuestionID)
		if err != nil {
			log.Printf("Error scanning hint row for question ID %d: %v", questionID, err)
			return nil, err
		}
		hints = append(hints, h)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating hint rows for question ID %d: %v", questionID, err)
		return nil, err
	}

	return hints, nil
}

func (us *UserService) DeleteHint(hintID int) error {
	// SQL query to delete the hint
	query := "DELETE FROM hints WHERE id = ?"

	// Execute the delete statement
	_, err := us.UserStore.DB.Exec(query, hintID)
	if err != nil {
		log.Printf("Error deleting hint with ID %d: %v", hintID, err)
		return err
	}

	return nil
}

func (us *UserService) UnlockHintForTeam(teamID int, hintID int, worth int) error {
	query := `
    INSERT OR IGNORE INTO team_hint_unlocked (team_id, hint_id)
    VALUES (?, ?)
    `
	_, err := us.UserStore.DB.Exec(query, teamID, hintID)
	if err != nil {
		log.Printf("Error unlocking hint %d for team %d: %v", hintID, teamID, err)
		return err
	}

	// Deduct the hint's worth from the team's points
	query = `UPDATE teams SET points = points - ? WHERE id = ?`

	_, err = us.UserStore.DB.Exec(query, worth, teamID)
	if err != nil {
		log.Printf("Error deducting team %d: %v", teamID, err)
		return err
	}

	return nil
}

func (us *UserService) HasTeamUnlockedHint(teamID int, hintID int) (bool, error) {
	query := `
    SELECT EXISTS(SELECT 1 FROM team_hint_unlocked
                  WHERE team_id = ? AND hint_id = ?)
    `
	var exists bool
	err := us.UserStore.DB.QueryRow(query, teamID, hintID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if team %d has unlocked hint %d: %v", teamID, hintID, err)
		return false, err
	}
	return exists, nil
}

func (us *UserService) GetHintById(id int) (string, int, error) {
	var hint string
	var worth int
	query := `SELECT hint, worth FROM hints WHERE id = ?`

	err := us.UserStore.DB.QueryRow(query, id).Scan(&hint, &worth)

	if err != nil {
		log.Printf("Error querying question with ID %d: %v", id, err)
		return hint, worth, err
	}

	log.Printf("Successfully retrieved question with ID: %d", id)
	return hint, worth, nil
}
