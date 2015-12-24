// weixin_sdk
package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	accessToken string
	timeStamp   int64
	signRes     signSt
	signResStr  []byte
)

type tokenResponse struct {
	Access_token string
	Expires_in   int
}

type ticketResponse struct {
	Errcode    int
	Errmsg     string
	Ticket     string
	Expires_in int
}

type signSt struct {
	TimeStamp int64
	Noncestr  string
	Sign      string
	AppID     string
}

type oauthToken struct {
	Access_token  string
	Expires_in    int
	Refresh_token string
	Openid        string
	Scope         string
}

type userInfo struct {
	Openid     string
	Nickname   string
	Sex        int
	Language   string
	City       string
	Province   string
	Headimgurl string
	Privilege  []interface{}
}

const (
	SIGN_MODULE string = "jsapi_ticket=$1&noncestr=$2&timestamp=$3&url="
	TimeExpire  int64  = 7100000
)

func GetSign(req *http.Request) string {
	req.ParseForm()
	if time.Now().Unix()-timeStamp < TimeExpire {
		return string(signResStr)
	}
	ticket := reqTicket()
	if ticket == "" {
		return "get ticket error"
	}
	timeStamp = time.Now().Unix()
	signRes.TimeStamp = timeStamp
	signRes.Noncestr = getNonceStr()
	signRes.AppID = AppID
	signRes.Sign = strings.Replace(SIGN_MODULE, "$1", ticket, 1)
	signRes.Sign = strings.Replace(signRes.Sign, "$2", signRes.Noncestr, 1)
	signRes.Sign = strings.Replace(signRes.Sign, "$3", strconv.FormatInt(signRes.TimeStamp, 10), 1)
	signRes.Sign = signRes.Sign + GameURL
	sha1Init := sha1.New()
	fmt.Println("sign:", signRes.Sign)
	io.WriteString(sha1Init, signRes.Sign)
	signRes.Sign = fmt.Sprintf("%x", sha1Init.Sum(nil))
	signResStr, _ = json.Marshal(signRes)
	fmt.Println(signResStr)
	return string(signResStr)
}

func getNonceStr() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%x", r.Intn(10000000))
}

func reqToken() string {
	var resp *http.Response
	var err error
	var result []byte
	var tokenResponse tokenResponse
	resp, err = http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppID + "&secret=" + AppKey)
	if err != nil {
		return ""
	}

	result, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return ""
	}

	err = json.Unmarshal(result, &tokenResponse)
	if err != nil {
		return ""
	}
	return tokenResponse.Access_token
}

func reqTicket() string {
	var resp *http.Response
	var err error
	var result []byte
	var ticketResponse ticketResponse
	accessToken = reqToken()
	if accessToken == "" {
		return "get tocken error"
	}

	resp, err = http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi")
	if err != nil {
		return ""
	}
	result, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return ""
	}

	err = json.Unmarshal(result, &ticketResponse)
	if err != nil {
		return ""
	}
	return ticketResponse.Ticket
}

func GetUserInfo(req *http.Request) string {
	req.ParseForm()
	GameURL = strings.Split(req.Referer(), "?")[0]
	code := req.Form.Get("code")
	if code == "" {
		return "please input code param"
	}

	token, openID := reqOAuthToken(code)
	if token != "" && openID != "" {
		return reqUserInfo(token, openID)
	}
	return "invalid code"
}

func reqUserInfo(token string, openID string) string {
	var result []byte
	var info userInfo
	resp, err := http.Get("https://api.weixin.qq.com/sns/userinfo?access_token=" + token + "&openid=" + openID)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	result, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	err = json.Unmarshal(result, &info)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Printf(string(result))
	return string(result)
}

func reqOAuthToken(code string) (string, string) {
	var result []byte
	var token oauthToken
	resp, err := http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + AppID + "&secret=" + AppKey + "&code=" + code + "&grant_type=authorization_code")
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	result, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	err = json.Unmarshal(result, &token)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	return token.Access_token, token.Openid
}
