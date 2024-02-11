package utills

import (
	"petpal-backend/src/configs"

	"github.com/golang-jwt/jwt/v4"
)

type LoginReq struct {
	LoginType string `json:logintype`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRes struct {
	accessToken string
	Type        string `json:logintype`
	Username    string `json:"username"`
}

var secretKey = configs.GetJWT_SECRET()

type JWTClaims struct {
	LoginType string `json:logintype`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func Login(req *LoginReq) (*LoginRes, error) {
	loginType := req.LoginType
	if loginType == "scvp" {

	} else if loginType == "user" {
		u, err := GetUserByEmail(req.Email)
		if err != nil {
			return &LoginUserRes{}, err
		}

		err = CheckPassword(req.Password, u.Password)
		if err != nil {
			return &LoginUserRes{}, err
		}
	} else {
		return &LoginUserRes{}, err
	}
	/*
		u, err := user.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return &LoginUserRes{}, err
		}

		err = utils.CheckPassword(req.Password, u.Password)
		if err != nil {
			return &LoginUserRes{}, err
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
			Username: u.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    strconv.Itoa(int(u.ID)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		})

		ss, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return &LoginUserRes{}, err
		}
		return &LoginUserRes{accessToken: ss, Username: u.Username, ID: strconv.Itoa(int(u.ID))}, nil
	*/
}
