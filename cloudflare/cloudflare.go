package cloudflare

import (
	"aws-helper/cfg"
	"aws-helper/log"
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"
)

func NewClient(cfg cfg.Cloudflare) (*Client, error) {
	api, err := cloudflare.NewWithAPIToken(cfg.ApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "new token")
	}
	api.AccountID = cfg.AccountId

	return &Client{
		cfg: &cfg,
		api: api,
	}, nil
}

type Client struct {
	cfg *cfg.Cloudflare
	api *cloudflare.API
}

func (c *Client) UpdateDnsRecord(name, newIp string) error {
	log.Infof("update cloudflare record ip to %v", newIp)
	recs, err := c.api.DNSRecords(context.TODO(), c.cfg.ZoneId, cloudflare.DNSRecord{Name: name})
	if err != nil {
		return errors.Wrap(err, "get records")
	}
	if len(recs) == 0 {
		rsp, err := c.api.CreateDNSRecord(context.TODO(), c.cfg.ZoneId, cloudflare.DNSRecord{
			Type:    "A",
			Name:    name,
			Content: newIp,
		})
		if err != nil {
			return errors.Wrap(err, "create record")
		}
		log.Infof("create %s success, result %+v", name, rsp)
		return nil
	}
	r := recs[0]
	r.Content = newIp
	err = c.api.UpdateDNSRecord(context.TODO(), c.cfg.ZoneId, r.ID, r)
	if err != nil {
		return errors.Wrap(err, "update record")
	}
	return nil
}
