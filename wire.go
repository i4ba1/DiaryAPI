package main

import (
	"database/sql"
	"github.com/google/wire"
	"github.com/i4ba1/DiaryAPI/RestAPI/diary_management"
)

func initProductAPI(db *sql.DB) diary_management.DiaryAPI {
	wire.Build(diary_management.ProvideDiaryRepository, diary_management.ProvideDiaryService, diary_management.ProvideDiaryAPI)

	return diary_management.DiaryAPI{}
}