package models

import (
	"strings"
	"time"
	"unicode"

	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
)

const (
	MinLenPassword = 6
	MinLenLogin    = 1
	MaxLenLogin    = 25
)

func hasNumeric(s string) bool {
	for _, char := range s {
		if unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func hasSpecialCharacters(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

//nolint:gochecknoinits
func init() {
	govalidator.CustomTypeTagMap.Set("login", func(i interface{}, o interface{}) bool {
		login, ok := i.(string)
		if !ok {
			return false
		}
		return len(login) >= MinLenLogin && len(login) <= MaxLenLogin
	})

	govalidator.CustomTypeTagMap.Set(
		"password",
		func(i interface{}, o interface{}) bool {
			subject, ok := i.(string)
			if !ok {
				return false
			}
			if len(subject) < MinLenPassword {
				return false
			}
			if !govalidator.HasLowerCase(subject) {
				return false
			}
			if !govalidator.HasUpperCase(subject) {
				return false
			}
			if !hasNumeric(subject) {
				return false
			}
			if !hasSpecialCharacters(subject) {
				return false
			}

			return true
		},
	)
}

type User struct {
	ID       uint64 `json:"id"       valid:"required"`
	Login    string `json:"login"    valid:"required,login"`
	Password string `json:"password" valid:"required,password"`
}

type UserWithoutPassword struct {
	ID        uint64    `json:"id"          valid:"required"`
	Login     string    `json:"login"       valid:"required,login"`
	CreatedAt time.Time `json:"created_at"  valid:"required,password"`
}

func (u *UserWithoutPassword) Trim() {
	u.Login = strings.TrimSpace(u.Login)
}

type UserWithoutID struct {
	Login    string `json:"login"    valid:"required,login"`
	Password string `json:"password" valid:"required,password"`
}

func (u *UserWithoutID) Trim() {
	u.Login = strings.TrimSpace(u.Login)
}

func (u *UserWithoutPassword) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	u.Login = sanitizer.Sanitize(u.Login)
}
