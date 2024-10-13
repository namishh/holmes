package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/namishh/holmes/services"
	"github.com/namishh/holmes/views/pages"
	"github.com/namishh/holmes/views/pages/hunt"
	"golang.org/x/crypto/bcrypt"
)

func (ah *AuthHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	if !ok {
		return errors.New("invalid type for key 'ISADMIN'")
	}
	sess, _ := session.Get(auth_sessions_key, c)
	// isError = false
	homeView := pages.Home(fromProtected)
	c.Set("ISERROR", false)
	if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
		return renderView(c, pages.HomeIndex(
			"Home",
			"",
			fromProtected,
			c.Get("ISERROR").(bool),
			homeView,
		))
	}

	return renderView(c, pages.HomeIndex(
		"Home",
		sess.Values[user_name_key].(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		homeView,
	))
}

func (ah *AuthHandler) Hunt(c echo.Context) error {
	questions, err := ah.UserServices.GetAllQuestionsWithStatus(c.Get(user_id_key).(int))
	hasCompleted, err := ah.UserServices.HasCompletedAllQuestions(c.Get(user_id_key).(int))
	if err != nil {
		return err
	}
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	quizview := hunt.Hunt(fromProtected, questions, hasCompleted)
	c.Set("ISERROR", false)
	return renderView(c, hunt.HuntIndex(
		"Hunt",
		c.Get(user_name_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		quizview,
	))
}

func (ah *AuthHandler) UnlockHint(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	hastaken, err := ah.UserServices.HasTeamUnlockedHint(c.Get(user_id_key).(int), id)
	if err != nil {
		return err
	}

	hint, worth, err := ah.UserServices.GetHintById(id)

	user, _ := ah.UserServices.CheckUsername(c.Get(user_name_key).(string))

	if !hastaken {
		if user.Points < worth {
			quizview := hunt.OutOfPoints()
			c.Set("ISERROR", true)
			fromProtected, _ := c.Get("FROMPROTECTED").(bool)
			return renderView(c, hunt.OutOfPointsIndex(
				"Hint",
				c.Get(user_name_key).(string),
				fromProtected,
				c.Get("ISERROR").(bool),
				quizview,
			))
		}
		err := ah.UserServices.UnlockHintForTeam(c.Get(user_id_key).(int), id, worth)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	quizview := hunt.Hint(fromProtected, hastaken, hint)
	c.Set("ISERROR", false)
	return renderView(c, hunt.HintIndex(
		"Hint",
		c.Get(user_name_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		quizview,
	))
}

func (ah *AuthHandler) Question(c echo.Context) error {
	errs := make(map[string]string)
	lvl, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	question, err := ah.UserServices.GetQuestionById(lvl)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching question")
	}
	media, err := ah.UserServices.GetMediaByQuestionId(lvl)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching media: %s", err))
	}

	hasCompleted, err := ah.UserServices.IsQuestionSolvedByTeam(c.Get(user_id_key).(int), lvl)
	if err != nil {
		return err
	}

	hints, err := ah.UserServices.GetHintsByQuestionID(lvl)
	if err != nil {
		return err
	}

	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth := sess.Values[user_type]; auth == "admin" {
			return c.String(http.StatusForbidden, "Admin cannot solve questions")
		}

		if hasCompleted {
			return c.String(http.StatusForbidden, "Question already solved")
		}

		answer := c.FormValue("answer")
		if bcrypt.CompareHashAndPassword([]byte(question.Answer), []byte(answer)) == nil {
			err = ah.UserServices.MarkQuestionAsCompleted(c.Get(user_id_key).(int), lvl)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error Validating: %s", err))
			}
			err = ah.UserServices.AddPointsToTeam(c.Get(user_id_key).(int), question.Points)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error adding Points: %s", err))
			}
			err = ah.UserServices.UpdateTeamLastAnsweredQuestion(c.Get(user_id_key).(int))
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error updating time: %s", err))
			}
			return c.Redirect(http.StatusFound, "/hunt")
		}

		errs["answer"] = "Incorrect Answer"
		quizview := hunt.Question(fromProtected, question, hasCompleted, media, errs, hints)
		c.Set("ISERROR", false)
		return renderView(c, hunt.QuestionIndex(
			"Solve",
			c.Get(user_name_key).(string),
			fromProtected,
			c.Get("ISERROR").(bool),
			quizview,
		))
	}

	quizview := hunt.Question(fromProtected, question, hasCompleted, media, errs, hints)
	c.Set("ISERROR", false)
	return renderView(c, hunt.QuestionIndex(
		"Solve",
		c.Get(user_name_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		quizview,
	))
}

func (ah *AuthHandler) Leaderboard(c echo.Context) error {

	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	users, err := ah.UserServices.GetLeaderbaord()

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching Leaderboard: %s", err))
	}

	user := services.User{}

	if c.Get(user_name_key).(string) == "admin" {
		user = services.User{ID: 0, Username: "Admin", Points: 0}
	} else {
		user, err = ah.UserServices.CheckUsername(c.Get(user_name_key).(string))
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching you: %s", err))
		}
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching you: %s", err))
	}

	quizview := hunt.Leaderboard(fromProtected, users, user)
	c.Set("ISERROR", false)
	return renderView(c, hunt.LeaderboardIndex(
		"Leaderboard",
		c.Get(user_name_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		quizview,
	))
}
