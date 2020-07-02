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
func (d *DiaryRepository) save(diary Diary) error {
	fmt.Println("Diary ID ==> ",diary.Id)
	insertQuery := "insert into diary(id, title, body, user_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6) returning id"
	err := d.db.QueryRow(insertQuery,
		diary.Id, diary.Title, diary.Body, diary.UserId, diary.CreatedAt, diary.UpdatedAt).Scan(&diary.Id)

	if err != nil {
		panic(err)
	}

	return nil
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
func (d *DiaryRepository) findDiaryById(diaryId string) (*sql.Row, error) {
	selectQuery := "select * from diary where id=$1"
	var diary Diary
	row := d.db.QueryRow(selectQuery, diaryId)
	err := row.Scan(&diary.Id, &diary.Title, &diary.Body, &diary.CreatedAt, &diary.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return row, nil
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
