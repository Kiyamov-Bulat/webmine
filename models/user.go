package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	u "webmine/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	MineModel `gorm:"embedded"`
	Email     string `json:"email" gorm:"unique_index"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Token     string `json:"token" sql:"-"`
}

func (user *User) Validate() (string, bool) {
	if !isValidPassword(user.Password) {
		return "Password is not valid", false
	}
	if !isValidEmail(user.Email) {
		return "Email is not valid", false
	}
	return "Requirement passed", true
}

func (user *User) Login() error {
	inputedPassword := user.Password
	err := GetDB().Table("users").Where("email = ?", user.Email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("email address not found")
		}
		return errors.New("connection error. please retry")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputedPassword))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("invalid login credentials. please try again")
	}
	user.Password = ""
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(u.Getenv("AUTH_TOKEN")))
	user.Token = tokenString
	return nil
}

func (user *User) Create() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Password = string(hashedPassword)
	GetDB().Create(user)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserId: user.ID})
	tokenString, _ := token.SignedString([]byte(u.Getenv("AUTH_TOKEN")))
	user.Token = tokenString
	user.Password = ""
}

func (user *User) Exist() bool {
	tmp := &User{}
	GetDB().Table("users").Where("email = ?", user.Email).First(tmp)
	return tmp.IsValid()
}

func GetUser(id uint) *User {
	model := &User{}
	setItemFromDB("users", model, id)
	return model
}

func isValidPassword(pass string) bool {
	return len(pass) > 6
}
func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

func initUsers() {
	usersJson := os.Getenv("USERS")
	if usersJson == "" {
		user := User{
			Name:      u.Getenv("USER_NAME"),
			Password:  u.Getenv("USER_PASSWORD"),
			Email:     u.Getenv("USER_EMAIL"),
			MineModel: MineModel{ID: 1},
		}
		if !user.Exist() {
			log.Println("here")
			user.Create()
		}
	} else {
		jsonBlob := []byte(usersJson)
		if json.Valid(jsonBlob) {
			users, err := parseUsersEnv(jsonBlob)
			if err != nil {
				log.Println("init users:", err)
				return
			}
			for _, user := range users {
				if !user.Exist() {
					user.Create()
				}
			}
		} else {
			log.Println("error: json: users env is not valid!")
		}
	}
}

func parseUsersEnv(jsonBlob []byte) ([]User, error) {
	var users []User
	err := json.Unmarshal(jsonBlob, &users)
	if err != nil {
		var user User
		err = json.Unmarshal(jsonBlob, &user)
		if err != nil {
			return nil, err
		}
		return []User{user}, err
	}
	return users, err
}
