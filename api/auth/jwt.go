package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	hmacSampleSecret        = []byte("Immortal budding, the leader is niubility")
	authUserkey             = "usr"
	basciAuthenticateFormat = "Bearer realm=\"Basic Product\", error=\"invalid_token\", error_description=\"%s\""
)

// SignJWTToken create a new hs256 signature token, specifying the user name and the time you want it to expire.
// See https://tools.ietf.org/html/rfc7519#section-4.1 for details
func SignJWTToken(userName string, expireAfter time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"iss":       "Budding Technology",
		"sub":       "Basic Product",
		authUserkey: userName,
		"exp":       time.Now().Add(expireAfter).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSampleSecret)
}

// JWTAuth for authorize
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")

		kv := strings.Split(authorization, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			// See https://tools.ietf.org/html/rfc6750#section-3 for details
			c.Header("WWW-Authenticate", fmt.Sprintf(basciAuthenticateFormat, "Invalid authorization format"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Parse token
		token, err := jwt.Parse(kv[1], VerifyJWTToken)
		if err != nil {
			c.Header("WWW-Authenticate", fmt.Sprintf(basciAuthenticateFormat, err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if !token.Valid {
			c.Header("WWW-Authenticate", fmt.Sprintf(basciAuthenticateFormat, "The access token expired"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("Invalid standard claims, token: %v", token))
			return
		}
		c.Set(gin.AuthUserKey, claims[authUserkey])
	}
}

// VerifyJWTToken verify a jwt token
func VerifyJWTToken(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return hmacSampleSecret, nil
}
