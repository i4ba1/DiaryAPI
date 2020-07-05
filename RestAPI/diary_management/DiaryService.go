package diary_management

import "github.com/go-redis/redis/v7"

type DiaryService struct {
	DiaryRepository DiaryRepository
	redis           *redis.Client
}

func ProvideDiaryService(dr DiaryRepository, client *redis.Client) DiaryService {
	return DiaryService{
		DiaryRepository: dr,
		redis: client,
	}
}

/**
Get diaries filter by year and quarter
*/
func (d *DiaryService) getDiariesFilterByYearAndQuarter(from string, to string) ([]Diary, error) {
	return d.DiaryRepository.findDiaryByYearAndQuarter(from, to)
	//return errors.New("not implemented")
}

/**
Get all diaries
*/
func (d *DiaryService) GetDiaries() ([]Diary, error) {
	return d.DiaryRepository.findAll()
	//return errors.New("not implemented")
}

/**
Create new diary
*/
func (d *DiaryService) createNewDiary(diary Diary) (int64, error) {
	return d.DiaryRepository.save(diary)
}

/**
Update existing diary
*/
func (d *DiaryService) updateDiary(updateDiary UpdateDiary) error {
	return d.DiaryRepository.updateDiary(updateDiary)
}

/**
Get single diary
*/
func (d *DiaryService) getDiaryById(diaryId int64) (Diary, error) {
	return d.DiaryRepository.findDiaryById(diaryId)
}

/**
Delete diary based on id
*/
func (d *DiaryService) deleteDiary(diaryId string) error {
	return d.DiaryRepository.deleteDiaryById(diaryId)
}
