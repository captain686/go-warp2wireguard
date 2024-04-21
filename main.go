package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/captain686/go-warp2wireguard/services"
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"os"
	"sync"
)

//go:embed template/clash_template
var clashTemplate string

var (
	mu          sync.Mutex
	serviceType string
)

func init() {
	flag.StringVar(&serviceType, "t", "wireguard", "operating mode [ wireguard | clash ]")
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
	if _, err := os.Stat(util.OutPath); os.IsNotExist(err) {
		// If the folder does not exist, create the folder
		err := os.MkdirAll(util.OutPath, 0755)
		if err != nil {
			return
		}
	}
}

// Account registration using public key
func task(id, accountToken string) {
	mu.Lock()
	defer mu.Unlock()
	token, err := services.FloWorker(id, accountToken)
	if err != nil {
		util.SleepBar("Waiting for retry", 10)
		task(id, accountToken)
	}
	log.Info("Start Account Info Query")
	util.SleepBar("Sleep Countdown", 60)
	err = services.Query(id, *token, util.AccountFilePath)
	if err != nil {
		util.SleepBar("Waiting for retry", 10)
		task(id, accountToken)
	}
}

// Earn traffic using account ID
func regAndFlow(serviceType string) error {
	if !util.FileCheck(util.KeyFilePath) || !util.FileCheck(util.AccountFilePath) {
		log.Info("User information not detected")
		err := services.AccountReg()
		if err != nil {
			return err
		}
		//util.SleepBar(60)
		util.SleepBar("Sleep Countdown", 10)
	}
	if !util.FileCheck(util.GeoLite2CityPath) {
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

	//if serviceType == "wireguard" {
	task(id, accountToken)
	//}

	if serviceType == "clash" {
		go func() {
			cronFunc := services.CronFunction{
				Function:     task,
				Id:           id,
				AccountToken: accountToken,
			}
			cronFunc.CronServe()
		}()
	}
	return nil
}

func speedTest(speedSort []util.Speed) []util.Speed {
	coonMap, err := services.WarpSpeedTest()
	if err != nil {
		log.Error(err)
		return nil
	}
	speedSort = services.SpeedSort(coonMap)
	err = util.WriteToCsv(speedSort)
	if err != nil {
		log.Error(err)
		return nil
	}
	return speedSort
}

func main() {
	flag.Parse()
	if serviceType != "wireguard" && serviceType != "clash" {
		flag.Usage()
		return
	}
	err := regAndFlow(serviceType)
	if err != nil {
		return
	}
	var speedSort []util.Speed

	mu.Lock()
	account, err := util.ReadAccountInfo(util.AccountFilePath)
	if err != nil {
		log.Error(err)
		return
	}
	mu.Unlock()
	publicKey := account.Config.Peers[0].PublicKey

	if serviceType == "clash" {
		go func() {
			newSpeedSort := speedTest(speedSort)
			err = util.WriteToCsv(newSpeedSort)
			if err != nil {
				log.Error(err)
				return
			}
		}()
		services.Subscribe(publicKey, clashTemplate)
	}
	if serviceType == "wireguard" {
		speedSort = speedTest(speedSort)
		err = util.WriteToCsv(speedSort)
		if err != nil {
			log.Error(err)
			return
		}

		err = services.GenerateConf(speedSort[0].Server, publicKey)
		if err != nil {
			log.Error(err)
			return
		}

	}
}
