package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	redditKey := "o9LvVjAhrThNDJqGgHwphQ"
	redditSecret := "vRgiA9Wy9N3qjsbfUky6PODnhZ9qLQ"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     "http://localhost:3001",
			WebsiteDomain: "http://localhost:3000",
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				/*
				   We use different credentials for different platforms when required. For example the redirect URI for Github
				   is different for Web and mobile. In such a case we can provide multiple providers with different client Ids.

				   When the frontend makes a request and wants to use a specific clientId, it needs to send the clientId to use in the
				   request. In the absence of a clientId in the request the SDK uses the default provider, indicated by `isDefault: true`.
				   When adding multiple providers for the same type (Google, Github etc), make sure to set `isDefault: true`.
				*/
				Providers: []tpmodels.TypeProvider{
					// We have provided you with development keys which you can use for testsing.
					// IMPORTANT: Please replace them with your own OAuth keys for production use.
					thirdparty.Google(tpmodels.GoogleConfig{
						// We use this for websites
						IsDefault:    true,
						ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
						ClientSecret: "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",
					}),
					thirdparty.Google(tpmodels.GoogleConfig{
						// we use this for mobile apps
						ClientID:     "1060725074195-c7mgk8p0h27c4428prfuo3lg7ould5o7.apps.googleusercontent.com",
						ClientSecret: "", // this is empty because we follow Authorization code grant flow via PKCE for mobile apps (Google doesn't issue a client secret for mobile apps).
					}),
					thirdparty.Github(tpmodels.GithubConfig{
						// We use this for websites
						IsDefault:    true,
						ClientID:     "467101b197249757c71f",
						ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
					}),
					thirdparty.Github(tpmodels.GithubConfig{
						// We use this for mobile apps
						ClientID:     "8a9152860ce869b64c44",
						ClientSecret: "00e841f10f288363cd3786b1b1f538f05cfdbda2",
					}),
					/*
					   For Apple signin, iOS apps always use the bundle identifier as the client ID when communicating with Apple. Android, Web and other platforms
					   need to configure a Service ID on the Apple developer dashboard and use that as client ID.
					   In the example below 4398792-io.supertokens.example.service is the client ID for Web. Android etc and thus we mark it as default. For iOS
					   the frontend for the demo app sends the clientId in the request which is then used by the SDK.
					*/
					thirdparty.Apple(tpmodels.AppleConfig{
						// For Android and website apps
						IsDefault: true,
						ClientID:  "4398792-io.supertokens.example.service",
						ClientSecret: tpmodels.AppleClientSecret{
							KeyId:      "7M48Y4RYDL",
							PrivateKey: "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
							TeamId:     "YWQCXGJRJL",
						},
					}),
					thirdparty.Apple(tpmodels.AppleConfig{
						// For iOS Apps
						ClientID: "4398792-io.supertokens.example",
						ClientSecret: tpmodels.AppleClientSecret{
							KeyId:      "7M48Y4RYDL",
							PrivateKey: "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
							TeamId:     "YWQCXGJRJL",
						},
					}),
					{
						ID: "reddit",
						Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
							if redirectURI == nil {
								temp := ""
								redirectURI = &temp
							}

							if authCodeFromRequest == nil {
								temp := ""
								authCodeFromRequest = &temp
							}

							return tpmodels.TypeProviderGetResponse{
								AccessTokenAPI: tpmodels.AccessTokenAPI{
									// this contains info about the token endpoint which exchanges the auth code with the access token and profile info.
									URL: "https://www.reddit.com/api/v1/access_token",
									Params: map[string]string{
										// example post params
										"client_id":     redditKey,
										"client_secret": redditSecret,
										"grant_type":    "authorization_code",
										"redirect_uri":  *redirectURI,
										"code":          *authCodeFromRequest,
										//...
									},
								},
								AuthorisationRedirect: tpmodels.AuthorisationRedirect{
									// this contains info about forming the authorisation redirect URL without the state params and without the redirect_uri param
									URL: "https://www.reddit.com/api/v1/authorize",
									Params: map[string]interface{}{
										"client_id":     redditKey,
										"scope":         "identity",
										"response_type": "code",
										"duration":      "permanent",
									},
								},
								GetClientId: func(userContext supertokens.UserContext) string {
									return redditKey
								},
								GetProfileInfo: func(authCodeResponse interface{}, userContext supertokens.UserContext) (tpmodels.UserInfo, error) {
									fmt.Println(authCodeResponse)
									accessToken := authCodeResponse.(map[string]interface{})["access_token"].(string)
									authHeader := "Bearer " + accessToken
									response, err := GetAuthRequest(authHeader)
									if err != nil {
										return tpmodels.UserInfo{}, err
									}
									userInfo := response.(map[string]interface{})
									log.Println(userInfo)
									_, emailOk := userInfo["email"]
									if !emailOk {
										return tpmodels.UserInfo{
											ID:    userInfo["id"].(string),
											Email: nil,
										}, nil
									}
									return tpmodels.UserInfo{
										ID: userInfo["id"].(string),
										Email: &tpmodels.EmailStruct{
											ID:         userInfo["email"].(string),
											IsVerified: userInfo["verified"].(bool),
										},
									}, nil
								},
							}
						},
					},
				},
			}),
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	http.ListenAndServe(":3001", corsMiddleware(
		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Handle your APIs..

			if r.URL.Path == "/sessioninfo" {
				session.VerifySession(nil, sessioninfo).ServeHTTP(rw, r)
				return
			}
		}))))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

func sessioninfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	sessionData, err := sessionContainer.GetSessionData()
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	w.WriteHeader(200)
	w.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error in converting to json"))
	} else {
		w.Write(bytes)
	}
}

func GetAuthRequest(authHeader string) (interface{}, error) {
	url := "https://oauth.reddit.com/api/v1/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return doGetRequest(req)
}

func doGetRequest(req *http.Request) (interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
