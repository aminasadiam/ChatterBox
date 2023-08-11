package restapi

import (
	"log"
	"net/http"

	"github.com/aminasadiam/ChatterBox/datalayer"
	"github.com/gorilla/mux"
)

var (
	STATIC_DIR = "/static/"
)

func RunApi(endpoint string, db datalayer.SQLhandler) error {
	r := mux.NewRouter()

	log.Fatalf("Started at : %s\n", endpoint)
	return http.ListenAndServe(endpoint, r)
}

func runApiOnRouter(r *mux.Router, db datalayer.SQLhandler) {
	siteRoutes(r, db)

	r.Methods("GET").PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("./templates/assets/"))))
}

func siteRoutes(r *mux.Router, db datalayer.SQLhandler) {
	handler := newShopRestApihandler(db)

	//r.Methods("GET").Path("/").HandlerFunc(handler.Index)

	// Athentication User Routes
	r.Methods("GET").Path("/login").HandlerFunc(handler.Login)
	r.Methods("POST").Path("/login").HandlerFunc(handler.PostLogin)
	r.Methods("GET").Path("/logout").HandlerFunc(handler.Logout)

	r.Methods("GET").Path("/register").HandlerFunc(handler.Register)
	r.Methods("POST").Path("/register").HandlerFunc(handler.PostRegister)
	r.Methods("GET").Path("/success").HandlerFunc(handler.SuccessRegister)
	r.Methods("GET").Path("/active/user/{activecode}").HandlerFunc(handler.ActiveUser)
}
