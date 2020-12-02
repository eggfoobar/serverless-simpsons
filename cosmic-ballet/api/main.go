package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

var (
	awsSession *session.Session
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
	awsSession = session.Must(session.NewSession())
	lambda.Start(handleUserInfo)
}

// Get the user info, retrieve the avatar url and create repeatable unique index per avatar to match to Simpson avatar
func handleUserInfo(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Parse out cookie from API Gateway event
	sessionToken := parseCookie(r.Headers, "session_access")
	if sessionToken == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       "access denied",
		}, nil
	}

	// Decode our base64 token
	encryptedToken, err := base64.StdEncoding.DecodeString(sessionToken)
	if err != nil {
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       "invalid token",
		}, nil
	}

	// Create a KMS client and the decrypt payload, which is our base64 decoded token, the metadata in the encrypted secret will be used by AWS
	svc := kms.New(awsSession)
	input := &kms.DecryptInput{
		CiphertextBlob: encryptedToken,
	}

	// Decrypt the token from the cookie
	output, err := svc.Decrypt(input)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       "invalid token",
		}, nil
	}

	// Grab our decrypted github token and use it to fetch user information
	token := string(output.Plaintext)
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("GetGithubAPI Error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to call github",
		}, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "authorization denied",
		}, nil
	}

	var userInfo UserInfo

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		fmt.Println(err)
	}

	// User avatar data is public, we don't need a token here, we fetch the info
	are, err := http.Get(userInfo.AvatarURL)
	if err != nil {
		fmt.Println("GetURL Error", err)
	}

	// We read in our bytes and do the work to map the avatar to a simpsons character
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

	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to create userinfo",
		}, nil
	}

	// We then send everything back to the frontend to present the information
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(userInfoBytes),
	}, nil

}

// This cookie parsing is done as a quick and very explicit teaching purpose, do not use in production.
// Always use a library that conforms to the Cookie standard to have consistant results.
func parseCookie(headers map[string]string, name string) string {
	var cookieString string
	if cookie, ok := headers["Cookie"]; ok {
		cookieString = cookie
	}
	if cookie, ok := headers["cookie"]; ok {
		cookieString = cookie
	}

	cookieArray := strings.Split(cookieString, ";")
	for _, cookie := range cookieArray {
		// Split at the first = symbol to avoid any n number of = symbols that might be present in the string
		kv := strings.SplitN(cookie, "=", 2)
		if len(kv) > 0 && kv[0] == name {
			return kv[1]
		}
	}
	return ""
}
