package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	HttpPort = "8887"
)

type Post struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" validate:"required,gt=2,lt=10" db:"title"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
	Deleted   bool   `json:"deleted" db:"deleted"`
}

type ErrorParams struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ErrorResponse struct {
	Message     string        `json:"message"`
	ErrorParams []ErrorParams `json:"errorParams"`
}

var validate *validator.Validate

func getDbConnection() (db *sql.DB) {
	db, err := sql.Open("postgres", getDbConnectSource())
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	return
}

func newPost() (post *Post) {
	createdAt := time.Now()

	post = new(Post)
	post.CreatedAt = createdAt.Format("2006-01-02 15:04:05.999999")
	post.UpdatedAt = createdAt.Format("2006-01-02 15:04:05.999999")
	post.Deleted = false
	return
}

func decorateErrorParams(err error) (errors []ErrorParams) {
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ErrorParams{Key: err.Field(), Value: err.Value()})
	}
	return
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load("./build/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getDbConnectSource() string {
	return "host=db user=" + goDotEnvVariable("DB_USER") +
		" password=" + goDotEnvVariable("DB_PW") +
		" dbname=" + goDotEnvVariable("DB_NAME") +
		" sslmode=disable"
}

func main() {

	e := echo.New()
	validate = validator.New()

	e.GET("/posts", func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if limit < 1 {
			limit = 10
		}
		if page > 1 {
			page--
		} else if page == 1 {
			page = 0
		}

		db := getDbConnection()
		defer db.Close()

		rows, err := db.Query("SELECT id, name, created_at , updated_at FROM posts WHERE deleted = false order by id OFFSET $1 LIMIT $2", page*limit, limit)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
		}
		var posts []Post
		for rows.Next() {
			post := Post{}
			err = rows.Scan(&post.Id, &post.Name, &post.CreatedAt, &post.UpdatedAt)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
			}
			posts = append(posts, post)
		}

		return c.JSON(http.StatusOK, posts)
	})

	e.GET("/posts/:post-id", func(c echo.Context) error {
		postId := c.Param("post-id")

		db := getDbConnection()
		defer db.Close()

		row := db.QueryRow("SELECT id, name, created_at , updated_at FROM posts WHERE id = $1 and deleted = false", postId)

		post := Post{}
		err := row.Scan(&post.Id, &post.Name, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Not found")
		}

		return c.JSON(http.StatusOK, post)
	})

	e.POST("/posts", func(c echo.Context) (err error) {
		p := newPost()
		decodeErr := json.NewDecoder(c.Request().Body).Decode(&p)
		if decodeErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, decodeErr.Error())
		}

		err = validate.Struct(p)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				ErrorResponse{ErrorParams: decorateErrorParams(err), Message: "Invalid request"},
			)
		}

		db := getDbConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO posts (name, created_at, updated_at, deleted) VALUES ($1, $2, $3, false ) RETURNING id`
		err = db.QueryRow(sqlStatement, p.Name, p.CreatedAt, p.UpdatedAt).Scan(&p.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Unknown error")
		}

		return c.JSON(http.StatusOK, p)
	})

	e.PUT("/posts/:post-id", func(c echo.Context) (err error) {
		postId := c.Param("post-id")
		db := getDbConnection()
		defer db.Close()
		row := db.QueryRow("SELECT id, name, created_at , updated_at FROM posts WHERE id = $1 and deleted = false", postId)
		post := Post{}
		err = row.Scan(&post.Id, &post.Name, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "")
		}

		form := newPost()
		decodeErr := json.NewDecoder(c.Request().Body).Decode(&form)
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
		post.UpdatedAt = time.Now().Format("2006-01-02 15:04:05.999999")

		db.QueryRow(`UPDATE posts SET name=$2, updated_at=$3 WHERE id=$1`, post.Id, post.Name, post.UpdatedAt)

		return c.JSON(http.StatusOK, post)
	})

	e.DELETE("/posts/:post-id", func(c echo.Context) (err error) {
		postId := c.Param("post-id")
		db := getDbConnection()
		defer db.Close()
		_, deleteError := db.Exec("SELECT id, name, created_at , updated_at FROM posts WHERE id = $1 and deleted = false", postId)
		if deleteError != nil {
			return echo.NewHTTPError(http.StatusNotFound, deleteError.Error())
		}
		db.QueryRow(`UPDATE posts SET deleted=true WHERE id=$1`, postId)

		return c.JSON(http.StatusOK, new(interface{}))
	})

	e.Logger.Fatal(e.Start(":" + HttpPort))
}
