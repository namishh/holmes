package handlers

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namishh/holmes/services"
	"github.com/namishh/holmes/views/pages/auth"
	"github.com/namishh/holmes/views/pages/panel"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func csrfMiddleware() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:_csrf",
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieHTTPOnly: true,
		CookieSecure:   true, // Set to true if using HTTPS
	})
}

func (ah *AuthHandler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if user, ok := sess.Values[user_type].(string); !ok || user != "admin" {
			c.Set("FROMPROTECTED", false)
			c.Set("ISADMIN", false)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set("ISADMIN", true)

		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			c.Set("FROMPROTECTED", false)
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if username, ok := sess.Values[user_name_key].(string); ok && len(username) != 0 {
			c.Set(user_name_key, username) // set the username in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		// fromProtected = true
		c.Set("FROMPROTECTED", true)

		return next(c)
	}
}

func (ah *AuthHandler) AdminHandler(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	sess, _ := session.Get(auth_sessions_key, c)
	if user, _ := sess.Values[user_type].(string); user == "admin" {
		return c.Redirect(http.StatusSeeOther, "/su")
	}

	errs := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {

		if c.FormValue("password") != os.Getenv("ADMIN_PASS") {
			c.Set("ISERROR", true)
			errs["pass"] = "Incorrect Password"

			adminLoginView := auth.AdminLogin(csrfToken, errs)
			c.Set("ISERROR", false)
			return renderView(c, auth.AdminLoginIndex(
				"Admin Panel",
				"admin",
				fromProtected,
				c.Get("ISERROR").(bool),
				adminLoginView,
			))
		} else {
			tzone := ""
			if len(c.Request().Header["X-Timezone"]) != 0 {
				tzone = c.Request().Header["X-Timezone"][0]
			}

			sess, _ := session.Get(auth_sessions_key, c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   60 * 60 * 24 * 7, // 1 week
				HttpOnly: true,
			}

			// Set user as authenticated, their username,
			// their ID and the client's time zone

			sess.Values = map[interface{}]interface{}{
				auth_key:      true,
				user_type:     "admin",
				user_id_key:   9999999,
				user_name_key: "admin",
				tzone_key:     tzone,
			}
			sess.Save(c.Request(), c.Response())

			return c.Redirect(http.StatusSeeOther, "/su")
		}

	}

	//sess, _ := session.Get(auth_sessions_key, c)
	// isError = false
	adminLoginView := auth.AdminLogin(csrfToken, errs)
	c.Set("ISERROR", false)
	return renderView(c, auth.AdminLoginIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}

func (ah *AuthHandler) AdminPageHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	users := make([]services.User, 0)
	questions := make([]services.Question, 0)

	users, err := ah.UserServices.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching users")
	}

	questions, err = ah.UserServices.GetAllQuestions()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching questions")
	}

	adminLoginView := panel.PanelHome(fromProtected, users, questions)
	c.Set("ISERROR", false)
	return renderView(c, panel.PanelIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}

func (ah *AuthHandler) AdminQuestionHandler(c echo.Context) error {
	errs := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		title := c.FormValue("title")

		if len(title) == 0 {
			c.Set("ISERROR", true)
			errs["title"] = "Title cannot be empty"
		}

		question := c.FormValue("question")
		if len(question) == 0 {
			c.Set("ISERROR", true)
			errs["question"] = "Question cannot be empty"
		}

		answer := c.FormValue("answer")
		if len(question) == 0 {
			c.Set("ISERROR", true)
			errs["answer"] = "Answer cannot be empty"
		}

		points := c.FormValue("points")
		i, err := strconv.Atoi(points)
		if err != nil || i == 0 {
			c.Set("ISERROR", true)
			errs["points"] = "Points cannot be empty"
		}

		if len(errs) > 0 {
			questionView := panel.PanelQuestion(fromProtected, errs)
			c.Set("ISERROR", false)
			return renderView(c, panel.PanelQuestionIndex(
				"Admin Panel",
				"admin",
				fromProtected,
				c.Get("ISERROR").(bool),
				questionView,
			))
		}

		// create the question
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		images, err := ah.UserServices.MakeArray("images", form, "IMG")
		if err != nil {
			return err
		}
		videos, err := ah.UserServices.MakeArray("videos", form, "VID")
		if err != nil {
			return err
		}
		audios, err := ah.UserServices.MakeArray("audios", form, "AUD")
		if err != nil {
			return err
		}
		log.Println(images, videos, audios)
		err = ah.UserServices.CreateQuestion(services.Question{Question: question, Title: title, Points: i, Answer: answer}, images, videos, audios)
		return c.Redirect(http.StatusSeeOther, "/su")
	}

	adminLoginView := panel.PanelQuestion(fromProtected, errs)
	c.Set("ISERROR", false)
	return renderView(c, panel.PanelQuestionIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}

func (ah *AuthHandler) AdminDeleteTeam(c echo.Context) error {
	teamID := c.Param("id")
	ti, err := strconv.Atoi(teamID)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))

	}

	ah.UserServices.DeleteTeam(ti)

	return c.Redirect(http.StatusSeeOther, "/su")
}

func (ah *AuthHandler) AdminDeleteQuestion(c echo.Context) error {
	qid := c.Param("id")
	ti, err := strconv.Atoi(qid)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))

	}

	ah.UserServices.DeleteQuestion(ti)

	return c.Redirect(http.StatusSeeOther, "/su")
}

func (ah *AuthHandler) AdminDeleteHint(c echo.Context) error {
	qid := c.Param("id")
	ti, err := strconv.Atoi(qid)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))

	}

	ah.UserServices.DeleteHint(ti)

	return c.Redirect(http.StatusSeeOther, "/su/hints")
}
func (ah *AuthHandler) AdminHintsHandler(c echo.Context) error {
	hints, err := ah.UserServices.GetHints()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching hints")
	}
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	adminLoginView := panel.PanelHints(fromProtected, hints)
	c.Set("ISERROR", false)
	return renderView(c, panel.PanelHintsIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}

func (ah *AuthHandler) AdminHintNewHandler(c echo.Context) error {
	errs := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		title := c.FormValue("title")
		level := c.FormValue("level")
		worth := c.FormValue("worth")

		if len(title) == 0 {
			c.Set("ISERROR", true)
			errs["title"] = "You can not scam people with an empty hint."
		}

		l, err := strconv.Atoi(level)
		if err != nil {
			c.Set("ISERROR", true)
			errs["level"] = "Invalid level"
		}

		_, err = ah.UserServices.GetQuestionById(l)
		if err != nil {
			c.Set("ISERROR", true)
			errs["level"] = "Invalid level"
		}

		w, err := strconv.Atoi(worth)
		if err != nil {
			c.Set("ISERROR", true)
			errs["worth"] = "Invalid worth"
		}

		if len(errs) > 0 {
			adminLoginView := panel.PanelNewHint(fromProtected, errs)
			c.Set("ISERROR", false)
			return renderView(c, panel.PanelNewHintIndex(
				"Admin Panel",
				"admin",
				fromProtected,
				c.Get("ISERROR").(bool),
				adminLoginView,
			))
		}

		err = ah.UserServices.CreateHint(services.Hint{Hint: title, ParentQuestionID: l, Worth: w})
		if err != nil {
			c.Set("ISERROR", true)
			errs["title"] = "Error creating hint"
		}

		return c.Redirect(http.StatusSeeOther, "/su/hints")
	}

	adminLoginView := panel.PanelNewHint(fromProtected, errs)
	c.Set("ISERROR", false)
	return renderView(c, panel.PanelNewHintIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}

func (ah *AuthHandler) AdminEditQuestionHandler(c echo.Context) error {
	qid := c.Param("id")
	errs := make(map[string]string)
	inputs := make(map[string]string)
	media := make(map[string][]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	t, err := strconv.Atoi(qid)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	question, err := ah.UserServices.GetQuestionById(t)

	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}

	inputs["title"] = question.Title
	inputs["question"] = question.Question
	inputs["points"] = strconv.Itoa(question.Points)

	media["images"], err = ah.UserServices.GetMedia(fmt.Sprintf("SELECT path FROM images where parent_question_id = %d", t))
	media["videos"], err = ah.UserServices.GetMedia(fmt.Sprintf("SELECT path FROM videos where parent_question_id = %d", t))
	media["audios"], err = ah.UserServices.GetMedia(fmt.Sprintf("SELECT path FROM audios where parent_question_id = %d", t))

	media["limages"], err = ah.UserServices.GetMedia(fmt.Sprintf("SELECT id FROM images where parent_question_id = %d", t))
	media["lvideos"], err = ah.UserServices.GetMedia(fmt.Sprintf("SELECT id FROM videos where parent_question_id = %d", t))
	media["laudios"], err = ah.UserServices.GetMedia(fmt.Sprintf("SELECT id FROM audios where parent_question_id = %d", t))

	if c.Request().Method == "POST" {

		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		images, err := ah.UserServices.MakeArray("images", form, "IMG")
		if err != nil {
			return err
		}
		videos, err := ah.UserServices.MakeArray("videos", form, "VID")
		if err != nil {
			return err
		}
		audios, err := ah.UserServices.MakeArray("audios", form, "AUD")
		if err != nil {
			return err
		}
		log.Println(images, videos, audios)
		err = ah.UserServices.CreateMedia(t, images, videos, audios)

		title := c.FormValue("title")
		qn := c.FormValue("question")
		points := c.FormValue("points")
		answer := c.FormValue("answer")

		if answer == "" {
			answer = question.Answer
		} else {
			by, err := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			answer = string(by)
		}

		if len(title) == 0 {
			c.Set("ISERROR", true)
			errs["title"] = "Empty title."
		}

		p, err := strconv.Atoi(points)
		if err != nil || p == 0 {
			c.Set("ISERROR", true)
			errs["points"] = "Invalid Points."
		}

		if len(errs) > 0 {
			view := panel.PanelEditQuestion(fromProtected, errs, inputs, media)

			c.Set("ISERROR", false)

			return renderView(c, panel.PanelEditQuestionIndex(
				"Edit",
				"",
				fromProtected,
				c.Get("ISERROR").(bool),
				view,
			))
		}

		err = ah.UserServices.UpdateQuestion(t, title, qn, p, answer)
		return c.Redirect(http.StatusSeeOther, "/su")
	}

	view := panel.PanelEditQuestion(fromProtected, errs, inputs, media)

	c.Set("ISERROR", false)

	return renderView(c, panel.PanelEditQuestionIndex(
		"Edit",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}

func (ah *AuthHandler) AdminDeleteImage(c echo.Context) error {
	qid := c.Param("name")
	n, err := strconv.Atoi(qid)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}
	ah.UserServices.DeleteMedia(n, "images")
	return c.Redirect(http.StatusSeeOther, "/su")
}

func (ah *AuthHandler) AdminDeleteAudio(c echo.Context) error {
	qid := c.Param("name")
	n, err := strconv.Atoi(qid)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}
	ah.UserServices.DeleteMedia(n, "audios")
	return c.Redirect(http.StatusSeeOther, "/su")
}

func (ah *AuthHandler) AdminDeleteVideo(c echo.Context) error {
	qid := c.Param("name")
	n, err := strconv.Atoi(qid)
	if err != nil {
		return echo.NewHTTPError(
			echo.ErrNotFound.Code,
			fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
	}
	ah.UserServices.DeleteMedia(n, "videos")
	return c.Redirect(http.StatusSeeOther, "/su")
}
