package auth

import (
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

func GenerateToken(u *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	ss, err := token.SignedString([]byte(secretKey))
	return ss, err
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
		ss, err := GenerateToken(u)
		if err != nil {
			return &LoginRes{}, err
		}

		return &LoginRes{accessToken: ss, LoginType: "user", Username: u.Username}, nil
	}
	return &LoginRes{}, fmt.Errorf("Invalid Login Type Request")
}
