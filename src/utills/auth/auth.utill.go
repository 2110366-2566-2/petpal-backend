package auth

import (
	"errors"
	"fmt"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	"petpal-backend/src/utills"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginReq struct {
	LoginType string `json:logintype`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRes struct {
	accessToken string
	LoginType   string `json:logintype`
	Username    string `json:"username"`
}

var secretKey = configs.GetJWT_SECRET()

type JWTClaims struct {
	LoginType string `json:logintype`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(u *models.User, loginType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username:  u.Username,
		LoginType: loginType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	ss, err := token.SignedString([]byte(secretKey))
	return ss, err
}

func DecodeToken(tokenString string) (*LoginRes, error) {
	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract and type-assert claims from the token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}
	return &LoginRes{accessToken: tokenString, LoginType: claims.LoginType, Username: claims.Username}, nil
}

func Login(db *models.MongoDB, req *LoginReq) (*LoginRes, error) {
	loginType := req.LoginType
	if loginType == "scvp" {

	} else if loginType == "user" {
		u, err := utills.GetUserByEmail(db, req.Email)
		if err != nil {
			return &LoginRes{}, err
		}
		err = CheckPassword(req.Password, u.Password)
		if err != nil {
			return &LoginRes{}, err
		}
		ss, err := GenerateToken(u, "user")
		if err != nil {
			return &LoginRes{}, err
		}

		return &LoginRes{accessToken: ss, LoginType: "user", Username: u.Username}, nil
	}
	return &LoginRes{}, fmt.Errorf("Invalid Login Type Request")
}
