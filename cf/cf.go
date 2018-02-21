package cf

import (
	"log"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"sevki.org/cf-controller/dns"
	"sevki.org/x/pretty"
	"sevki.org/x/reconcile"
)

type action int

type cf struct {
	api *cloudflare.API

	zone   string
	zoneID string
}

// New returns the state of a cloudflare zone
func New(apiKey string, apiEmail string, zone string) (reconcile.State, error) {
	api, err := cloudflare.New(apiKey, apiEmail)
	if err != nil {
		return nil, err
	}
	zoneID, err := api.ZoneIDByName(zone)
	if err != nil {
		return nil, err
	}
	return &cf{api: api, zone: zone, zoneID: zoneID}, nil
}

func (cf *cf) Add(key string, v interface{}) {
	switch x := v.(type) {
	case dns.Record:
		if resp, err := cf.api.CreateDNSRecord(cf.zoneID, cloudflare.DNSRecord(x)); err != nil {
			log.Println(err)
		} else {
			log.Println(pretty.JSON(resp))
		}

	case cloudflare.DNSRecord:
		if _, err := cf.api.CreateDNSRecord(cf.zoneID, x); err != nil {
			log.Println(err)
		}
	case []cloudflare.DNSRecord:
		for _, record := range x {
			if _, err := cf.api.CreateDNSRecord(cf.zoneID, record); err != nil {
				log.Println(err)
			}
		}
	default:
		log.Println("cf: value must be a record")
	}
}

func (cf *cf) Delete(key string) {
	records, err := cf.api.DNSRecords(cf.zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return
	}
	for _, record := range records {
		if record.Content == key {
			cf.api.DeleteDNSRecord(cf.zone, record.ID)
		}
	}
}

func (cf *cf) Get(key string) interface{} {
	records, err := cf.api.DNSRecords(cf.zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return nil
	}
	for _, record := range records {
		if record.Content == key {
			return dns.Record(record)
		}
	}
	return nil
}

func (cf *cf) Update(key string, v interface{}) {
	if record, ok := v.(cloudflare.DNSRecord); ok {
		cf.api.UpdateDNSRecord(cf.zoneID, key, record)

	}
}

func (cf *cf) Walk(wf reconcile.StateWalkFunc) {
	records, err := cf.api.DNSRecords(cf.zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return
	}
	for _, record := range records {
		wf(record.Content, dns.Record(record))
	}
}
