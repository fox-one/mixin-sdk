package mixin

import (
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	blazeHostURL = "wss://mixin-blaze.zeromesh.net/"
)

func init() {
	if url := os.Getenv("MIXIN_API_HOST_URL"); len(url) > 0 {
		SetHostURL(url)
	}

	if url := os.Getenv("MIXIN_BLAZE_HOST_URL"); len(url) > 0 {
		SetBlazeHostURL(url)
	}
}

func SetHostURL(url string) {
	httpClient = httpClient.SetHostURL(url)
}

func SetBlazeHostURL(u string) {
	url, err := url.Parse(u)
	if err != nil {
		log.WithError(err).Errorln("parse blaze host url failed")
		return
	}

	blazeHostURL = url.String()
}
