package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hako/branca"
)

var (
	clientID            = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret        = os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURL         = os.Getenv("SIMPSONS_CONFIG_REDIRECT_URL")
	scope               = os.Getenv("SIMPSONS_CONFIG_SCOPES")
	superSecretPassword = os.Getenv("SIMPSONS_CONFIG_ENCRYPTING_SECRET")
	secureCookie        = os.Getenv("SIMPSONS_CONFIG_SECURE_COOKIE")
)

// UserInfo contains the main api object that gets returned to the front end, returns just your avatar url and login name in github
type UserInfo struct {
	LoginName string           `json:"login"`
	AvatarURL string           `json:"avatar_url"`
	Simpson   SimpsonCharacter `json:"simpson"`
}

// SimpsonCharacter contains the main information for a simpsons character that we care to present
type SimpsonCharacter struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Quote     string `json:"quote"`
}

func main() {

	mux := http.NewServeMux()
	// These handlers represent the API for our application, they initiate and consume requests meant for applications
	mux.HandleFunc("/api/start/github", loggingMiddleware(startOAuthWithGitHub))
	mux.HandleFunc("/api/callback/github", loggingMiddleware(handleCallbackFromGitHub))
	mux.HandleFunc("/api/userinfo", loggingMiddleware(handleUserInfo))

	// These handlers are all about our UI, they simply server up very basic HTML files from the file system, nothing fancy
	mux.HandleFunc("/login.html", loggingMiddleware(serveLogin))
	mux.HandleFunc("/", loggingMiddleware(serveIndex))

	// These are just some more handlers for the assets that the HTML files will request, makes it a bit easier to manage
	// Just know that all they do is simply serve any file located in the ./public/images/ and ./public/css/
	// For security reasons this is locked to just those
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images/"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css/"))))
	mux.Handle("/robots.txt", http.FileServer(http.Dir("./public/html/")))

	// The application starts and listens on port 8080
	fmt.Println("-> Starting Server on Port 8080")
	fmt.Println(http.ListenAndServe(":8080", mux))
}

// Some helpful middlware to simply log the http requests coming in
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("-> Time %s | Method %s | Path %s\n", time.Now().Format(time.RFC822), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

/*

Authorization portion of the application

*/

// Initiate the OAuth flow with github
func startOAuthWithGitHub(w http.ResponseWriter, r *http.Request) {

	// Generate a random state for the user request
	state := uuid.New().String()

	// Create Query Params for the request
	formData := url.Values{}
	formData.Set("client_id", clientID)
	formData.Set("redirect_uri", redirectURL)
	formData.Set("scope", scope)
	// Note: we tell Github what the state is
	formData.Set("state", state)
	formData.Set("allow_signup", "false")

	// Set the cookie on the client side to store the session state for 5 minutes, the flow shouldn't take longer than that
	stateCookie := &http.Cookie{
		Name:     "session_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		Secure:   secureCookie == "true",
		HttpOnly: secureCookie == "true",
	}
	fmt.Println(stateCookie.String())
	http.SetCookie(w, stateCookie)

	// Construct our request to github
	// Here we use the formData object we created to encode the data in the proper query param format and append it to the URL
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?%s", formData.Encode())

	// Finally we send the user to github to authorize with them
	http.Redirect(w, r, redirectURL, 302)
}

// Handle the callback hop after the user has been authenticated with github
func handleCallbackFromGitHub(w http.ResponseWriter, r *http.Request) {

	// We first pull out the state from the cookie, if not present, something is wrong and we throw a 401 status code
	sess, err := r.Cookie("session_state")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		w.Write([]byte("invalid state"))
		return
	}

	// Now that we have the state from the cookie, we pull the state coming from Github and match them
	// This ensures that the request is coming from us, since it's the state I gave to the user and to Github
	if r.URL.Query().Get("state") != sess.Value {
		w.WriteHeader(401)
		w.Write([]byte("invalid state"))
		return
	}

	// Now we create the final request to exchange the callback code for a user access token
	formData := url.Values{}
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("code", r.URL.Query().Get("code"))
	formData.Set("state", sess.Value)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	val := make(map[string]string)

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&val)
	if err != nil {
		fmt.Println(err)
	}

	accessToken, ok := val["access_token"]
	if !ok {
		fmt.Printf("Map value does not contain access_token | %+v\n", val)
		w.WriteHeader(401)
		w.Write([]byte("invalid token"))
		return
	}

	b := branca.NewBranca(superSecretPassword)
	encryptedToken, err := b.EncodeToString(accessToken)
	if err != nil {
		fmt.Println(err)
	}

	accessCookie := &http.Cookie{
		Name:     "session_access",
		Value:    encryptedToken,
		Path:     "/",
		Secure:   secureCookie == "true",
		HttpOnly: secureCookie == "true",
	}
	http.SetCookie(w, accessCookie)

	http.Redirect(w, r, "/", 302)
}

/*

Main logic of the application api

*/

// Get the user info, retrieve the avatar url and create repeatable unique index per avatar to match to Simpson avatar
func handleUserInfo(w http.ResponseWriter, r *http.Request) {
	sess, err := r.Cookie("session_access")
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("session_timeout"))
		return
	}

	b := branca.NewBranca(superSecretPassword)

	token, err := b.DecodeToString(sess.Value)
	if err != nil {
		fmt.Println("DecryptToken", err)
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GetGithubAPI Error", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var userInfo UserInfo

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		fmt.Println(err)
	}

	are, err := http.Get(userInfo.AvatarURL)
	if err != nil {
		fmt.Println("GetURL Error", err)
	}

	imageBytes, err := ioutil.ReadAll(are.Body)
	if err != nil {
		fmt.Println("ReadImage Error", err)
	}

	sha := sha1.New()
	sha.Write(imageBytes)
	hashSumByte := sha.Sum(nil)
	hexVal := hex.EncodeToString(hashSumByte)

	simpsonsIndex := 0
	if intBase16, success := new(big.Int).SetString(hexVal, 16); success {
		source := rand.NewSource(intBase16.Int64())
		random := rand.New(source)
		simpsonsIndex = random.Intn(len(characterIndex))
	}

	userInfo.Simpson = simpsonsCharacters[characterIndex[simpsonsIndex]]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

/*

Presentaion portion, frontend display of the application

*/

// Present the index.html file
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/html/index.html")
}

// Present the login.html file
func serveLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/html/login.html")
}
