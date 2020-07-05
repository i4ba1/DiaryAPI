package diary_management

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/i4ba1/DiaryAPI/user_management"
	"net/http"
	"strconv"
	"time"

	//"time"
)


type DiaryAPI struct {
	DiaryService DiaryService
	UserService user_management.UserService
}

func ProvideDiaryAPI(d DiaryService) DiaryAPI {
	return DiaryAPI{
		DiaryService: d,
	}
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

func (d *DiaryAPI) isAuthorized(w http.ResponseWriter, r *http.Request){
	//Extract the access token metadata
	metadata, err := d.UserService.ExtractTokenMetadata(r)
	//fmt.Println("isAuthorized ===> ",err, " metadata ===> ",metadata)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = d.UserService.FetchAuth(metadata, d.DiaryService.redis)
	fmt.Println("FetchAuth err ===> ",err, " metadata ===> ",metadata)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	return
}

/**
Create new dairy note, /createNewDiary, with POST method
*/
func (d *DiaryAPI) CreateNewDiaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/createNewDiary")

	var diary Diary
	currentTime := time.Now()
	createdDate := currentTime.Format(layoutISO)
	diary.UpdatedAt = createdDate
	diary.CreatedAt = createdDate
	err := json.NewDecoder(r.Body).Decode(&diary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.isAuthorized(w, r)
	defer r.Body.Close()

	fmt.Println(diary.Id)
	if _, err := d.DiaryService.createNewDiary(diary); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, diary)
}

/**
Rest handler to update Diary /updateDiary with POST method
*/
func (d *DiaryAPI) UpdateDiaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/updateDiaryHandler")
	var updateDiary UpdateDiary

	err := json.NewDecoder(r.Body).Decode(&updateDiary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.isAuthorized(w, r)
	defer r.Body.Close()

	fmt.Println(updateDiary.Id)
	if err := d.DiaryService.updateDiary(updateDiary); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, updateDiary)
}

/**
Find all dairy, /GetDiaryByYearAndQuarter/{year}/{quarter} with GET method
*/
func (d *DiaryAPI) GetDiaryByYearAndQuarter(w http.ResponseWriter, r *http.Request) {
	//t := time.Now()

	params := mux.Vars(r)
	year, err 		:= 	strconv.Atoi(params["year"])
	quarter,_ 		:=	strconv.Atoi(params["quarter"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	var from = ""
	var to = ""

	d.isAuthorized(w, r)
	fmt.Println(quarter)
	if quarter == 1 {
		from = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter)) + "-" +"01"
		to = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter+2)) + "-" + "01"
	} else if quarter == 2 {
		quarter += 2
		from = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter)) + "-" +"01"
		to = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter+2)) + "-" + "01"
	} else if quarter == 3 {
		quarter += 4
		from = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter)) + "-" +"01"
		to = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter+2)) + "-" + "01"
	} else if quarter == 4 {
		quarter += 6
		from = strconv.Itoa(year) + "-" + ("0"+strconv.Itoa(quarter)) + "-" +"01"
		to = strconv.Itoa(year) + "-" + strconv.Itoa(quarter+2) + "-" + "01"
	}

	fmt.Println("From ===> ",from)

	diaries, err := d.DiaryService.getDiariesFilterByYearAndQuarter(from, to)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, diaries)
}

/**
Find diary by diaryId, /findDiaryById/{id}, GET method
*/
func (d *DiaryAPI) GetDiaryById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	diaryId, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid diary ID")
		return
	}

	d.isAuthorized(w, r)
	row, err2 := d.DiaryService.getDiaryById(int64(diaryId))
	if err2 != nil {
		switch err2 {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err2.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, row)
}

/**
Delete diary by diaryId, /deleteDiary/{id}, GET method
*/
func (d *DiaryAPI) DeleteDiary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	diaryId := params["id"]
	if len(diaryId) == 0 || diaryId == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	d.isAuthorized(w, r)
	err := d.DiaryService.deleteDiary(diaryId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

/**
Get all diaries, /getAllDiary, GET method
*/
func (d *DiaryAPI) GetAllDiary(w http.ResponseWriter, r *http.Request) {
	d.isAuthorized(w, r)
	diaries, err := d.DiaryService.GetDiaries()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, diaries)
}
