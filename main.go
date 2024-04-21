package main

import (
	_ "embed"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/captain686/go-warp2wireguard/services"
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"os"
	"strings"
	"sync"
	"time"
)

//go:embed template/clash_template
var clashTemplate string

var (
	mu         sync.Mutex
	serverType string
)

func init() {
	flag.StringVar(&serverType, "t", "wireguard", "operating mode [ wireguard | clash ]")
	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage:\n")
		if err != nil {
			return
		}
		flag.PrintDefaults()
	}
	if _, err := os.Stat(util.ConfigPath); os.IsNotExist(err) {
		// If the folder does not exist, create the folder
		err := os.MkdirAll(util.ConfigPath, 0755)
		if err != nil {
			return
		}
	}
}

// Account registration using public key

func task(id, accountToken string) error {
	mu.Lock()
	defer mu.Unlock()
	token, err := services.FloWorker(id, accountToken)
	if err != nil {
		//log.Error(err)
		return err
	}
	log.Info("Start Account Info Query")
	util.SleepBar(60)
	err = services.Query(id, *token, util.AccountFilePath)
	if err != nil {
		return err
	}
	return nil
}

// Earn traffic using account ID
func regAndFlow() error {
	if !util.ConfigFileCheck(util.KeyFilePath) || !util.ConfigFileCheck(util.AccountFilePath) {
		log.Info("User information not detected")
		err := services.AccountReg()
		if err != nil {
			return err
		}
		//util.SleepBar(60)
		util.SleepBar(10)
	}
	if !util.ConfigFileCheck(util.GeoLite2CityPath) {
		err := util.DownloadGeoLite()
		if err != nil {
			return err
		}
	}
	accountInfo, err := util.ReadAccountInfo(util.AccountFilePath)
	if err != nil || accountInfo == nil {
		return err
	}
	account := *accountInfo
	id := account.ID
	accountToken := account.Token

	err = task(id, accountToken)
	if err != nil {
		util.RetryBar(10)
		_ = regAndFlow()
	}
	go func() {
		for {
			//util.SleepBar(60 * 30)
			time.Sleep(time.Minute * 30)
			err := task(id, accountToken)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func clashGenerate(publicKey, node string, clash util.Clash) *util.Clash {
	var (
		server string
		port   string
	)

	res := strings.Split(node, ":")
	if len(res) == 2 {
		server = res[0]
		port = res[1]
	}
	location, err := util.Ip2Location(server)
	if err != nil {
		return nil
	}
	emoji, err := util.RandomEmoji()
	if err != nil {
		return nil
	}
	proxyName := fmt.Sprintf("%s %s %s:%s", emoji, *location, server, port)
	clash.Proxies = *util.MergeNode(&clash.Proxies, proxyName, server, port, publicKey)
	clash.AutoSelect.Proxies = append(clash.AutoSelect.Proxies, proxyName)
	for i := range clash.ProxyGroups {
		clash.ProxyGroups[i].Proxies = append(clash.ProxyGroups[i].Proxies, proxyName)
	}
	return &clash
}

func main() {
	flag.Parse()
	if serverType != "wireguard" && serverType != "clash" {
		flag.Usage()
		return
	}
	err := regAndFlow()
	if err != nil {
		return
	}
	coonMap, err := services.WarpSpeedTest()
	if err != nil {
		return
	}
	clash := util.InitClash()
	mu.Lock()
	account, err := util.ReadAccountInfo(util.AccountFilePath)
	if err != nil {
		log.Error(err)
		return
	}
	mu.Unlock()
	publicKey := account.Config.Peers[0].PublicKey
	speedSort := services.SpeedSort(coonMap)
	count := 1
	if serverType == "clash" {
		for _, v := range speedSort {
			if count > 60 {
				break
			}
			node := v.Server
			clash = *clashGenerate(publicKey, node, clash)
			count++
		}
		services.Subscribe(clash, clashTemplate)
	}
	if serverType == "wireguard" {
		if _, err := os.Stat(util.OutPath); os.IsNotExist(err) {
			// If the folder does not exist, create the folder
			err := os.MkdirAll(util.OutPath, 0755)
			if err != nil {
				return
			}
		}
		speedResultFile, err := os.Create(util.OutFilePath)
		if err != nil {
			log.Error(err)
			return
		}
		defer func(speedResultFile *os.File) {
			err := speedResultFile.Close()
			if err != nil {
				return
			}
		}(speedResultFile)
		writer := csv.NewWriter(speedResultFile)
		defer writer.Flush()
		var data [][]string
		for _, v := range speedSort {
			tmp := []string{v.Server, fmt.Sprintf("%d ms", v.TimeOut)}
			data = append(data, tmp)
		}
		for _, row := range data {
			err := writer.Write(row)
			if err != nil {
				log.Error(err)
				return
			}
		}
		err = services.GenerateConf(speedSort[0].Server, publicKey)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
