package jwt

import (
	"log"
	"net/http"
	"time"

	"github.com/aminasadiam/ChatterBox/datalayer"
	"github.com/aminasadiam/ChatterBox/security/password"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("ThisIsAOnlineShop")

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Cliams struct {
	Username string
	jwt.StandardClaims
}

func SignIn(w http.ResponseWriter, login LoginInfo, dbhandler datalayer.SQLhandler) bool {
	user, err := dbhandler.GetUserByEmail(login.Email)
	if err != nil {
		return false
	}

	if ok := password.CheckComparePass(login.Password, user.Password); !ok {
		return false
	}

	expTime := time.Now().Add(time.Hour * 24 * 60)

	claims := Cliams{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return false
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expTime,
	})

	return true
}

func CheckUserIslogedIn(r *http.Request) bool {
	c, err := r.Cookie("token")
	if err != nil {
		return false
	}

	token, err := jwt.ParseWithClaims(c.Value, &Cliams{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false
	}

	if _, ok := token.Claims.(*Cliams); ok && token.Valid {
		return true
	}

	return false
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		return
	}

	c.MaxAge = -1

	http.SetCookie(w, c)
}

func CheckEmailIsExist(email string, handler datalayer.SQLhandler) bool {
	user, err := handler.GetUserByEmail(email)
	if err != nil {
		log.Fatal(err)
	}

	niluser := datalayer.User{}
	return user != niluser
}

func CheckUserIsActive(email string, handler datalayer.SQLhandler) bool {
	user, err := handler.GetUserByEmail(email)
	if err != nil {
		log.Fatal(err)
	}

	if user.IsActive {
		return true
	} else {
		return false
	}
}

func GetUsername(r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		return ""
	}

	token, err := jwt.ParseWithClaims(c.Value, &Cliams{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return ""
	}

	if claims, ok := token.Claims.(*Cliams); ok && token.Valid {
		return claims.Username
	}

	return ""
}
