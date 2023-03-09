package users

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type GoogleUser struct {
// 	Email string `json:"email"`
// 	Name  string `json:"name"`
// }

// func GetGoogleUser(token string, user *GoogleUser) error {
// 	client := &http.Client{}

// 	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)

// 	if err != nil {
// 		return err
// 	}

// 	// Query params
// 	q := req.URL.Query()
// 	q.Add("access_token", token)
// 	req.URL.RawQuery = q.Encode()

// 	// Authorization
// 	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

// 	// Execute query
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	// Decode body
// 	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
// 		return err
// 	}

// 	return nil
// }
