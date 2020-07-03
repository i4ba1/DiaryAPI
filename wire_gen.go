package main

import (
	"database/sql"
	"github.com/i4ba1/DiaryAPI/RestAPI/diary_management"
	"github.com/i4ba1/DiaryAPI/user_management"
)

func InitDiaryAPI(db *sql.DB) diary_management.DiaryAPI {
	productRepository 	:= diary_management.ProvideDiaryRepository(db)
	productService 		:= diary_management.ProvideDiaryService(productRepository)
	productAPI 			:= diary_management.ProvideDiaryAPI(productService)
	return productAPI
}

func InitUserAPI(db *sql.DB) user_management.UserAPI{
	userRepository 	:= user_management.ProvideUserRepository(db)
	userService 	:= user_management.ProvideUserService(userRepository)
	userAPI			:= user_management.ProvideUserAPI(userService)
	return userAPI
}
