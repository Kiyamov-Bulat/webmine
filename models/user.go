package models

import (
	"errors"
	"os"
	"strings"

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
	Email     string `json:"email"`
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
			return errors.New("Email address not found")
		}
		return errors.New("Connection error. Please retry")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputedPassword))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Пароль не совпадает!!
		return errors.New("Invalid login credentials. Please try again")
	}
	user.Password = ""
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("AUTH_TOKEN")))
	user.Token = tokenString
	return nil
}

/*
func (user *User) Login() error {
	if user.Email != "art@gmail.com" && user.Email != "bul@gmail.com" {
		return errors.New("Email is not vald")
	}
	if user.Password != "ART_PASS" && user.Password != "BUL_PASS" {
		return errors.New("Password is not vald")
	}
	user.Token = "aGVsbG8gd29ybGQgMTIzIQo="
	user.ID = 1
	return nil
}
*/

func (user *User) Create() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Password = string(hashedPassword)
	GetDB().Create(user)
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserId: user.ID})
	tokenString, _ := token.SignedString([]byte(os.Getenv("AUTH_TOKEN")))
	user.Token = tokenString
	user.Password = ""
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
	var users [2]User

	users[0].Email = os.Getenv("BUL_EMAIL")
	users[1].Email = os.Getenv("ART_EMAIL")
	users[0].Password = os.Getenv("BUL_PASS")
	users[1].Password = os.Getenv("ART_PASS")
	users[0].Name = "Bulat"
	users[1].Name = "Artur"
	users[0].ID = 1
	users[1].ID = 2
	for _, user := range users {
		user.Create()
	}
}
