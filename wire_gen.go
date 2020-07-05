package main

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"github.com/i4ba1/DiaryAPI/RestAPI/diary_management"
	"github.com/i4ba1/DiaryAPI/user_management"
)

func InitDiaryAPI(db *sql.DB, client *redis.Client) diary_management.DiaryAPI {
	diaryRepository		:= diary_management.ProvideDiaryRepository(db)
	diaryService		:= diary_management.ProvideDiaryService(diaryRepository, client)
	diaryAPI 			:= diary_management.ProvideDiaryAPI(diaryService)
	return diaryAPI
}

func InitUserAPI(db *sql.DB, client *redis.Client) user_management.UserAPI{
	userRepository 	:= user_management.ProvideUserRepository(db)
	userService 	:= user_management.ProvideUserService(userRepository, client)
	userAPI			:= user_management.ProvideUserAPI(userService)
	return userAPI
}
