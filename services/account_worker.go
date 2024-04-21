package services

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ApiUrl     = "https://api.cloudflareclient.com"
	ApiVersion = "v0a1922"
)

var headers = HeaderMap("{\"Host\": \"api.cloudflareclient.com\",\n\"User-Agent\": \"okhttp/3.12.1\",\n\"Accept\": \"application/json\",\n\"Cf-Client-Version\": \"a-6.3-1922\",\n\"Content-Type\": \"application/json\",\n\"Accept-Encoding\": \"gzip, deflate, br\"}")

func randomElement(elements []string) string {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(elements))))
	if err != nil {
		log.Fatal(err)
	}
	indexInt := index.Int64()
	return elements[indexInt]
}

func generateRandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Error(err)
		return base64.URLEncoding.EncodeToString([]byte("zdsgniuoabsrezdfhdtrh"))[:length]
	}
	str := base64.URLEncoding.EncodeToString(randomBytes)[:length]
	return str
}

func Query(id, token, accountFilePath string) error {
	if !util.FileCheck(util.CookieFilePath) {
		err := fmt.Errorf("cookie is empty")
		log.Error(err)
		return err
	}
	cookie, err := os.ReadFile(util.CookieFilePath)
	if err != nil {
		log.Error("Read Cookie File Error")
		return err
	}
	queryUrl := fmt.Sprintf("%s/%s/reg/%s", ApiUrl, ApiVersion, id)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	headers["Cookie"] = string(cookie)
	request := Req{
		Method: "GET",
		Header: headers,
		Host:   queryUrl,
		Gzip:   true,
	}
	response, err := request.Requests()
	if err != nil {
		log.Error("Error getting account information", err)
		return err
	}
	if response.StatusCode != 200 {
		log.Error("Account Query Requests Status Code", response.StatusCode)
		return fmt.Errorf("requests Error")
	}
	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Error("gzip NewReader error:", err)
		return err
	}
	defer func(reader *gzip.Reader) {
		err := reader.Close()
		if err != nil {
			return
		}
	}(reader)
	accountInfo := util.Response{}
	err = yaml.NewDecoder(reader).Decode(&accountInfo)
	if err != nil {
		log.Error("Yaml format Error ", err)
		return err
	}
	accountInfo.Token = token
	log.Info(fmt.Sprintf("premium data: %d", accountInfo.Account.PremiumData))
	log.Info(fmt.Sprintf("quota   data: %d", accountInfo.Account.Quota))
	err = util.Save2Yaml(accountInfo, accountFilePath)
	if err != nil {
		log.Error("Account Save Error: ", err)
		return err
	}
	log.Info("Account information query successful")
	return nil
}

func postRequest(publicKey, referrer string) (*http.Response, error) {
	installID := generateRandomString(43)
	data := util.PostData{
		FcmToken:   fmt.Sprintf("%s:APA91b%s", installID, generateRandomString(134)),
		InstallId:  installID,
		WarpEnable: true,
		Key:        publicKey,
		Local:      "en_US",
		Model:      "pc",
		Tos:        util.GetTimestamp(),
		Type:       randomElement([]string{"Android", "iOS"}),
	}
	if referrer != "" {
		data.Referrer = referrer
	}
	jsonPostData, err := json.Marshal(data)
	if err != nil {
		log.Error("JSON marshal error:", err)
		return nil, err
	}

	requests := Req{
		Method:  "POST",
		Header:  headers,
		Gzip:    true,
		Host:    fmt.Sprintf("%s/%s/reg", ApiUrl, ApiVersion),
		Data:    jsonPostData,
		Timeout: 30 * time.Second,
	}

	response, err := requests.Requests()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Error("response.StatusCode", response.StatusCode)
		return nil, err
	}

	return response, nil
}

func FloWorker(accountId, accountToken string) (*string, error) {
	privateKey, err := NewPrivateKey()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("VPN traffic recharging")
	//accountKey := *key
	response, err := postRequest(privateKey.Public().String(), accountId)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(response.Body)
	cookies := response.Header["Set-Cookie"]
	if len(cookies) != 0 {
		setCookie := strings.Split(cookies[0], ";")
		if len(setCookie) != 0 {
			err := os.WriteFile(util.CookieFilePath, []byte(setCookie[0]), 0644)
			if err != nil {
				log.Error("cookie file write error", err)
				return nil, err
			}
		}
	}
	log.Info(fmt.Sprintf("VPN traffic recharging Response Status_Code %d", response.StatusCode))
	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Error("gzip NewReader error:", err)
		return nil, err
	}
	defer func(reader *gzip.Reader) {
		err := reader.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(reader)
	log.Info("VPN traffic recharge successful")
	accountData := util.Response{}
	err = yaml.NewDecoder(reader).Decode(&accountData)
	if err != nil {
		log.Error("Yaml format Error ", err)
		return nil, err
	}
	return &accountToken, nil
}

func AccountReg() error {
	log.Info("Start account registration")
	privateKey, err := NewPrivateKey()
	if err != nil {
		log.Error(err)
		return err
	}
	keyData := util.Key{
		PrivateKey: privateKey.String(),
		PublicKey:  privateKey.Public().String(),
	}
	err = util.Save2Yaml(keyData, util.KeyFilePath)
	if err != nil {
		log.Error(err)
		return err
	}
	publicKey := privateKey.Public()
	response, err := postRequest(publicKey.String(), "")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(response.Body)

	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Error("gzip NewReader error:", err)
		return err
	}
	defer func(reader *gzip.Reader) {
		err := reader.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(reader)
	accountData := util.Response{}
	err = yaml.NewDecoder(reader).Decode(&accountData)
	if err != nil {
		log.Error("Yaml format Error ", err)
		return err
	}
	err = util.Save2Yaml(accountData, util.AccountFilePath)
	if err != nil {
		log.Error("Account Save Error: ", err)
		return err
	}
	fmt.Println("--------------------------------warning--------------------------------")
	fmt.Println("|                    Account registration successful                  |")
	fmt.Println("|            The account information is in the config directory       |")
	fmt.Println(fmt.Sprintf("|            account license :  %s            |", accountData.Account.License))
	fmt.Println("-----------------------------------------------------------------------")
	return nil
}
