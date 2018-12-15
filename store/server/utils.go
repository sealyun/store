package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	clientID     = "89c1b05d77fb1c92a1ef"
	clientSecret = "541ddd76e65abeabd12ad9f8b02f6601394d3ad0"
)

//User is
type User struct {
	Login       string `json:"login"`
	ID          int    `json:"id"`
	AvatarURL   string `json:"avatar_url"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	SiteAdmin   bool   `json:"site_admin"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"Location"`
	Email       string `json:"email"`
	Hireable    bool   `json:"hireable"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	PublicGists int    `json:"public_gists"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

//GetLoginURL is
func GetLoginURL(state string) string {
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s", clientID, state)
	return url
}

//GetGithubAccessToken is
func GetGithubAccessToken(id, secret, code string) (token string, err error) {
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		id, secret, code)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	token = string(body)

	s1 := strings.Split(token, "=")
	if len(s1) == 0 {
		return "", errors.New("can't fetch token")
	}
	s2 := strings.Split(s1[1], "&")
	if len(s2) == 0 {
		return "", errors.New("can't fetch token")
	}
	token = s2[0]

	fmt.Printf("sucess fetch token : %s", token)
	return
}

//GetUserInfo is
//https://api.github.com/user?access_token=access_token
func GetUserInfo(token string) (user *User, err error) {
	url := fmt.Sprintf("https://api.github.com/user?access_token=%s", token)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	u := &User{}

	err = json.Unmarshal(body, u)
	if err != nil {
		return nil, err
	}

	user = u

	fmt.Printf("user info: %v", user)

	return
}
