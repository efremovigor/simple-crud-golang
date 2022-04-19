package main

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

var validate *validator.Validate

func runWebServer(port string) {
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
	post, deleteError := getPost(context.Param("post-id"))
	if deleteError != nil {
		return echo.NewHTTPError(http.StatusNotFound, deleteError.Error())
	}
	post.remove()
	return context.JSON(http.StatusOK, new(interface{}))
}

func updatePostHandler(context echo.Context) error {
	post, err := getPost(context.Param("post-id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	form := newPost()
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
	_ = post.savePost()

	return context.JSON(http.StatusOK, post)
}

func addPostHandler(context echo.Context) error {
	post := newPost()
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

	if err != post.savePost() {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
	}

	return context.JSON(http.StatusOK, post)
}

func getPostHandler(context echo.Context) error {
	post, err := getPost(context.Param("post-id"))
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

	posts, err := getPosts(page, limit)

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
