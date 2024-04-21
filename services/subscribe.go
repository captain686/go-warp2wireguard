package services

import (
	_ "embed"
	"fmt"
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"net/http"
	"strings"
	"text/template"
)

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

func flashSubs(publicKey string) util.Clash {
	var (
		speedSort []util.Speed
		err       error
	)
	clash := util.InitClash()
	if util.FileCheck(util.OutFilePath) {
		speedSort, err = util.ReadFromCsv(util.OutFilePath)
		if err != nil {
			//log.Error(err)
			return clash
		}
	}
	count := 1
	for _, v := range speedSort {
		if count > 60 {
			break
		}
		node := v.Server
		clash = *clashGenerate(publicKey, node, clash)
		count++
	}
	return clash
}

func Subscribe(publicKey, clashTemplate string) {
	tmpl, err := template.New("clash").Parse(clashTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	handleTemplate := func(w http.ResponseWriter, r *http.Request) {
		// Create data object
		// Perform template rendering and pass data
		clash := flashSubs(publicKey)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		err := tmpl.Execute(w, clash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Register handler function and start HTTP server
	http.HandleFunc("/", handleTemplate)
	log.Info("Http Server Listen in http://127.0.0.1:8888")
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Error(err)
		return
	}

}
