package restapi

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aminasadiam/ChatterBox/datalayer"
	"github.com/aminasadiam/ChatterBox/security/generator"
	"github.com/aminasadiam/ChatterBox/security/jwt"
	"github.com/aminasadiam/ChatterBox/security/password"
	"github.com/aminasadiam/ChatterBox/sender"
)

type ShopRestApiHandler struct {
	dbhandler datalayer.SQLhandler
}

func newShopRestApihandler(dbHandler datalayer.SQLhandler) *ShopRestApiHandler {
	return &ShopRestApiHandler{
		dbhandler: dbHandler,
	}
}

func (handler *ShopRestApiHandler) Login(w http.ResponseWriter, r *http.Request) {
	if jwt.CheckUserIslogedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	temp := template.Must(template.ParseFiles("./templates/Site/login.html", "./templates/base.html", "./templates/menu.html"))
	temp.ExecuteTemplate(w, "base", nil)
}
func (handler *ShopRestApiHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}

	email := r.FormValue("email")
	pass := r.FormValue("password")

	if email == "" {
		msg := datalayer.LoginMessage{
			Message: "لطفا ایمیل خود را وارد کنید",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/login.html", "./templates/base.html", "./templates/menu.html"))
		temp.ExecuteTemplate(w, "base", msg)
		return
	}

	if pass == "" {
		msg := datalayer.LoginMessage{
			Message: "لطفا رمز عبور خود را وارد کنید",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/login.html", "./templates/base.html"))
		temp.ExecuteTemplate(w, "base", msg)
		return
	}

	isExist := jwt.CheckEmailIsExist(email, handler.dbhandler)
	if isExist {
		msg := datalayer.LoginMessage{
			Message: "ایمیل وارد شده صحیح نمی باشد",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/login.html", "./templates/base.html", "./templates/menu.html"))
		temp.ExecuteTemplate(w, "base", msg)
		return
	}

	isActive := jwt.CheckUserIsActive(email, handler.dbhandler)
	if !isActive {
		msg := datalayer.LoginMessage{
			Message: "ایمیل وارد شده فعال نشده است",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/login.html", "./templates/base.html", "./templates/menu.html"))
		temp.ExecuteTemplate(w, "base", msg)
		return
	}

	ok := jwt.SignIn(w, jwt.LoginInfo{
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(pass),
	}, handler.dbhandler)

	if !ok {
		msg := datalayer.LoginMessage{
			Message: "ایمیل یا رمز عبور اشتباه است",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/login.html", "./templates/base.html", "./templates/menu.html"))
		temp.ExecuteTemplate(w, "base", msg)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (handler *ShopRestApiHandler) Logout(w http.ResponseWriter, r *http.Request) {
	jwt.Logout(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (handler *ShopRestApiHandler) Register(w http.ResponseWriter, r *http.Request) {
	if jwt.CheckUserIslogedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	temp := template.Must(template.ParseFiles("./templates/Site/register.html", "./templates/base.html", "./templates/menu.html"))
	temp.ExecuteTemplate(w, "base", nil)
}

func (handler *ShopRestApiHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	confirmPass := r.FormValue("confirm_password")

	emailIsExist := jwt.CheckEmailIsExist(email, handler.dbhandler)
	if emailIsExist {
		message := datalayer.RegisterMessage{
			Message: "ایمیل وارد شده تکراری می باشد",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/register.html", "./templates/base.html", "./templates/menu.html"))
		temp.ExecuteTemplate(w, "base", message)
		return
	}

	if pass != confirmPass {
		message := datalayer.RegisterMessage{
			Message: "Password and Confirm Password not match",
		}
		temp := template.Must(template.ParseFiles("./templates/Site/register.html", "./templates/base.html", "./templates/menu.html"))
		temp.ExecuteTemplate(w, "base", message)
		return
	}

	password_hash, err := password.HashPassword(pass)
	if err != nil {
		log.Fatalln(err)
		return
	}

	activeCode, err := generator.GenerateActivationCode()
	if err != nil {
		log.Fatalln(err)
		return
	}

	user := datalayer.User{
		Username:     username,
		Email:        email,
		Password:     password_hash,
		ActiveCode:   activeCode,
		RegisterDate: time.Now(),
		IsActive:     false,
		IsDelete:     false,
		IsAdmin:      false,
	}

	err = handler.dbhandler.AddUser(user)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// send activation email
	sender.SendActiveEmail(username, email, activeCode)

	tmp := template.Must(template.ParseFiles("./templates/Site/successregister.html", "./templates/base.html", "./templates/menu.html"))
	tmp.ExecuteTemplate(w, "base", nil)
}

func (handler *ShopRestApiHandler) SuccessRegister(w http.ResponseWriter, r *http.Request) {
	if jwt.CheckUserIslogedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	tmp := template.Must(template.ParseFiles("./templates/Site/successregister.html", "./templates/base.html", "./templates/menu.html"))
	tmp.ExecuteTemplate(w, "base", nil)
}

func (handler *ShopRestApiHandler) ActiveUser(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("activecode")
	if code == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err := handler.dbhandler.ActiveUser(code)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	tmp := template.Must(template.ParseFiles("./templates/Site/successactive.html", "./templates/base.html", "./templates/menu.html"))
	tmp.ExecuteTemplate(w, "base", nil)
}
