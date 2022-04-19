package main

import "time"

type Post struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" validate:"required,gt=2,lt=10" db:"title"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
	Deleted   bool   `json:"deleted" db:"deleted"`
}

func newPost() (post *Post) {
	post = new(Post)
	post.setDefaultDates()
	post.Deleted = false
	return
}

func (post *Post) setDefaultDates() {
	post.setDefaultCreated()
	post.setDefaultUpdated()
}

func (post *Post) setDefaultCreated() {
	post.CreatedAt = time.Now().Format(datetimeLayer)
}

func (post *Post) setDefaultUpdated() {
	post.UpdatedAt = time.Now().Format(datetimeLayer)
}

func (post *Post) savePost() (err error) {
	db := getDbConnection()
	defer db.Close()
	if post.Id != 0 {
		post.setDefaultUpdated()
		db.QueryRow(`UPDATE posts SET name=$2, updated_at=$3 WHERE id=$1`, post.Id, post.Name, post.UpdatedAt)
		return nil
	} else {
		sqlStatement := `INSERT INTO posts (name, created_at, updated_at, deleted) VALUES ($1, $2, $3, false ) RETURNING id`
		return db.QueryRow(sqlStatement, post.Name, post.CreatedAt, post.UpdatedAt).Scan(&post.Id)
	}
}

func getPost(id string) (post Post, err error) {
	db := getDbConnection()
	defer db.Close()
	row := db.QueryRow("SELECT id, name, created_at , updated_at FROM posts WHERE id = $1 and deleted = false", id)
	post = Post{}
	err = row.Scan(&post.Id, &post.Name, &post.CreatedAt, &post.UpdatedAt)
	return
}

func getPosts(page int, limit int) (posts []Post, err error) {

	db := getDbConnection()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, created_at , updated_at FROM posts WHERE deleted = false order by id OFFSET $1 LIMIT $2", page*limit, limit)
	if err != nil {
		return
	}

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Name, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

func (post *Post) remove() {
	db := getDbConnection()
	defer db.Close()
	db.QueryRow(`UPDATE posts SET deleted=true WHERE id=$1`, post.Id)
}
