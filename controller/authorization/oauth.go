package authorization

import (
	"github.com/bytedance/sonic"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	state       = "Gauss curvature"
	OauthConfig *oauth2.Config
)

/*
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		Scope string `json:"scope"`
		TokenType string `json:"token_type"`
	}
*/
type UserResponse struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func InitialConfig() {
	OauthConfig = &oauth2.Config{
		ClientID:     "Iv1.d7d4884211aa1791",
		ClientSecret: "00b199cff9f402f4daa0b97ce698719044b71951",
		RedirectURL:  "http://127.0.0.1:9999/av/callback",
		Scopes:       []string{"read:user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

// If you use PKCE, need to enable them
/*func generateCodeVerifier() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	codeVerifier := base64.URLEncoding.EncodeToString(randomBytes)
	return codeVerifier, nil
}

func generateCodeChallenge(codeVerifier string) string {
	codeChallenge := sha256.Sum256([]byte(codeVerifier))
	return base64.URLEncoding.EncodeToString(codeChallenge[:])
}

func generateCode() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}*/

//	 GitHub Oauth Login
//
//		@Tags			GitHub Oauth
//		@Success		200	{string} {"UserName":userResp.Name,"AvatarURL":userResp.AvatarURL,}
//		@router			/authorization/login [get]
func GitLogin(context *gin.Context) {
	// init oauth config
	InitialConfig()

	// If you use PKCE:
	/*var PKCE bool = false
	if PKCE == true {
		codeVerifier, err := generateCodeVerifier()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code verifier"})
			return
		}
		code, err := generateCode()
		codeChallenge := generateCodeChallenge(codeVerifier)
		context.SetCookie("code_verifier", codeVerifier, 3600, "", "", false, true)
	}*/

	url := OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	context.Redirect(http.StatusTemporaryRedirect, url)
}

func GitCallBack(context *gin.Context) {
	code := context.Query("code")
	/*
		codeVerifier, err := context.Cookie("code_verifier")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve code verifier"})
			return
		}

			token, err := oauthConfig.Exchange(context, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	*/
	/* Prepare the data for the POST request
	data := url.Values{}
	data.Set("client_id", OauthConfig.ClientID)
	data.Set("client_secret", OauthConfig.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", OauthConfig.RedirectURL)
	data.Set("state", state)
	request, err := http.NewRequest("POST", OauthConfig.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()
	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	fmt.Println("This is body")
	fmt.Println(string(body))
	// The body contains the access token
	context.JSON(http.StatusOK, gin.H{"access_token": string(body)})
	var tokenResp tokenResponse
	Parser := sonic.Unmarshal(body, &tokenResp)
	if Parser != nil {
		log.Printf("Failed to parse token response: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}
	access_token := tokenResp.AccessToken
	GetUserDetails(context, access_token)
	*/
	token, err := OauthConfig.Exchange(context, code)
	if err != nil {
		log.Println("Failed to exchange token", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	GetUserDetails(context, token.AccessToken)
}

func GetUserDetails(context *gin.Context, token string) {
	// Create a new request
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add the access token to the Authorization header
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln("Failed to close body")
			return
		}
	}(resp.Body)

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Unmarshal the JSON response
	var userResp UserResponse
	err = sonic.Unmarshal(body, &userResp)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user response"})
		return
	}

	// The userResp contains the user's name, email, and avatar URL
	context.HTML(http.StatusOK, "callback.html", gin.H{
		"UserName":  userResp.Name,
		"AvatarURL": userResp.AvatarURL,
	})
}
