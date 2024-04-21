package services

import (
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"os"
	"text/template"
)

type WireGuardConfig struct {
	PrivateKey string
	PublicKey  string
	Endpoint   string
}

var confTemplate = `[Interface]
PrivateKey = {{ .PrivateKey }}
Address = 172.16.0.2/32, 2606:4700:110:8b6a:de1:9549:27d5:8d6d/128
DNS = 1.1.1.1, 1.0.0.1, 2606:4700:4700::1111, 2606:4700:4700::1001
MTU = 1280

[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = 0.0.0.0/1, 128.0.0.0/1, ::/1, 8000::/1
Endpoint = {{ .Endpoint }}
PersistentKeepalive = 15
`

func (wg WireGuardConfig) initConf() error {
	tmpl, err := template.New("conf").Parse(confTemplate)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(util.WireGuardConfPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, wg)
	if err != nil {
		return err
	}
	return nil
}

func GenerateConf(endpoint, publicKey string) error {
	key, err := util.ReadKey(util.KeyFilePath)
	if err != nil {
		log.Error(err)
		return nil
	}
	wg := WireGuardConfig{
		PrivateKey: key.PrivateKey,
		PublicKey:  publicKey,
		Endpoint:   endpoint,
	}
	err = wg.initConf()
	if err != nil {
		return err
	}
	return nil
}
