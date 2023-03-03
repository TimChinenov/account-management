package posts

import (
	"context"
	"example/account-management/services/storage"
	"example/account-management/services/tokens"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostResponse struct {
	Id            uint   `json:"id"`
	Username      string `json:"username"`
	Body          string `json:"body"`
	UpvoteCount   uint   `json:"upvoteCount"`
	DownvoteCount uint   `json:"downvoteCount"`
}

type Post struct {
	UserId uint   `json:"userId"`
	Body   string `json:"body"`
}

type Vote struct {
	UserId uint `json:"userId"`
	PostId uint `json:"postId"`
}

type PostResponsePaginated struct {
	Posts           []PostResponse `json:"posts"`
	Page            uint           `json:"page"`
	Count           uint           `json:"count"`
	NextPageUrl     string         `json:"nextPageUrl"`
	PreviousPageUrl string         `json:"previousPageUrl"`
}

type PostFactory struct {
	storage.Storage
}

func (factory PostFactory) Create(c *gin.Context) {
	var newPost Post

	if err := c.BindJSON((&newPost)); err != nil {
		return
	}

	if !isCurrentUser(c, newPost.UserId) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "not authorized to create posts"})
		return
	}

	query := `INSERT INTO posts (user_id, body) VALUES ($1, $2);`
	_, err := factory.Storage.QueryContext(context.Background(), query, newPost.UserId, newPost.Body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, "")
}

func (factory PostFactory) Upvote(c *gin.Context) {
	var vote Vote

	if err := c.BindJSON((&vote)); err != nil {
		return
	}

	// check that a mapping has not already been made, if it has been made it needs to be a downvote
	query := `SELECT vote_type FROM user_post_votes WHERE user_id = $1`
	row := factory.Storage.QueryRowContext(context.Background(), query, vote.UserId)

	var vote_type uint
	if err := row.Scan(&vote_type); err == nil {
		// if there is no error, that means a record was found
		if vote_type == 1 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Already upvoted"})
		}
	}

	query = `UPDATE posts SET upvote_count = upvote_count + 1 WHERE post_id = $1`
	_, err := factory.Storage.QueryContext(context.Background(), query, vote.PostId)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, "")
}

func (factory PostFactory) Downvote(c *gin.Context) {
	// get user id, get post id
	// update post count
	// add mapping
}

func (factory PostFactory) Search(c *gin.Context) {
	pageData, _ := c.GetQuery("page")
	pageCountData, _ := c.GetQuery("page_count")

	page, err := strconv.Atoi(pageData)
	if err != nil {
		page = 1
	}

	pageCount, err := strconv.Atoi(pageCountData)
	if err != nil {
		pageCount = 10
	}

	query := `SELECT posts.id, username, body, upvote_count, downvote_count FROM posts
		INNER JOIN users ON user_id = users.id
		ORDER BY posts.id DESC
		OFFSET $1
		LIMIT $2;`
	rows, queryErr := factory.Storage.QueryContext(context.Background(), query, (page-1)*pageCount, pageCount)

	if queryErr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": queryErr})
		return
	}

	posts := make([]PostResponse, 0)
	for rows.Next() {
		var post PostResponse
		queryErr = rows.Scan(&post.Id, &post.Username, &post.Body, &post.UpvoteCount, &post.DownvoteCount)

		if queryErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse posts"})
			return
		}

		posts = append(posts, post)
	}

	queryErr = rows.Err()
	if queryErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to generate"})
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func isCurrentUser(c *gin.Context, userId uint) bool {
	tokenUserId, err := tokens.ExtractTokenId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	if userId != tokenUserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}
