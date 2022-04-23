package web

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"net/http"
	dbModel "simple-crud-golang/db/model"
	"strconv"
)

var validate *validator.Validate

func RunWebServer(port string) {
	e := echo.New()
	validate = validator.New()

	e.GET("/posts", getPostsHandler)

	e.GET("/posts/:post-id", getPostHandler)

	e.POST("/posts", addPostHandler)

	e.PUT("/posts/:post-id", updatePostHandler)

	e.DELETE("/posts/:post-id", deletePostHandler)

	e.Logger.Fatal(e.Start(":" + port))
}

func deletePostHandler(context echo.Context) error {
	post, deleteError := dbModel.GetPost(context.Param("post-id"))
	if deleteError != nil {
		return echo.NewHTTPError(http.StatusNotFound, deleteError.Error())
	}
	post.Remove()
	return context.JSON(http.StatusOK, new(interface{}))
}

func updatePostHandler(context echo.Context) error {
	post, err := dbModel.GetPost(context.Param("post-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	form := dbModel.NewPost()
	decodeErr := json.NewDecoder(context.Request().Body).Decode(&form)
	if decodeErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, decodeErr.Error())
	}

	err = validate.Struct(form)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrorResponse{ErrorParams: decorateErrorParams(err), Message: "Invalid request"},
		)
	}

	post.Name = form.Name
	_ = post.SavePost()

	return context.JSON(http.StatusOK, post)
}

func addPostHandler(context echo.Context) error {
	post := dbModel.NewPost()
	decodeErr := json.NewDecoder(context.Request().Body).Decode(&post)
	if decodeErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, decodeErr.Error())
	}

	err := validate.Struct(post)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrorResponse{ErrorParams: decorateErrorParams(err), Message: "Invalid request"},
		)
	}

	if err != post.SavePost() {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
	}

	return context.JSON(http.StatusOK, post)
}

func getPostHandler(context echo.Context) error {
	post, err := dbModel.GetPost(context.Param("post-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	return context.JSON(http.StatusOK, post)
}

func getPostsHandler(context echo.Context) error {
	limit, _ := strconv.Atoi(context.QueryParam("limit"))
	page, _ := strconv.Atoi(context.QueryParam("page"))
	if limit < 1 {
		limit = 10
	}
	if page > 1 {
		page--
	} else if page == 1 {
		page = 0
	}

	posts, err := dbModel.GetPosts(page, limit)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
	}

	return context.JSON(http.StatusOK, posts)
}

type ErrorParams struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ErrorResponse struct {
	Message     string        `json:"message"`
	ErrorParams []ErrorParams `json:"errorParams"`
}

func decorateErrorParams(err error) (errors []ErrorParams) {
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ErrorParams{Key: err.Field(), Value: err.Value()})
	}
	return
}
