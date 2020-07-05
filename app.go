package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/i4ba1/DiaryAPI/RestAPI"
	"github.com/i4ba1/DiaryAPI/user_management"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"regexp"
)

type App struct {
	Router 	*mux.Router
	DB 		*sql.DB
	Client	*redis.Client
}

func (a *App) Initialize(){
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	fmt.Println("App, dsn ===> ",dsn)
	a.Client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := a.Client .Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Client ===> ",a.Client)

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
	var u user_management.Account
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

	emailRegex, errEmail := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if errEmail != nil {
		respondWithError(w, http.StatusExpectationFailed, "Incorrect email regex pattern")
		return
	}

	isEmailValid := emailRegex.MatchString(u.Email)
	if !isEmailValid {
		respondWithError(w, http.StatusExpectationFailed, "Incorrect email pattern")
		return
	}

	passLength := len(u.Password)
	if passLength < 6 || passLength > 32 {
		respondWithError(w, http.StatusExpectationFailed, "Incorrect password pattern, minimum 6 character and maximum 32 character" +
			"and one special character")
		return
	}

	if validPassword(u.Password) != nil{
		respondWithError(w, http.StatusExpectationFailed, "Incorrect password pattern, should contain at least one uppercase letter, one lowercase letter, one number, " +
			"and one special character")
		return
	}
	/*passwordRegex := regexp.MustCompile("^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[*.!@$%^&(){}[]:;<>,.?/~_+-=|]).{6,32}$")
	if !passwordRegex.MatchString(u.Password) {
		respondWithError(w, http.StatusExpectationFailed, "Incorrect password pattern, should contain 6-32 characters and must\nhave at least one uppercase letter, one lowercase letter, one number, and one special character")
		return
	}*/

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
	//go build && ./DiaryAPI
	a.Router = mux.NewRouter()
	a.Router.StrictSlash(true)
	diaryAPI 	:= InitDiaryAPI(a.DB, a.Client)
	userAPI		:= InitUserAPI(a.DB, a.Client)
	subRouter := a.Router.PathPrefix("/api").Subrouter()

	//a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	subRouter.HandleFunc("/user/signUp", a.signUp).Methods("POST")
	subRouter.HandleFunc("/user/login", userAPI.SignInHandler).Methods("POST")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	//Diary
	//subRouter.HandleFunc("/getAllDiary", diaryAPI.GetAllDiary).Methods("GET")
	subRouter.HandleFunc("/diary/createNewDiary", diaryAPI.CreateNewDiaryHandler).Methods("POST")
	subRouter.HandleFunc("/diary/updateDiary", diaryAPI.UpdateDiaryHandler).Methods("POST")
	subRouter.HandleFunc("/diary/getDiaryByYearAndQuarter/{year}/{quarter}", diaryAPI.GetDiaryByYearAndQuarter).Methods("GET")
	subRouter.HandleFunc("/diary/getDiaryById/{id}", diaryAPI.GetDiaryById).Methods("GET")
	subRouter.HandleFunc("/diary/getAllDiary", diaryAPI.GetAllDiary).Methods("GET")
	subRouter.HandleFunc("/diary/deleteDiary/{id}", diaryAPI.DeleteDiary).Methods("GET")
}