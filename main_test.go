package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/i4ba1/DiaryAPI/RestAPI"
	"github.com/i4ba1/DiaryAPI/RestAPI/diary_management"
	"time"

	//"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M){
	a.Initialize()

	//ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	a.DB.Exec("DELETE FROM account")
	//a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func clearTableDiary(){
	a.DB.Exec("Delete from diary")
}

func ensureTableExists(){
	if _, err := a.DB.Exec(userTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func TestEmptyTable(t *testing.T){
	clearTable()
	req,_ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(request *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, request)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateNewDiary(t *testing.T){
	fmt.Println("/CreateNew Diary")
	diary, _ := json.Marshal(diary_management.Diary{
		Id:        uuid.New(),
		Title:     "List kegiatan",
		Body:      "Baca quran",
		UserId:    "3bc60588-5780-4254-8ed1-29c1bde55e13",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	request, err := http.NewRequest("POST", "/createNewDiary", bytes.NewReader(diary))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Body-Type", "application/json")
	response := executeRequest(request)
	fmt.Println("Response =====> "+response.Body.String())
	checkResponseCode(t, http.StatusOK, response.Code)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v caused by %v",
			status, http.StatusOK, response.Body.String())
	}
}

func TestCreateUser(t *testing.T) {

	clearTable()

	//	var jsonStr = []byte(`{"name":"uways", "sure_name": "muhammad uways", "email":"muhammaduways@outlook.co.id","username":"uways","password":"123","salt":"123","locked":false, "disabled":true}`)
	fmt.Println("UUID ===> ", uuid.New())
	payload, _ := json.Marshal(RestAPI.CreateNewUser{
		Id:       uuid.Must(uuid.NewRandom()),
		Name:     "Uways",
		SureName: "Muhammad Uways",
		Email:    "muhammaduways12@outlook.co.id",
		Username: "uways12",
		Password: "123",
		Locked:   false,
		Disabled: true,
	})
	//"username":"uways","password":"123","salt":"123","locked":false, "disabled":true
	request, err := http.NewRequest("POST", "/signUp", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Body-Type", "application/json")
	response := executeRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v caused by %v",
			status, http.StatusOK, response.Body.String())
	}

	fmt.Println("Response =====> "+response.Body.String())
	/*expected := `{"name":"Uways","sure_name":"Muhammad Uways","email":"muhammaduways@outlook.co.id","username":"uways","password":"123","salt":"123","locked":false,"disabled":true}`
	fmt.Println("Response =====> "+response.Body.String())
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expected)
	}*/
}

const userTableCreationQuery = `CREATE TABLE IF NOT EXISTS public."account"
(
    id uuid DEFAULT uuid_generate_v4 (),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    sure_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    email character(100) COLLATE pg_catalog."default" NOT NULL,
    password character(1) COLLATE pg_catalog."default" NOT NULL,
    salt text COLLATE pg_catalog."default" NOT NULL,
    locked boolean NOT NULL,
    disabled boolean NOT NULL,
    username character(100) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT pk_id PRIMARY KEY (id),
    CONSTRAINT unique_email UNIQUE (email)
        INCLUDE(email),
    CONSTRAINT unique_username UNIQUE (username)
)`