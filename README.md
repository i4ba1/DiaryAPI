### Untuk perintah jalankan API nya dengan:
go build && ./DiaryAPI

### App Dockerize:
File Dockerfile, yang digunakan untuk dockerize app. Namun saya tidak build karena cukup besar jadi untuk upload nya. Perintah:
sudo docker build . -t go-dock

### Dockerize PostgreSQL
Pada folder postgresql-docker terdapat 2 file, docker-compose.yml dan backup untuk restore database. Perintah jalankan PostgreSQL:
sudo docker-compose up

### List rest api, base url: http://localhost:8789
1. /api/user/signUp
2. /api/user/login
3. /api/diary/createNewDiary
4. /diary/updateDiary
5. /diary/getDiaryByYearAndQuarter/{year}/{quarter}
6. /diary/getDiaryById/{id}
7. /diary/getAllDiary
8. /diary/deleteDiary/{id}

### Saya juga telah berikan Unit Testing nya pada file main_test.go yang menjalankan setiap rest yang di buat
 
