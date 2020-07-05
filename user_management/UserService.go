package user_management

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type UserService struct {
	userRepository UserRepository
	client         *redis.Client
}

func ProvideUserService(repo UserRepository, redisClient *redis.Client) UserService {
	return UserService{
		userRepository: repo,
		client:         redisClient,
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
	token := jwt.New(jwt.SigningMethodHS256)
	atClaims := token.Claims.(jwt.MapClaims)
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.ExpiredAt
	//at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	fmt.Println("Access Token ===> ", td.AccessToken)

	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := token.Claims.(jwt.MapClaims)
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.ExpiredAt
	//rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	//fmt.Println("REFRESH_SECRET ===> ",os.Getenv("REFRESH_SECRET"))
	//fmt.Println(rt.SignedString([]byte(os.Getenv("REFRESH_SECRET"))))
	//fmt.Println()
	td.RefreshToken, err = token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	/*fmt.Println("")
	fmt.Println("Refresh Token ===> ", td.RefreshToken)
	fmt.Println("")*/

	if err != nil {
		//fmt.Println("Error Refresh Token ===> ", td.RefreshToken)
		fmt.Println("")
		return nil, err
	}
	/*fmt.Println("Token Details AccessUuid ==> ",td.AccessUuid)
	fmt.Println("Token Details RefreshToken ==> ",td.RefreshToken)
	fmt.Println("Token Details AccessToken ==> ",td.AccessToken)
	fmt.Println("Token Details RefreshUuid ==> ",td.RefreshUuid)
	fmt.Println("Token Details ExpiredAt ==> ",td.ExpiredAt)
	fmt.Println("Token Details ExpiredRt ==> ",td.ExpiredRt)*/
	return td, nil
}

func (a *UserService) CreateAuth(userId int64, td *TokenDetails) error {
	fmt.Println("====== CreateAuth ====== ", a.client)
	fmt.Println()
	at := time.Unix(td.ExpiredAt, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.ExpiredRt, 0)
	now := time.Now()

	//fmt.Println("AccessUuid ====> ", td.AccessUuid, " at ", at, " accessUuid ", td.AccessUuid, " - ", strconv.Itoa(int(userId)))
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

func (a *UserService) FetchAuth(authD *AccessDetails, redisClient *redis.Client) (int64, error) {
	fmt.Println("FetchAuth() ===> ",a.client)
	if a.client == nil {
		a.client = redisClient
	}

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

func (a *UserService) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := a.VerifyToken(r)
	//fmt.Println("ExtractTokenMetadata ===> ",err, "token ==> ",token)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	//fmt.Println("Token is valid ===> ",token.Valid)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     int64(userId),
		}, nil
	}
	return nil, err
}

/**
Parse, validate, and return a token.
keyFunc will receive the parsed token and should return the key for validating.
*/
func (a *UserService) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := a.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *UserService) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (a *UserService) TokenValid(r *http.Request) error {
	token, err := a.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
