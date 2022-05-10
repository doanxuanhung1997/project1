package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/url"
	"houze_ops_backend/config"
	"strings"
	"time"
)

type Claims struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
	jwt.StandardClaims
}

//generate token
func GenerateJWT(id string, phoneNumber string, email string, role int) (string, error) {
	env := config.GetEnvValue()
	expirationTime := time.Now().UTC().Add(time.Duration(env.Server.ExpireToken) * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Id:          id,
		Email:       email,
		PhoneNumber: phoneNumber,
		Role:        role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string

	var jwtKey = []byte(env.Secret.JwtSecretKey)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

//generate token
func GenerateRefreshToken(id string) (string, error) {
	env := config.GetEnvValue()
	expirationTime := time.Now().UTC().Add(time.Duration(env.Server.ExpireToken) * 14 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string

	var jwtKey = []byte(env.Secret.JwtSecretKey)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

var (
	ErrEmptyAuthHeader   = errors.New("auth header is empty")
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
)

//get token from cookie or header
func JwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)
	if authHeader == "" {
		authHeader, _ = JwtFromCookie(c, "Token")
		return authHeader, nil
	} else {
		if authHeader == "" {
			return "", ErrEmptyAuthHeader
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Token") {
			return "", ErrInvalidAuthHeader
		}

		return parts[1], nil
	}

}

//get token from cookie
func JwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Request.Cookie(key)
	if cookie != nil {
		authHeader, _ := url.QueryUnescape(cookie.Value)
		if authHeader == "" {
			return "", ErrEmptyAuthHeader
		}
		return authHeader, nil
	}
	return "", nil
}

//init model token information
type TokenInfo struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
}

//get info listener login by authorization token
func AuthenticateToken(c *gin.Context) (data TokenInfo, err error) {
	var jwtKey = []byte(config.GetSecret())
	var token, _ = JwtFromHeader(c, "Authorization")
	claims := &Claims{}
	tkn, errTK := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// check token invalid
	if errTK != nil || !tkn.Valid {
		return data, errTK
	}

	data.Id = claims.Id
	data.PhoneNumber = claims.PhoneNumber
	data.Email = claims.Email
	data.Role = claims.Role
	return data, nil
}

//get info listener login by authorization token
func GetInfoByToken(token string) (data TokenInfo, err error) {
	var jwtKey = []byte(config.GetSecret())
	claims := &Claims{}
	tkn, errTK := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// check token invalid
	if errTK != nil || !tkn.Valid {
		return data, errTK
	}

	data.Id = claims.Id
	data.PhoneNumber = claims.PhoneNumber
	data.Email = claims.Email
	data.Role = claims.Role
	return data, nil
}

//Logout
func ExtractTokenMetadata(c *gin.Context) error {

	return nil
}
