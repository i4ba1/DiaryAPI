package user_management

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func ProvideUserRepository(db *sql.DB) UserRepository {
	repository := UserRepository{db: db}
	return repository
}

/**
to query user by username or email and password
*/
func (d *UserRepository) findUserByUserNameOrEmailAndPass(dto LoginDto) (Account, error) {
	fmt.Println("======== findUserByUserNameOrEmailAndPass ======")
	queryLogin := "select * from account where username=$1 or email=$2 and password=$3"

	var a Account
	row := d.db.QueryRow(queryLogin, dto.Username, dto.Email, dto.Pass)
	var err = row.Scan(&a.Id, &a.Name, &a.SureName, &a.Email, &a.Password, &a.Username)
	return a, err
}

/**
Repository for select query account by username
*/
func (d *UserRepository) findUserByUserName(username string) error {
	queryLogin := "select * from account where username=$1"
	var a Account
	return d.db.QueryRow(queryLogin, username).Scan(&a.Id, &a.Name, &a.SureName, &a.Email, &a.Password, &a.Username)
}
