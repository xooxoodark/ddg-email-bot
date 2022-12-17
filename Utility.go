package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var baseurl = "https://quack.duckduckgo.com/"

var client = resty.New()

func RequestOTP(UserName string) {
	var params = map[string]string{}
	params["user"] = UserName
	rq, err := client.R().SetQueryParams(params).Get(baseurl + "api/auth/loginlink")
	if err != nil {
		fmt.Println(err)
	}
	if rq.StatusCode() != 200 {
		fmt.Println(rq.String())
	}
}

type OPTResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	User   string `json:"user"`
}

type DDGUser struct {
	AccessToken string `json:"access_token"`
}
type TokenResponse struct {
	DDUser DDGUser `json:"user"`
}

func GetToken(OTP, UserName string) (string, error) {
	rq, err := client.R().
		SetQueryParams(map[string]string{"otp": OTP, "user": UserName}).
		SetResult(&OPTResponse{}).Get(baseurl + "api/auth/login")
	if err != nil {
		return "", err
	}
	opt := rq.Result().(*OPTResponse)

	rq, err = client.R().SetAuthToken(opt.Token).SetResult(&TokenResponse{}).Get(baseurl + "api/email/dashboard")
	if err != nil {
		return "", err
	}
	tokenresponse := rq.Result().(*TokenResponse)

	return tokenresponse.DDUser.AccessToken, nil
}

type Email struct {
	Address string `json:"address"`
}

func Generate(token Token) (string, error) {
	rq, err := client.R().
		SetAuthToken(token.Token).
		SetResult(&Email{}).Post(baseurl + "api/email/addresses")

	if err != nil {
		return "", err
	}
	email := rq.Result().(*Email)

	return email.Address + "@duck.com", nil

}
