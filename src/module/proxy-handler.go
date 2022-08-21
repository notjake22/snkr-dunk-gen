package module

import (
	"encoding/base64"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (t *SnkrDunkTask) proxyHandler() error {
	proxyU := getProxy()

	if proxyU != "" {
		parts := strings.Split(proxyU, ":")
		if len(parts) == 4 {
			newProxy := fmt.Sprintf("%s:%s@%s:%s", parts[2], parts[3], parts[0], strings.Replace(parts[1], "]", "", -1))
			proxyUrl, err := url.Parse("http://" + newProxy)
			if err != nil {
				sentry.CaptureException(err)
				return err
			}
			ipPort, err := url.Parse(fmt.Sprintf("http://%s:%s", proxyUrl.Hostname(), proxyUrl.Port()))
			if err != nil {
				sentry.CaptureException(err)
				return err
			}

			pw, _ := proxyUrl.User.Password()
			userpass := fmt.Sprintf("%s:%s", proxyUrl.User.Username(), pw)
			t.client = http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(ipPort),
					ProxyConnectHeader: http.Header{
						"Proxy-Authorization": {fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userpass)))},
					},
				},
			}

		} else if len(parts) == 2 {
			newProxy := fmt.Sprintf("%s:%s", parts[0], strings.Replace(parts[1], "]", "", -1))
			proxyUrl, err := url.Parse("http://" + newProxy)
			if err != nil {
				sentry.CaptureException(err)
				return err
			}
			ipPort, err := url.Parse(fmt.Sprintf("http://%s:%s", proxyUrl.Hostname(), proxyUrl.Port()))
			if err != nil {
				sentry.CaptureException(err)
				return err
			}

			t.client = http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(ipPort),
				},
			}

		} else {
			log.Println("Bad proxy, using local")
		}
	} else {
		log.Println("No proxies found in proxies.txt continuing with local host")
	}
	return nil
}

func getProxy() string {
	file, err := os.ReadFile("proxies.txt")
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	proxies := string(file)
	if proxies == "" {
		return ""
	}
	proxies = strings.Replace(proxies, "\r", "", -1)
	proxiesArray := strings.Split(proxies, "\n")

	return proxiesArray[rand.Intn(len(proxiesArray))]
}
