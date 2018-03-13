package main

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/pbkdf2"
)

var (
	saltLength            = 16
	minimunPasswordLength = 6
	// 72 hour, 3 day
	defaultExpireHour = 72

	jwtSecretKey = "dae6471e3bbe7c5049c9cb2781e7111614670313"
	jwttokenName = "user"
)

// help func for User
func makeSalt() string {
	b := make([]byte, saltLength)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func hashPassword(password string, salt string) string {
	decodeSalt, _ := base64.StdEncoding.DecodeString(salt)
	dk := pbkdf2.Key([]byte(password), decodeSalt, 10000, 64, sha1.New)
	return base64.StdEncoding.EncodeToString(dk)
}

type SimpleCount struct {
	Count int `db:"count" json:"count"`
}

// NullTime struct is for PostgreSQL datetime column which is nullable.
// https://stackoverflow.com/questions/24564619/nullable-time-time-in-golang
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

type Member struct {
	Id      int64  `json:"id" db:"id"`           // id
	Account string `json:"account" db:"account"` // account
}

func (m Member) GenToken() string {
	log.Println("[models GenToken()] id: ", m.Id)
	log.Println("[models GenToken()] account: ", m.Account)
	var tokenString string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      m.Id,
		"account": m.Account,
		"exp":     time.Now().Add(time.Hour * time.Duration(defaultExpireHour)).Unix(),
	})
	tokenString, _ = token.SignedString([]byte(jwtSecretKey))
	return tokenString
}

func (m Member) GenPasswordResetToken() string {
	var tokenString string

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      m.Id,
		"account": m.Account,
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	})
	tokenString, _ = token.SignedString([]byte(jwtSecretKey))

	return tokenString
}

func parseToken(token string) (claims jwt.MapClaims) {
	t, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {

		// Always check the signing method
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			// if tk.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", tk.Header["alg"])
		}
		// Return the key for validation
		return []byte(jwtSecretKey), nil
	})

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims
	} else {
		log.Println(err)
	}
	return claims
}
