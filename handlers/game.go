package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
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

	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		answer := c.FormValue("answer")
		log.Print(bcrypt.CompareHashAndPassword([]byte(question.Answer), []byte(answer)))
		if bcrypt.CompareHashAndPassword([]byte(question.Answer), []byte(answer)) == nil {
			err = ah.UserServices.MarkQuestionAsCompleted(c.Get(user_id_key).(int), lvl)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error Validating: %s", err))
			}

			return c.Redirect(http.StatusFound, "/hunt")
		}

		errs["answer"] = "Incorrect Answer"
		quizview := hunt.Question(fromProtected, question, hasCompleted, media, errs)
		c.Set("ISERROR", false)
		return renderView(c, hunt.QuestionIndex(
			"Solve",
			c.Get(user_name_key).(string),
			fromProtected,
			c.Get("ISERROR").(bool),
			quizview,
		))
	}

	quizview := hunt.Question(fromProtected, question, hasCompleted, media, errs)
	c.Set("ISERROR", false)
	return renderView(c, hunt.QuestionIndex(
		"Solve",
		c.Get(user_name_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		quizview,
	))
}
