package util

import (
	"crypto/rand"
	"github.com/charmbracelet/log"
)

type Proxies struct {
	IP               string
	Name             string
	Port             string
	PrivateKey       string
	PublicKey        string
	RemoteDNSResolve bool
	Server           string
	Type             string
	UDP              bool
}

type ProxyGroups struct {
	Name     string
	Proxies  []string
	Type     string
	Interval int
}

type AutoSelect struct {
	ProxyGroups
	URL string
}

type Clash struct {
	Proxies     []Proxies
	AutoSelect  AutoSelect
	ProxyGroups []ProxyGroups
}

func InitClash() Clash {
	clash := Clash{
		Proxies:     []Proxies{},
		AutoSelect:  AutoSelect{},
		ProxyGroups: []ProxyGroups{},
	}
	proxyGroups := []string{"🚀 节点选择", "📲 电报信息", "🎯 全球直连", "🐟 漏网之鱼", "🍎 苹果服务", "Ⓜ️ 微软服务", "🛑 全球拦截", "🍃 应用净化", "🌍 国外媒体", "📢 谷歌FCM"}
	clash.AutoSelect.Name = "♻️ 自动选择"
	clash.AutoSelect.Type = "url-test"
	clash.AutoSelect.URL = "http://www.gstatic.com/generate_204"
	clash.AutoSelect.Interval = 300
	for _, groupName := range proxyGroups {
		newGroup := InitNewProxyGroup()
		newGroup.Name = groupName
		switch groupName {
		case "🎯 全球直连":
			newGroup.Proxies = append(newGroup.Proxies, "DIRECT")
		case "🛑 全球拦截", "🍃 应用净化":
			newGroup.Proxies = append(newGroup.Proxies, "DIRECT")
			newGroup.Proxies = append(newGroup.Proxies, "REJECT")
		case "🚀 节点选择":
			newGroup.Proxies = append(newGroup.Proxies, "♻️ 自动选择")
		default:
			newGroup.Proxies = append(newGroup.Proxies, "🚀 节点选择")
			newGroup.Proxies = append(newGroup.Proxies, "♻️ 自动选择")
			newGroup.Proxies = append(newGroup.Proxies, "🎯 全球直连")
		}
		clash.ProxyGroups = append(clash.ProxyGroups, newGroup)
	}
	return clash
}

func InitProxy() Proxies {
	return Proxies{
		RemoteDNSResolve: false,
		UDP:              true,
		Type:             "wireguard",
		IP:               "172.16.0.2",
	}
}

func InitNewProxyGroup() ProxyGroups {
	proxyGroup := ProxyGroups{
		Type:     "select",
		Interval: 300,
	}
	return proxyGroup
}

func RandomEmoji() (string, error) {
	emoji := [][]int{
		// Emoticons icons
		{128513, 128591},
		// Transport and map symbols
		{128640, 128704},
	}

	// Generate a random number between 0 and total number of emojis
	maxRange := 0
	for _, r := range emoji {
		maxRange += r[1] - r[0] + 1
	}

	var randomBytes [1]byte
	_, err := rand.Read(randomBytes[:])
	if err != nil {
		log.Error(err)
		return "", err
	}
	randomValue := int(randomBytes[0])

	selectedEmoji := ""
	count := 0
	for _, r := range emoji {
		minVul := r[0]
		maxVul := r[1]
		for i := minVul; i <= maxVul; i++ {
			if count == randomValue%maxRange {
				selectedEmoji = string(rune(i))
				break
			}
			count++
		}
		if selectedEmoji != "" {
			break
		}
	}

	return selectedEmoji, nil
}

func MergeNode(proxies *[]Proxies, name, server, port, publicKey string) *[]Proxies {
	key, err := ReadKey(KeyFilePath)
	if err != nil {
		log.Error(err)
		return nil
	}
	privateKey := key.PrivateKey
	temp := InitProxy()
	temp.Server = server
	temp.Port = port
	temp.PrivateKey = privateKey
	temp.PublicKey = publicKey
	temp.Name = name
	*proxies = append(*proxies, temp)
	return proxies
}
