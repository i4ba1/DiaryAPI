package diary_management

import (
	"database/sql"
)

type DiaryService struct {
	DiaryRepository DiaryRepository
}

func ProvideDiaryService(dr DiaryRepository) DiaryService{
	return DiaryService{DiaryRepository: dr}
}

/**
Get all diaries
 */
func (d *DiaryService) getDiaries(from string, to string) ([]Diary, error) {
	return d.DiaryRepository.findDiaryByYearAndQuarter(from, to)
	//return errors.New("not implemented")
}

/**
Create new diary
 */
func (d *DiaryService) createNewDiary(diary Diary) error {
	return d.DiaryRepository.save(diary)
}

/**
Update existing diary
*/
func (d *DiaryService) updateDiary(updateDiary UpdateDiary) error{
	return d.DiaryRepository.updateDiary(updateDiary)
}

/**
Get single diary
 */
func (d *DiaryService) getDiaryById(diaryId string) (*sql.Row, error){
	return d.DiaryRepository.findDiaryById(diaryId)
}

/**
Delete diary based on id
 */
func (d *DiaryService) deleteDiary(diaryId string) error {
	return d.DiaryRepository.deleteDiaryById(diaryId)
}




