package services

import (
	_ "embed"
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"net/http"
	"text/template"
)

func Subscribe(clash util.Clash, clashTemplate string) {
	tmpl, err := template.New("clash").Parse(clashTemplate)
	if err != nil {
		log.Error(err)
		return
	}
	handleTemplate := func(w http.ResponseWriter, r *http.Request) {
		// Create data object
		// Perform template rendering and pass data
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
