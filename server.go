package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"encoding/json"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

// client
var (
	AppID     = "APP01"
	AppSecret = "APPSEC"
)

// User from Json
type User struct {
	Id       string `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Read file
	bytes, err := ioutil.ReadFile("user.json")
	if err != nil {
		log.Fatal(err)
	}

	// Decode
	var users []User
	if err := json.Unmarshal(bytes, &users); err != nil {
		log.Fatal(err)
	}

	// Token Store
	manager := manage.NewDefaultManager()

	// SetPasswordTokenCfg set the password grant token config
	cfg := &manage.Config{
		// access token expiration time (default: time.Hour * 2)
		AccessTokenExp: time.Minute * 2,
		// refresh token expiration time (default: time.Hour * 24 * 7)
		RefreshTokenExp: time.Hour * 24 * 7,
		// whether to generate the refreshing token (default: true)
		IsGenerateRefresh: false,
	}
	manager.SetPasswordTokenCfg(cfg)

	// TODO: imple RDBMS
	manager.MustTokenStorage(store.NewFileTokenStore("token.db"))

	// client memory store
	clientStore := store.NewClientStore()
	// form JSON
	for _, p := range users {
		clientStore.Set(AppID, &models.Client{ID: p.Id, Secret: AppSecret})
		//fmt.Println(p.Id)
		//fmt.Println(p.Name)
		//fmt.Println(p.Password)
	}
	manager.MapClientStorage(clientStore)

	// http srv
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetAllowedGrantType(oauth2.PasswordCredentials)

	// client authentication from GET pram
	// e.g. &client_id=APP01&client_secret=APPSEC
	srv.SetClientInfoHandler(server.ClientFormHandler)

	// password authentication from GET pram
	// e.g. &username=admin&password=123456
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		//if username == "admin" && password == "123456" {
		//	userID = "000000"
		//	return
		//}
		//err = fmt.Errorf("user not found")
		//return

		userID = ""
		for _, p := range users {
			if username == p.Name && password == p.Password {
				userID = p.Id
				return
			}
		}
		if userID == "" {
			err = fmt.Errorf("user not found")
		}
		return
	})

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		token, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	})

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}
