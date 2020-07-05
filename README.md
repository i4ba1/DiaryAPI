### Untuk perintah jalankan API nya dengan:
go build && ./DiaryAPI

### App Dockerize:
File Dockerfile, yang digunakan untuk dockerize app. Perintah:
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

### Automated Test
go test -run <Nama Method> -v. Test pada file main_test
 
