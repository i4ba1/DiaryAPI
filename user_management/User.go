package user_management

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
)

type Account struct {
	Id       int64		`json:"id"`
	Name     string    	`json:"name"`
	SureName string    	`json:"sure_name"`
	Email    string    	`json:"Email"`
	Username string    	`json:"Username"`
	Password string    	`json:"password"`
}

type CreateNewUser struct {
	Id       int64		`json:"id"`
	Name     string     `json:"name"`
	SureName string     `json:"sure_name"`
	Email    string     `json:"Email"`
	Username string     `json:"Username"`
	Password string     `json:"password"`
}

type NewUserDto struct {
	Id       int64
	Name     string
	SureName string
	Email    string
	Username string
	Password string     `validate:"min=6,max=32,regexp=^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[*.!@$%^&(){}[]:;<>,.?/~_+-=|]).{6,32}$"`
}

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type LoginDto struct {
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Pass     string `json:"Pass"`
}

type AccessDetails struct {
	AccessUuid string
	UserId   int64
}

type TokenDetails struct {
	AccessToken  string		`json:"access_token"`
	RefreshToken string		`json:"refresh_token"`
	AccessUuid   string		`json:"access_uuid"`
	RefreshUuid  string		`json:"refresh_uuid"`
	ExpiredAt    int64		`json:"expired_at"`
	ExpiredRt    int64		`json:"expired_rt"`
}

/**
Get single account or account detail
*/
func (u *Account) getUser(db *sql.DB) error {
	return db.QueryRow("select id, name, sure_name, Email, Username from account where id=$1", u.Id).
		Scan(&u.Id, &u.Name, &u.SureName, &u.Email, &u.Username)
	//return errors.New("not implemented")
}

func (u *Account) EmailWasUsed(db *sql.DB) bool{
	err, _ := db.Query("select Email from account where Email=$1", u.Email)
	fmt.Println("Error ==> ",err.Err())
	if err != nil {
		return true
	}
	return false
}

/**
Update account data based on id
*/
func (u *Account) updateUser(db *sql.DB) error {
	_, err := db.Exec("update account set name=$1, sure_name=$2, Email=$3, Username=$4, password=$5 where u.id=$6", u.Name, u.SureName, u.Email)
	return err
	//return errors.New("not implemented")
}

/**
delete selected account based on id
*/
func (u *Account) deleteUser(db *sql.DB) error {
	_, err := db.Exec("delete from account where u.id=$1", u.Id)
	return err
	//return errors.New("not implemented")
}

// GeneratePassword is used to generate a new password hash for storing and
// comparing at a later date.
func GeneratePassword(c *PasswordConfig, password string) (string, error) {

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return b64Hash, nil
}

/**
* Registration for new account
 */
func (u *Account) CreateUser(db *sql.DB) error {
	config := &PasswordConfig{
		time: 1,
		memory: 64 * 1024,
		threads: 4,
		keyLen: 32,
	}
	hashPass, err := GeneratePassword(config, u.Password)
	if err == nil{
		u.Password = hashPass
		err = db.QueryRow("insert into account(id, name, sure_name, Email, Username, password) values($1, $2, $3, $4, $5, $6) returning id",
			u.Id, u.Name, u.SureName, u.Email, u.Username, u.Password).Scan(&u.Id)
	}
	fmt.Println(hashPass)

	if err != nil {
		return err
	}
	return nil
	//return errors.New("not implemented")
}

/**
Get list of all users
*/
func (u *Account) getUsers(db *sql.DB, start, count int) ([]Account, error) {
	rows, err := db.Query("select id, name, sure_name, Email, Username from account limit $1 offset $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []Account

	for rows.Next() {
		var u Account
		if err := rows.Scan(&u.Id, &u.Name, &u.SureName, &u.Email, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
	//return errors.New("not implemented")
}
