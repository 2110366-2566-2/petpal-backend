package auth

import (
	"errors"
	"fmt"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	user_utills "petpal-backend/src/utills"
	scvp_utills "petpal-backend/src/utills/serviceprovider"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = configs.GetJWT_SECRET()

type JWTClaims struct {
	LoginType string `json:logintype`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(Username string, email string, loginType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username:  Username,
		LoginType: loginType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	ss, err := token.SignedString([]byte(secretKey))
	return ss, err
}

func DecodeToken(tokenString string) (*models.LoginRes, error) {
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
	return &models.LoginRes{AccessToken: tokenString, LoginType: claims.LoginType, Username: claims.Username}, nil
}

func Login(db *models.MongoDB, req *models.LoginReq) (*models.LoginRes, error) {
	loginType := req.LoginType
	if loginType == "scvp" {
		u, err := scvp_utills.GetSVCPByEmail(db, req.Email)
		if err != nil {
			return &models.LoginRes{}, err
		}
		err = CheckPassword(req.Password, u.SVCPPassword)
		if err != nil {
			return &models.LoginRes{}, err
		}
		ss, err := GenerateToken(u.SVCPUsername, u.SVCPEmail, "scvp")
		if err != nil {
			return &models.LoginRes{}, err
		}
		return &models.LoginRes{AccessToken: ss, LoginType: "scvp", Username: u.SVCPUsername}, nil
	} else if loginType == "user" {
		u, err := user_utills.GetUserByEmail(db, req.Email)
		if err != nil {
			return &models.LoginRes{}, err
		}
		err = CheckPassword(req.Password, u.Password)
		if err != nil {
			return &models.LoginRes{}, err
		}
		ss, err := GenerateToken(u.Username, u.Email, "user")
		if err != nil {
			return &models.LoginRes{}, err
		}

		return &models.LoginRes{AccessToken: ss, LoginType: "user", Username: u.Username}, nil
	}
	return &models.LoginRes{}, fmt.Errorf("Invalid Login Type Request")
}
