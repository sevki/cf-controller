package main

import (
	"log"
	"os"
	"time"

	"sevki.org/cf-controller/cf"
	"sevki.org/cf-controller/gcp"
	"sevki.org/x/reconcile"
)

func doe(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// "us-east1-b"
func main() {
	gcp, err := gcp.New(os.Getenv("GCP_KEY_FILE"), os.Getenv("GCP_PROJECT"), os.Getenv("GCP_ZONE"))
	if err != nil {
		log.Fatal(err)
	}
	cf, err := cf.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"), os.Getenv("CF_ZONE"))
	if err != nil {
		log.Fatal(err)
	}
	reconLoop(cf, gcp)
}

func reconLoop(cloudflare reconcile.State, cloudprovider reconcile.State) {
	for {
		reconcile.Reconcile(cloudflare, cloudprovider)
		time.Sleep(time.Second * 30)
	}
}
