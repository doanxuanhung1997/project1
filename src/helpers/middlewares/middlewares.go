package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"houze_ops_backend/configs"
	"houze_ops_backend/helpers/common"
	"strings"
	"time"
)

type Claims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  int    `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(id int, email string, role int) (string, error) {
	expirationTime := time.Now().UTC().Add(time.Duration(configs.GetEnvConfig().ServerExpireToken) * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Id:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string

	var jwtKey = []byte(configs.GetEnvConfig().JWTSecretKey)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

var (
	ErrEmptyAuthHeader   = errors.New("auth header is empty")
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
)

// JwtFromHeader get token from header
func JwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)
	if common.IsEmpty(authHeader) {
		return "", ErrEmptyAuthHeader
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && (parts[0] == "Token" || parts[0] == "Bearer")) {
		return "", ErrInvalidAuthHeader
	}
	return parts[1], nil

}

// AuthenticateToken get info listener login by authorization token
func AuthenticateToken(c *gin.Context) (claims Claims, err error) {
	var jwtKey = []byte(configs.GetEnvConfig().JWTSecretKey)
	var token, _ = JwtFromHeader(c, "Authorization")
	tkn, errTK := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// check token invalid
	if errTK != nil || !tkn.Valid {
		return claims, errTK
	}

	return claims, nil
}
