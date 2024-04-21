package util

import (
	"fmt"
	"time"
)

type Response struct {
	ID      string `yaml:"id"`
	Type    string `yaml:"type"`
	Model   string `yaml:"model"`
	Name    string `yaml:"name"`
	Key     string `yaml:"key"`
	Account struct {
		ID                       string    `yaml:"id"`
		AccountType              string    `yaml:"account_type"`
		Created                  time.Time `yaml:"created"`
		Updated                  time.Time `yaml:"updated"`
		PremiumData              int       `yaml:"premium_data"`
		Quota                    int       `yaml:"quota"`
		Usage                    int       `yaml:"usage"`
		WarpPlus                 bool      `yaml:"warp_plus"`
		ReferralCount            int       `yaml:"referral_count"`
		ReferralRenewalCountdown int       `yaml:"referral_renewal_countdown"`
		Role                     string    `yaml:"role"`
		License                  string    `yaml:"license"`
	} `yaml:"account"`
	Config struct {
		ClientID string `yaml:"client_id"`
		Peers    []struct {
			PublicKey string `yaml:"public_key"`
			Endpoint  struct {
				V4    string `yaml:"v4"`
				V6    string `yaml:"v6"`
				Host  string `yaml:"host"`
				Ports []int  `yaml:"ports"`
			} `yaml:"endpoint"`
		} `yaml:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `yaml:"v4"`
				V6 string `yaml:"v6"`
			} `yaml:"addresses"`
		} `yaml:"interface"`
		Services struct {
			HTTPProxy string `yaml:"http_proxy"`
		} `yaml:"services"`
	} `yaml:"config"`
	Token           string    `yaml:"token"`
	WarpEnabled     bool      `yaml:"warp_enabled"`
	WaitlistEnabled bool      `yaml:"waitlist_enabled"`
	Created         time.Time `yaml:"created"`
	Updated         time.Time `yaml:"updated"`
	Tos             time.Time `yaml:"tos"`
	Place           int       `yaml:"place"`
	Locale          string    `yaml:"locale"`
	Enabled         bool      `yaml:"enabled"`
	InstallID       string    `yaml:"install_id"`
	FcmToken        string    `yaml:"fcm_token"`
	Referrer        string    `yaml:"referrer"`
}

type PostData struct {
	FcmToken   string `json:"fcm_token"`
	InstallId  string `json:"install_id"`
	Key        string `json:"key"`
	Local      string `json:"local"`
	Model      string `json:"model"`
	Tos        string `json:"tos"`
	Type       string `json:"type"`
	WarpEnable bool   `json:"warp_enable"`
	Referrer   string `json:"referrer"`
}

type Key struct {
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
}

var (
	ConfigPath        = "config"
	OutPath           = "out"
	KeyFilePath       = fmt.Sprintf("%s/key.yml", ConfigPath)
	AccountFilePath   = fmt.Sprintf("%s/account.yml", ConfigPath)
	CookieFilePath    = fmt.Sprintf("%s/cookie", ConfigPath)
	GeoLite2CityPath  = fmt.Sprintf("%s/GeoLite2-City.mmdb", ConfigPath)
	WireGuardConfPath = fmt.Sprintf("%s/wireguard.conf", ConfigPath)
	OutFilePath       = fmt.Sprintf("%s/data.csv", OutPath)
)

func GetTimestamp() string {
	return time.Now().Format(time.RFC3339Nano)
}
