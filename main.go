package main

import (
	"context"
	_ "ddns/log"
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/dns/v1"
)

var PROJECT_ID = ""
var MANAGED_ZONE = ""
var NAME = ""
var GOOGLE_APPLICATION_CREDENTIALS = ""
var GOOGLE_SERVICE_ACCOUNT_KEY_JSON_BASE64 = ""

func init() {
	PROJECT_ID = os.Getenv("PROJECT_ID")
	MANAGED_ZONE = os.Getenv("MANAGED_ZONE")
	NAME = os.Getenv("NAME")
	GOOGLE_APPLICATION_CREDENTIALS = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	GOOGLE_SERVICE_ACCOUNT_KEY_JSON_BASE64 = os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_JSON_BASE64")

	saJson, err := base64.StdEncoding.DecodeString(GOOGLE_SERVICE_ACCOUNT_KEY_JSON_BASE64)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = os.WriteFile(GOOGLE_APPLICATION_CREDENTIALS, saJson, 0644)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

}

func main() {

	update()

	for range time.Tick(time.Minute) {

		update()
	}
}

func update() {

	log.Info().
		Str("PROJECT_ID", PROJECT_ID).
		Str("MANAGED_ZONE", MANAGED_ZONE).
		Str("NAME", NAME).
		Msg("Start")

	eip := expect()
	aip := actual()

	log.Info().
		Str("eip", eip).
		Str("aip", aip).
		Msg("")

	if eip == aip {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dnsService, err := dns.NewService(ctx)

	if err != nil {
		log.Fatal().Err(err).Send()
	}
	call := dnsService.ResourceRecordSets.Create(
		PROJECT_ID,
		MANAGED_ZONE,
		&dns.ResourceRecordSet{
			Kind:    "dns#resourceRecordSet",
			Name:    NAME,
			Rrdatas: []string{aip + "."},
			Ttl:     300,
			Type:    "A",
		})

	record, err := call.Do()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Interface("record", record).Msg("")
}

func actual() string {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return string(b)
}

func expect() string {
	ips, err := net.LookupIP(NAME)
	if err != nil {
		log.Info().Err(err).Send()
		return ""
	}

	if len(ips) == 0 {
		return ""
	}

	return ips[0].To4().String()
}
