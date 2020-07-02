package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/i4ba1/DiaryAPI/RestAPI"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"regexp"
)

type App struct {
	Router 	*mux.Router
	DB 		*sql.DB
}

func (a *App) Initialize(){
	a.DB = RestAPI.InitDB()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8789", a.Router))
}

/**
New user to signUp
 */
func (a *App) signUp(w http.ResponseWriter, r *http.Request) {
	var u RestAPI.Account
	err := json.NewDecoder(r.Body).Decode(&u)
	if  err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/**
	Check if email was used or was registered in Database
	 */
	if !u.EmailWasUsed(a.DB) {
		respondWithError(w, http.StatusConflict, "Email was used")
		return
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(u.Email) {
		respondWithError(w, http.StatusExpectationFailed, "Incorrect email pattern")
		return
	}

	passwordRegex := regexp.MustCompile("^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[@#$%^&-+=()])(?=\\\\S+$).{8, 20}$")
	if !passwordRegex.MatchString(u.Password) {
		respondWithError(w, http.StatusExpectationFailed, "Incorrect password pattern, should contain 6-32 characters and must\nhave at least one uppercase letter, one lowercase letter, one number, and one special character")
		return
	}

	defer r.Body.Close()

	if err := u.CreateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Body-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) initializeRoutes() {
	a.Router = mux.NewRouter()
	a.Router.StrictSlash(true)
	diaryAPI := InitDiaryAPI(a.DB)
	subRouter := a.Router.PathPrefix("/api").Subrouter()

	//a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	subRouter.HandleFunc("/signUp", a.signUp).Methods("POST")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	//Diary
	//subRouter.HandleFunc("/getAllDiary", diaryAPI.GetAllDiary).Methods("GET")
	subRouter.HandleFunc("/diary/createNewDiary", diaryAPI.CreateNewDiaryHandler).Methods("POST")
	subRouter.HandleFunc("/diary/updateDiary", diaryAPI.UpdateDiaryHandler).Methods("POST")
	subRouter.HandleFunc("/diary/GetDiaryByYearAndQuarter/{year}/{quarter}", diaryAPI.GetDiaryByYearAndQuarter).Methods("GET")
	subRouter.HandleFunc("/diary/getDiaryById/{id}", diaryAPI.GetDiaryById).Methods("GET")
	subRouter.HandleFunc("/diary/deleteDiary/{id}", diaryAPI.DeleteDiary).Methods("GET")
}