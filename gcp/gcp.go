package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"sevki.org/cf-controller/dns"
	"sevki.org/x/names"
	"sevki.org/x/reconcile"
)

const ttl = 120

type gcp struct {
	client          *compute.Service
	projectID, zone string
}

// New returns the instance state of a project and zone
func New(keyFile, projectID, zone string) (reconcile.State, error) {
	bytz, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(bytz, compute.ComputeReadonlyScope, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	c := conf.Client(context.Background())
	client, err := compute.New(c)
	if err != nil {
		return nil, err
	}

	return &gcp{
		client:    client,
		projectID: projectID,
		zone:      zone,
	}, nil
}

func (g *gcp) Add(string, interface{})    { return } // We can't add shit to GCP, nor do we want to.
func (g *gcp) Delete(string)              {}         // We can't delete shit from GCP, nor do we want to.
func (g *gcp) Update(string, interface{}) {}         // We can't delete shit from GCP, nor do we want to.
func (g *gcp) Get(key string) interface{} {
	list, err := g.client.Instances.List(g.projectID, g.zone).Do()
	if err != nil {
		return nil
	}
	for _, v := range list.Items {
		id := fmt.Sprintf("%x", v.Id)
		_, z := path.Split(v.Zone)
		ip := v.NetworkInterfaces[0].AccessConfigs[0].NatIP
		val := cloudflare.DNSRecord{
			Name:    "sevki.cloud",
			Type:    "A",
			Content: ip,
			TTL:     ttl,
			Meta: map[string]interface{}{
				"zone":        z,
				"instance_id": id,
			},
		}
		if key == ip {
			return dns.Record(val)
		}
	}
	return nil
}
func (g *gcp) fqdn(s string) string {
	name := strings.ToLower(names.For(s))
	return fmt.Sprintf("%s.sevki.cloud", name)
}
func (g *gcp) Walk(wf reconcile.StateWalkFunc) {
	list, err := g.client.Instances.List(g.projectID, g.zone).Do()
	if err != nil {
		return
	}
	for _, v := range list.Items {
		ip := v.NetworkInterfaces[0].AccessConfigs[0].NatIP
		val := cloudflare.DNSRecord{
			Name:    "sevki.cloud",
			Type:    "A",
			Content: ip,
			TTL:     ttl,
			Meta: map[string]interface{}{
				"zone":        v.Zone,
				"instance_id": fmt.Sprintf("%x", v.Id),
			},
		}
		wf(ip, dns.Record(val))
	}

}
