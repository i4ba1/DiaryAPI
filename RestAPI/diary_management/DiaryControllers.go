package diary_management

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"

	//"time"
)

type DiaryAPI struct {
	DiaryService DiaryService
}

func ProvideDiaryAPI(d DiaryService) DiaryAPI{
	return DiaryAPI{DiaryService: d}
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

/**
Create new dairy note, /createNewDiary, with POST method
 */
func (d *DiaryAPI) CreateNewDiaryHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("/createNewDiary")

	var diary Diary
	diary.Id = uuid.NewV4()
	diary.UpdatedAt = time.Now()
	diary.CreatedAt = time.Now()
	err := json.NewDecoder(r.Body).Decode(&diary)
	if  err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println(diary.Id)
	if err := d.DiaryService.createNewDiary(diary); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, diary)
}

/**
Rest handler to update Diary /updateDiary with POST method
 */
func (d *DiaryAPI) UpdateDiaryHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("/updateDiaryHandler")
	var updateDiary UpdateDiary
	updateDiary.UpdatedAt = time.Now()

	err := json.NewDecoder(r.Body).Decode(&updateDiary)
	if  err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
func (d *DiaryAPI) GetDiaryByYearAndQuarter(w http.ResponseWriter, r *http.Request){
	//t := time.Now()

	var quarterDto GetQuarterlyDiary
	err := json.NewDecoder(r.Body).Decode(&quarterDto)
	if  err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quarter := quarterDto.Quarter
	var from = ""
	var to = ""

	if quarter == 1 {
		from = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
		to = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter+2) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
	}else if quarter == 2{
		quarter += 2
		from = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
		to = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter+2) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
	}else if quarter == 3{
		quarter += 4
		from = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
		to = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter+2) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
	}else if quarter == 4{
		quarter += 6
		from = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
		to = strconv.Itoa(quarterDto.Year) + "-" + strconv.Itoa(quarter+2) +"-"+"01"+" "+"00"+":"+"00"+":"+"00"
	}

	diaries, err := d.DiaryService.getDiaries(from, to)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, diaries)
}

/**
Find diary by diaryId, /findDiaryById/{id}, GET method
 */
func (d *DiaryAPI) GetDiaryById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	diaryId := params["id"]
	if len(diaryId) == 0 || diaryId == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid diary ID")
		return
	}

	//d := diary{ID: id}
	row, err2 := d.DiaryService.getDiaryById(diaryId)
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
func (d *DiaryAPI) DeleteDiary(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	diaryId := params["id"]
	if len(diaryId) == 0 || diaryId == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	err := d.DiaryService.deleteDiary(diaryId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


