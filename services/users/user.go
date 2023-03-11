package users

import (
	"context"
	"database/sql"
	"example/account-management/services/storage"
	"example/account-management/services/tokens"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type UserRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UpdateScoreRequst struct {
	Score int `json:"score"`
}

type UserFactory struct {
	storage.Storage
}

type UserStore interface {
	Create(*gin.Context)
}

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return &userStore{db: db}
}

func (u *userStore) Create(c *gin.Context) {
	var newUser UserRequest

	if err := c.BindJSON((&newUser)); err != nil {
		return
	}

	hashedPassword, err := hashPassword(newUser.Password)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	query := `INSERT INTO users (username, password, score) VALUES ($1, $2, $3);`
	_, err = u.db.QueryContext(context.Background(), query, newUser.Username, hashedPassword, 0)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query = `SELECT id, username, score FROM users WHERE username=$1;`
	row := u.db.QueryRowContext(context.Background(), query, newUser.Username)

	var id int
	var username string
	var score int

	err = row.Scan(&id, &username, &score)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var userResponse UserResponse = UserResponse{ID: id, Username: username, Score: score}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

func (factory UserFactory) Create(c *gin.Context) {
	var newUser UserRequest

	if err := c.BindJSON((&newUser)); err != nil {
		return
	}

	hashedPassword, err := hashPassword(newUser.Password)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	query := `INSERT INTO users (username, password, score) VALUES ($1, $2, $3);`
	_, err = factory.Storage.QueryContext(context.Background(), query, newUser.Username, hashedPassword, 0)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query = `SELECT id, username, score FROM users WHERE username=$1;`
	row := factory.Storage.QueryRowContext(context.Background(), query, newUser.Username)

	var id int
	var username string
	var score int

	err = row.Scan(&id, &username, &score)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var userResponse UserResponse = UserResponse{ID: id, Username: username, Score: score}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

func (factory UserFactory) Get(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to parse user id"})
		return
	}

	var userResponse UserResponse

	userResponse, err = factory.GetUserById(userID)

	if err != nil || userResponse.ID == 0 || userResponse.Username == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.IndentedJSON(http.StatusFound, userResponse)
}

func (factory UserFactory) Search(c *gin.Context) {
	usernameSearch, err := c.GetQuery("username")
	usernameSearch = strings.TrimSpace(usernameSearch)

	if !err || usernameSearch == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no parameters found"})
		return
	}

	query := `SELECT id, username, score FROM users WHERE username LIKE '%' || $1 || '%' LIMIT 10`
	rows, queryErr := factory.Storage.QueryContext(context.Background(), query, usernameSearch)

	if queryErr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to search"})
		return
	}

	users := make([]UserResponse, 0)
	for rows.Next() {
		var user UserResponse
		queryErr = rows.Scan(&user.ID, &user.Username, &user.Score)

		if queryErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse users"})
			return
		}

		users = append(users, user)
	}

	queryErr = rows.Err()
	if queryErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to generate"})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (factory UserFactory) UpdatePoints(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to parse user id"})
		return
	}

	var scoreRequest UpdateScoreRequst

	if err := c.BindJSON((&scoreRequest)); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse score"})
		return
	}

	query := `UPDATE users SET score = $1 WHERE id = $2`
	row := factory.Storage.QueryRowContext(context.Background(), query, scoreRequest.Score, userID)

	if row.Err() != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to update score"})
	}

	userResponse, err := factory.GetUserById(userID)

	if err != nil || userResponse.ID == 0 || userResponse.Username == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.IndentedJSON(http.StatusFound, userResponse)
}

func (factory UserFactory) Login(c *gin.Context) {
	var loginUser UserRequest

	if err := c.BindJSON((&loginUser)); err != nil {
		return
	}

	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := factory.Storage.QueryRowContext(context.Background(), query, loginUser.Username)

	if row.Err() != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid username"})
		return
	}

	var foundId int
	var foundUsername string
	var foundPassword string

	row.Scan(&foundId, &foundUsername, &foundPassword)

	if !checkPasswordHash(loginUser.Password, foundPassword) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	token, err := tokens.CreateToken(foundId)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to generate token"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

func (factory UserFactory) GetUserById(userID int) (UserResponse, error) {
	query := `SELECT id, username, score FROM users WHERE id=$1;`
	row := factory.Storage.QueryRowContext(context.Background(), query, userID)

	var userResponse UserResponse

	err := row.Scan(&userResponse.ID, &userResponse.Username, &userResponse.Score)

	return userResponse, err
}

func (factory UserFactory) CurrentUser(c *gin.Context) {
	userId, err := tokens.ExtractTokenId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := factory.GetUserById(int(userId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
