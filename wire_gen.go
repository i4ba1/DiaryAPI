package main

import (
	"database/sql"
	"github.com/i4ba1/DiaryAPI/RestAPI/diary_management"
)

func InitDiaryAPI(db *sql.DB) diary_management.DiaryAPI {
	productRepository := diary_management.ProvideDiaryRepository(db)
	productService := diary_management.ProvideDiaryService(productRepository)
	productAPI := diary_management.ProvideDiaryAPI(productService)
	return productAPI
}
