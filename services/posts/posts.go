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
	UserId   uint `json:"userId"`
	PostId   uint `json:"postId"`
	VoteType uint `json:"voteType"`
}

type VoteResponse struct {
	PostId        uint   `json:"postId"`
	Body          string `json:"body"`
	VoteType      uint   `json:"voteType"`
	UpvoteCount   uint   `json:"upvoteCount"`
	DownvoteCount uint   `json:"downvoteCount"`
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

func (factory PostFactory) Vote(c *gin.Context) {
	var vote Vote

	if err := c.BindJSON((&vote)); err != nil {
		return
	}

	if vote.VoteType != 0 || vote.VoteType != 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Vote type invalide"})
		return
	}

	// Check if this user has voted on this post before.
	alreadyVotedQuery := `SELECT id, vote_type FROM user_post_votes WHERE user_id = $1 && post_id = $2`
	alreadyVotedRow := factory.Storage.QueryRowContext(context.Background(), alreadyVotedQuery, vote.UserId, vote.PostId)

	var userPostVotesId uint = 0
	var voteType uint
	alreadyVotedRow.Scan(userPostVotesId, voteType)

	// If they have voted on this post before...
	if userPostVotesId != 0 {
		if voteType == vote.VoteType {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "already voted for post"})
			return
		}

		incrementVote := `downvote_count`
		decrementVote := `upvote_count`

		// Increment the new vote type and decrement the former vote type.
		if vote.VoteType == 1 {
			incrementVote = `upvote_count`
			decrementVote = `downvote_count`
		}

		updateCountQuery := `UPDATE posts SET $1 = $1 + 1, $2 = $2 - 1 WHERE post_id = $3`
		_, err := factory.Storage.QueryContext(context.Background(), updateCountQuery, incrementVote, decrementVote, vote.PostId)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update votes"})
			return
		}

		// Update the post to user mapping with the new vote type.
		updateUserPostVotesQuery := `UPDATE user_post_votes SET vote_type = $1 WHERE user_id = $2 && post_id = $3`
		_, err = factory.Storage.QueryContext(context.Background(), updateUserPostVotesQuery, voteType, vote.UserId, vote.PostId)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update votes"})
			return
		}

		// Return the vote response
		voteResponse := factory.getPostForUser(vote.UserId, vote.PostId)
		c.IndentedJSON(http.StatusOK, voteResponse)
		return
	}

	// If they have not voted before...
	incrementVote := `downvote_count`
	if vote.VoteType == 1 {
		incrementVote = `upvote_count`
	}

	// Increment the vote count of the posts.
	updateCountQuery := `UPDATE posts SET $1 = $1 + 1 WHERE post_id = $3`
	_, err := factory.Storage.QueryContext(context.Background(), updateCountQuery, incrementVote, vote.PostId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update votes"})
		return
	}

	// Add a mapping between the user and the post
	insertPostVoteMapping := `INSERT INTO user_post_votes (user_id, post_id, vote_type) VALUES ($1, $2, $3)`
	_, err = factory.Storage.QueryContext(context.Background(), insertPostVoteMapping, vote.UserId, vote.PostId, voteType)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to add user to post"})
		return
	}

	voteResponse := factory.getPostForUser(vote.UserId, vote.PostId)
	c.IndentedJSON(http.StatusOK, voteResponse)
	return
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

func (factory PostFactory) getPostForUser(userId uint, postId uint) VoteResponse {
	query := `SELECT posts.id, body, upvote_count, downvote_count, vote_type FROM user_post_votes
		INNER JOIN posts ON post_id = posts.id
		WHERE user_id = $1 && post_id = $2`
	row := factory.Storage.QueryRowContext(context.Background(), query, userId, postId)

	var voteResponse VoteResponse
	row.Scan(&voteResponse.PostId, &voteResponse.Body, &voteResponse.UpvoteCount, &voteResponse.DownvoteCount, &voteResponse.VoteType)
	return voteResponse
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
