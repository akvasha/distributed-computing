package main

import (
	"DC-homework-1/authentication/dbClient"
	"DC-homework-1/authentication/mqClient"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	ACCESSTOKEN       int = 0
	REFRESHTOKEN      int = 1
	CONFIRMATIONTOKEN int = 2
)

type Config struct {
	TokenLength               uint64
	AccessTokenLifetime       time.Duration
	RefreshTokenLifetime      time.Duration
	ConfirmationTokenLifetime time.Duration
	ConfirmationAddress       string
}

var db dbClient.Client
var config Config
var mq mqClient.MQClient

type errorResponse struct {
	Error string `json:"error"`
}

func responseError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintln(err)})
}

func (c *Config) Init() (err error) {
	if c.TokenLength, err = strconv.ParseUint(os.Getenv("TOKEN_LENGTH"), 10, 64); err != nil {
		return
	}
	if c.AccessTokenLifetime, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_LIFETIME")); err != nil {
		return
	}
	if c.RefreshTokenLifetime, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_LIFETIME")); err != nil {
		return
	}
	if c.ConfirmationTokenLifetime, err = time.ParseDuration(os.Getenv("CONFIRMATION_TOKEN_LIFETIME")); err != nil {
		return
	}
	c.ConfirmationAddress = os.Getenv("CONFIRM_ADDRESS")
	return
}

func initToken(len uint64) string {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func createHash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 17)
	return string(bytes), err
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	hash, err := createHash(req.Password)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
	}
	user := dbClient.User{
		Username:    req.Username,
		Password:    hash,
		Email:       req.Email,
		Phone:       req.Phone,
		PhoneStatus: false,
	}
	if err = db.AddUser(user); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if err = sendPhoneConfirmation(user.Username, user.Phone); err != nil {
		responseError(w, err, http.StatusInternalServerError)
	}
	return
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	user, err := db.GetUser(req.Username)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		responseError(w, err, http.StatusUnauthorized)
		return
	}
	if !user.PhoneStatus {
		if err = sendPhoneConfirmation(user.Username, user.Phone); err != nil {
			responseError(w, err, http.StatusInternalServerError)
			return
		}
		responseError(w, errors.New("Phone is not confirmed, we have sent you new token, please confirm your phone."+
			"Maybe your token has expired."), http.StatusForbidden)
		return
	}
	var isDuplicate = true
	var accessToken, refreshToken string
	for isDuplicate {
		isDuplicate = false
		accessToken, refreshToken = initToken(config.TokenLength), initToken(config.TokenLength)
		accessTokenData := dbClient.TokenData{
			Token:    accessToken,
			Type:     ACCESSTOKEN,
			Lifetime: time.Now().Add(config.AccessTokenLifetime),
			Username: req.Username,
		}
		refreshTokenData := dbClient.TokenData{
			Token:    refreshToken,
			Type:     REFRESHTOKEN,
			Lifetime: time.Now().Add(config.RefreshTokenLifetime),
			Username: req.Username,
		}
		if err = db.AddToken(accessTokenData); err != nil {
			if err == dbClient.ErrorDuplicateToken {
				isDuplicate = true
			} else {
				responseError(w, err, http.StatusInternalServerError)
				return
			}
		}
		if err = db.AddToken(refreshTokenData); err != nil {
			if err == dbClient.ErrorDuplicateToken {
				isDuplicate = true
			} else {
				responseError(w, err, http.StatusInternalServerError)
				return
			}
		}
	}
	_ = json.NewEncoder(w).Encode(SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func sendPhoneConfirmation(username string, phone string) (err error) {
	var isDuplicate = true
	var confirmationToken string
	for isDuplicate {
		isDuplicate = false
		confirmationToken = initToken(config.TokenLength)
		confirmationTokenData := dbClient.TokenData{
			Token:    confirmationToken,
			Type:     CONFIRMATIONTOKEN,
			Lifetime: time.Now().Add(config.ConfirmationTokenLifetime),
			Username: username,
		}
		if err = db.AddToken(confirmationTokenData); err != nil {
			if err == dbClient.ErrorDuplicateToken {
				isDuplicate = true
			} else {
				return
			}
		}
	}
	reg := fmt.Sprintf("Confrim registration via token: %s", confirmationToken)
	err = mq.SendMessage(mqClient.Message{
		Receiver: phone,
		Text:     reg,
	})
	return
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateAccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func updateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	var token dbClient.TokenData
	var err error
	if token, err = db.GetToken(req.RefreshToken); err != nil {
		if err == dbClient.ErrNotFound {
			responseError(w, errors.New("Invalid token"), http.StatusUnauthorized)
		} else {
			responseError(w, err, http.StatusInternalServerError)
		}
		return
	}
	if token.Type != REFRESHTOKEN {
		responseError(w, errors.New("Provide refreshToken"), http.StatusUnauthorized)
	}
	if token.Lifetime.Before(time.Now()) {
		responseError(w, errors.New("Expired token"), http.StatusUnauthorized)
	}
	accessToken, refreshToken := initToken(config.TokenLength), initToken(config.TokenLength)
	accessTokenData := dbClient.TokenData{
		Token:    accessToken,
		Lifetime: time.Now().Add(config.AccessTokenLifetime),
		Username: token.Username,
	}
	refreshTokenData := dbClient.TokenData{
		Token:    refreshToken,
		Lifetime: time.Now().Add(config.RefreshTokenLifetime),
		Username: token.Username,
	}
	if err = db.AddToken(accessTokenData); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if err = db.AddToken(refreshTokenData); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(UpdateAccessResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func getAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var token dbClient.TokenData
	var err error
	if token, err = db.GetToken(r.Header.Get("auth")); err != nil {
		if err == dbClient.ErrNotFound {
			responseError(w, errors.New("Invalid token"), http.StatusUnauthorized)
		} else {
			responseError(w, err, http.StatusInternalServerError)
		}
		return
	}
	if token.Type != ACCESSTOKEN {
		responseError(w, errors.New("Provide accessToken"), http.StatusUnauthorized)
		return
	}
	if token.Lifetime.Before(time.Now()) {
		responseError(w, errors.New("Expired token"), http.StatusUnauthorized)
		return
	}
}

type confirmationResponse struct {
	Message string `json:"message"`
}

func confirmationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, ok := params["token"]; ok {
		var token dbClient.TokenData
		var err error
		if token, err = db.GetToken(params["token"]); err != nil {
			if err == dbClient.ErrNotFound {
				responseError(w, errors.New("Invalid token"), http.StatusUnauthorized)
			} else {
				responseError(w, err, http.StatusInternalServerError)
			}
			return
		}
		if token.Type != CONFIRMATIONTOKEN {
			responseError(w, errors.New("Provide confirmationToken"), http.StatusUnauthorized)
			return
		}
		if token.Lifetime.Before(time.Now()) {
			responseError(w, errors.New("Expired token"), http.StatusUnauthorized)
			return
		}
		if err = db.ConfirmPhone(token.Username); err != nil {
			responseError(w, err, http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(confirmationResponse{Message: "Registration completed"})
		return
	} else {
		responseError(w, errors.New("Provide token"), http.StatusBadRequest)
		return
	}
}

func main() {
	var err error
	if err = config.Init(); err != nil {
		log.Fatal(err)
	}
	if db, err = dbClient.InitClient(); err != nil {
		log.Fatal(err)
	}
	if mq, err = mqClient.InitMQClient(); err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/signUp", signUpHandler).Methods("POST")
	r.HandleFunc("/signIn", signInHandler).Methods("POST")
	r.HandleFunc("/refresh", updateAccessTokenHandler).Methods("PUT")
	r.HandleFunc("/validate", getAccessTokenHandler).Methods("GET")
	r.HandleFunc("/confirm/{token}", confirmationHandler).Methods("GET")
	log.Println("Authentication server started")
	log.Fatal(http.ListenAndServe(":8000", r))
}
