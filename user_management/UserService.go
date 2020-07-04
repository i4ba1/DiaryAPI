package user_management

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type UserService struct {
	userRepository UserRepository
	client	*redis.Client
}

func ProvideUserService(repo UserRepository, redisClient *redis.Client) UserService {
	return UserService{
		userRepository: repo,
		client: redisClient,
	}
}

func (a *UserService) signIn(dto LoginDto) (Account, error) {
	return a.userRepository.findUserByUserNameOrEmailAndPass(dto)
}

func (a *UserService) findByUsername(username string) error {
	return a.userRepository.findUserByUserName(username)
}

func (a *UserService) CreateToken(userId int64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.ExpiredAt = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.ExpiredRt = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + strconv.Itoa(int(userId))

	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	//Create Token
	token := jwt.New(jwt.SigningMethodHS512)
	atClaims:= token.Claims.(jwt.MapClaims)
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.ExpiredAt
	//at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	fmt.Println("Access Token ===> ",td.AccessToken)

	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims:= token.Claims.(jwt.MapClaims)
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.ExpiredAt
	//rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	//fmt.Println("REFRESH_SECRET ===> ",os.Getenv("REFRESH_SECRET"))
	//fmt.Println(rt.SignedString([]byte(os.Getenv("REFRESH_SECRET"))))
	//fmt.Println()
	td.RefreshToken, err = token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	fmt.Println("")
	fmt.Println("Refresh Token ===> ",td.RefreshToken)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (a *UserService) CreateAuth(userId int64, td *TokenDetails) error {
	fmt.Println("====== CreateAuth ====== ",a.client)
	fmt.Println()
	at := time.Unix(td.ExpiredAt, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.ExpiredRt, 0)
	now := time.Now()

	fmt.Println("AccessUuid ====> ",td.AccessUuid, " at ",at, " accessUuid ",td.AccessUuid, " - ",strconv.Itoa(int(userId)))
	errAccess := a.client.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now))
	//fmt.Print("Error on auth access token ", errAccess)

	if errAccess != nil {
		return errAccess.Err()
	}

	errRefresh := a.client.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	fmt.Print("Error on auth refresh token ", errRefresh)
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (a *UserService) FetchAuth(authD *AccessDetails) (int64, error) {
	userid, err := a.client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	i := int64(userID)
	if authD.UserId != int64(userID) {
		return 0, errors.New("unauthorized")
	}
	return i, nil
}

func (a *UserService) DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := a.client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (a *UserService) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.AccessUuid, authD.UserId)
	//delete access token
	deletedAt, err := a.client.Del(authD.AccessUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := a.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}
