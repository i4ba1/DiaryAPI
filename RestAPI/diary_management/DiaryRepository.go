package diary_management

import (
	"database/sql"
	"fmt"
	//"github.com/google/uuid"
	//"time"
)

type DiaryRepository struct {
	db *sql.DB
}

func ProvideDiaryRepository(db *sql.DB) DiaryRepository{
	repository := DiaryRepository{db:db}
	return repository
}

/**
Create/Save new diary
 */
func (d *DiaryRepository) save(diary Diary) (int64, error) {
	fmt.Println("Diary ID ==> ",diary.Id)
	var Id int64
	Id = 1
	var curretDateTime string
	queryCurretDateTime := "select current_date + current_time as created_at"
	rows, errQuery := d.db.Query(queryCurretDateTime)

	if errQuery != nil {
		return 0, nil
	}

	for rows.Next() {
		if err := rows.Scan(&curretDateTime); err != nil{
			return 0, err
		}
	}

	insertQuery := "insert into diary(title, body, user_id, created_at, updated_at) values($1,$2,$3,$4,$5) returning id"
	row := d.db.QueryRow(insertQuery,
		diary.Title, diary.Body, diary.UserId, curretDateTime, curretDateTime)
	fmt.Println("row ===> ",row)
	err := row.Scan(&Id)

	if err != nil {
		return 0, err
	}

	return Id, nil
}

func (d *DiaryRepository) findDiaryByYearAndQuarter(from string, to string) ([]Diary, error){
	rows, err := d.db.Query("select * from diary where updated_at >= $1 and updated_at < $2", from, to)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var diaries []Diary

	for rows.Next() {
		var d Diary
		if err := rows.Scan(&d.Id, &d.Title, &d.Body, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		diaries = append(diaries, d)
	}

	return diaries, nil
}

/**
Get all diary
 */
func (d *DiaryRepository) findAll() ([]Diary, error){
	selectAll := "select * from diary order by created_at asc"
	rows, err := d.db.Query(selectAll)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diaries []Diary
	for rows.Next() {
		var diary Diary
		if err := rows.Scan(&diary.Id, &diary.Title, &diary.Body, &diary.UserId, &diary.UpdatedAt, &diary.CreatedAt); err != nil{
			return nil, err
		}
		diaries = append(diaries, diary)
	}

	fmt.Println("Size of list diary ===> ",len(diaries))
	return diaries, nil
}

/**
Update diary data title, body, updated_at based on diaryId
 */
func (d *DiaryRepository) updateDiary(diary UpdateDiary) error{
	fmt.Println("Diary ID ==> ",diary.Id)
	updateQuery := "update diary set title=$1, body=$2, updated_at=$3 where id=$4"
	_, err := d.db.Exec(updateQuery,
		diary.Title, diary.Body, diary.UpdatedAt, diary.Id)

	return err
}

/**
Find diary by diaryId
 */
func (d *DiaryRepository) findDiaryById(diaryId int64) (Diary, error) {
	selectQuery := "select * from diary where id=$1"
	var diary Diary
	row := d.db.QueryRow(selectQuery, diaryId)
	err := row.Scan(&diary.Id, &diary.Title, &diary.Body, &diary.UserId, &diary.CreatedAt, &diary.UpdatedAt)

	if err != nil {
		return Diary{}, err
	}

	return diary, nil
}

/**
Delete diary by diaryId
 */
func (d *DiaryRepository) deleteDiaryById(diaryId string) error{
	getCurrentDiary := "select * from diary where id=$1"
	var diary Diary

	row := d.db.QueryRow(getCurrentDiary, diaryId)
	deleteQuery := "delete from diary where id=$1"
	var err error
	if row != nil {
		_, err = d.db.Exec(deleteQuery, diary.Id)
	}

	if err != nil{
		return err
	}

	return err
}
