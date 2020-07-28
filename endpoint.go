package mixin

import (
	"os"
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
	httpClient.SetHostURL(url)
}

func SetBlazeHostURL(url string) {
	blazeHostURL = url
}
