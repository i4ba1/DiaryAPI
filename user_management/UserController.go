package user_management

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	fmt.Println(err)
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

	respondWithJSON(w, http.StatusOK, Account)
}
