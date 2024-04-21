package util

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/oschwald/geoip2-golang"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

type Release struct {
	HTMLURL string `json:"html_url"`
}

func checkLastReleases() (*string, error) {
	owner := "P3TERX"      // 替换为仓库所有者的用户名或组织名
	repo := "GeoLite.mmdb" // 替换为仓库名称

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	releaseVersion := strings.Split(release.HTMLURL, "/")

	downloadUrl := fmt.Sprintf("https://ghproxy.com/https://github.com/P3TERX/GeoLite.mmdb/releases/download/%s/GeoLite2-City.mmdb", releaseVersion[len(releaseVersion)-1])
	return &downloadUrl, nil
}

func DownloadGeoLite() error {
	// Create output file
	fileURL, err := checkLastReleases()
	if err != nil {
		return err
	}
	out, err := os.Create(GeoLite2CityPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	// Send HTTP GET request to obtain file content
	resp, err := http.Get(*fileURL)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Download GeoLite.mmdb")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}(resp.Body)

	// Copy the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Download Success")
	return nil
}

func Ip2Location(target string) (*string, error) {
	db, err := geoip2.Open(GeoLite2CityPath)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer func(db *geoip2.Reader) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)
	ip := net.ParseIP(target)
	record, err := db.City(ip)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	country := record.RegisteredCountry.Names["en"]
	return &country, nil
}
