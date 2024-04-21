package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/charmbracelet/log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Req struct {
	Method   string
	Host     string
	Timeout  time.Duration
	Data     []byte
	ProxyUrl string
	Header   map[string]string
	Redirect bool
	NoVerify bool
	Gzip     bool
}

var (
	dnsResolverIP        = "1.1.1.1:53" // Google DNS resolver.
	dnsResolverProto     = "udp"        // Protocol to use for the DNS resolver
	dnsResolverTimeoutMs = 5000         // Timeout (ms) for the DNS resolver (optional)
)

func randomUa() string {
	ua := browser.Random()
	return ua
}

func HeaderMap(jsonData ...string) map[string]string {
	var headers = make(map[string]string)
	headers["User-Agent"] = randomUa()
	headers["Connection"] = "close"
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	if len(jsonData) > 0 {
		for _, arg := range jsonData {
			err := json.Unmarshal([]byte(arg), &headers)
			if err != nil {
				log.Error(err)
				return headers
			}
		}
	}
	//headersMap := headerMap{headers: headers}
	return headers
}

// Requests (*int, *[]byte, error)
func (requests Req) Requests() (*http.Response, error) {
	// Specify dns server
	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
				}
				return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
			},
		},
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	http.DefaultTransport.(*http.Transport).DialContext = dialContext

	//proxy
	//proxyUrl := "http://127.0.0.1:8080"
	// var client = &http.Client{Timeout: time.Second * 15}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12,
		},
		ForceAttemptHTTP2: false,
		// From http.DefaultTransport
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if requests.ProxyUrl != "" {
		proxy, _ := url.Parse(requests.ProxyUrl)

		if !requests.NoVerify {
			tr = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		} else {
			tr = &http.Transport{
				Proxy:           http.ProxyURL(proxy),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	} else {
		if requests.NoVerify {
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   requests.Timeout, //超时时间
	}

	req, err := http.NewRequest(requests.Method, requests.Host, bytes.NewBuffer(requests.Data))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	headers := requests.Header
	if len(headers) == 0 {
		headers = HeaderMap()
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if !requests.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("\n[O] disable Redirect")
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		if !strings.Contains(fmt.Sprint(err), "disable Redirect") {
			log.Error(fmt.Sprintf("%s %v", requests.Host, err))
		}
		return resp, err
	}
	if resp.StatusCode != 200 {
		log.Error(fmt.Sprintf("%s Status Code: %d", requests.Host, resp.StatusCode))
	}
	return resp, nil
}
