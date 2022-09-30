package tokens

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/classmethod/aurl/utils"
)

const TOKEN_STORAGE_DIR = "~/.aurl/token"
const TOKEN_STORAGE_FORMAT = TOKEN_STORAGE_DIR + "/%s.json"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
	Expires      int32  `json:"expires"` // broken Facebook spelling of expires_in
	Timestamp    int64  `json:"timestamp"`
}

func (tokenResponse TokenResponse) IsExpired() bool {
	return false // TODO
}

func New(tokenResponseString *string) (TokenResponse, error) {
	if tokenResponseString == nil {
		return TokenResponse{}, errors.New("nil response")
	}
	var tokenResponse TokenResponse
	jsonParser := json.NewDecoder(strings.NewReader(*tokenResponseString))
	if err := jsonParser.Decode(&tokenResponse); err != nil {
		log.Printf("Failed to parse token response: %v", err)
		return TokenResponse{}, err
	}
	tokenResponse.Timestamp = time.Now().Unix()
	return tokenResponse, nil
}

func LoadTokenResponseString(profileName string) (*string, error) {
	path := utils.ExpandPath(fmt.Sprintf(TOKEN_STORAGE_FORMAT, profileName))
	buf, err := ioutil.ReadFile(path)
	if err == nil {
		s := string(buf)
		return &s, nil
	}
	return nil, err
}

func SaveTokenResponseString(profileName string, tokenResponseString *string) {
	os.Mkdir(utils.ExpandPath(TOKEN_STORAGE_DIR), os.FileMode(0755))
	path := utils.ExpandPath(fmt.Sprintf(TOKEN_STORAGE_FORMAT, profileName))
	content := []byte(*tokenResponseString)
	err := ioutil.WriteFile(path, content, os.FileMode(0600))
	if err != nil {
		log.Printf("Failed to save token response: %v", err.Error())
	}
}
