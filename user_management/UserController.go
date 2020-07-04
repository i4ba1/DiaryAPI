package user_management

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type UserAPI struct {
	userService UserService
}

func ProvideUserAPI(u UserService) UserAPI {
	return UserAPI{userService: u}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Body-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *UserAPI) SignInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/login")

	var loginDto LoginDto
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	config := &PasswordConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}

	fmt.Println("Username ==> ", loginDto.Username)
	hashPass, err := GeneratePassword(config, loginDto.Pass)
	err = a.userService.findByUsername(loginDto.Username)

	fmt.Println("Username is ",err)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Username not found, please make sure the Username is correct")
		return
	}

	loginDto.Pass = hashPass
	var Account, errLogin = a.userService.signIn(loginDto)
	if errLogin != nil {
		respondWithError(w, http.StatusNotFound, "Incorrect Username/Email or Password")
		return
	}

	ts, err := a.userService.CreateToken(Account.Id)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := a.userService.CreateAuth(Account.Id, ts)
	if saveErr != nil {
		respondWithError(w, http.StatusUnprocessableEntity, saveErr.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	respondWithJSON(w, http.StatusOK, tokens)
}

func (a *UserAPI) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var tokenDetail TokenDetails
	err := json.NewDecoder(r.Body).Decode(&tokenDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	refreshToken := tokenDetail.RefreshToken

	//verify the token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		fmt.Println("the error: ", err)
		respondWithError(w, http.StatusUnauthorized, "Refresh token was expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			respondWithJSON(w, http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		i := int64(userId)
		if err != nil {
			respondWithJSON(w, http.StatusUnprocessableEntity, err)
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := a.userService.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			respondWithJSON(w, http.StatusUnauthorized, delErr.Error())
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := a.userService.CreateToken(i)
		if createErr != nil {
			respondWithJSON(w, http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := a.userService.CreateAuth(i, ts)
		if saveErr != nil {
			respondWithJSON(w, http.StatusForbidden, saveErr)
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		respondWithJSON(w, http.StatusCreated, tokens)
	} else {
		respondWithJSON(w, http.StatusUnauthorized, "refresh expired")
	}
}

func (a *UserAPI) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func (a *UserAPI) VerifyToken(r *http.Request) (*jwt.Token, error) {
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

func (a *UserAPI) TokenValid(r *http.Request) error {
	token, err := a.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (a *UserAPI) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := a.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
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

func (a *UserAPI) Logout(w http.ResponseWriter, r *http.Request) {
	metadata, err := a.ExtractTokenMetadata(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	delErr := a.userService.DeleteTokens(metadata)
	if delErr != nil {
		respondWithError(w, http.StatusUnauthorized, delErr.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Successfully logged out")
}
