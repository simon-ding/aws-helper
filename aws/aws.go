package aws

import (
	"aws-helper/log"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/pkg/errors"
)

func CreateLightsailClient(instanceName, region, cfgDir string) (*Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigFiles(
			[]string{cfgDir},
		),
		//types.RegionNameUsWest2
		config.WithRegion(region),
	)
	if err != nil {
		return nil, errors.Wrap(err, "load aws config")
	}
	client := lightsail.NewFromConfig(cfg)
	c := &Client{Client: client, instanceName: instanceName}
	return c, nil
}

type Client struct {
	*lightsail.Client
	instanceName string
}

func (c *Client) CreateIp(name string) error {
	in := lightsail.AllocateStaticIpInput{
		StaticIpName: &name,
	}
	_, err := c.AllocateStaticIp(context.TODO(), &in)
	if err != nil {
		return errors.Wrap(err, "new ip")
	}

	log.Infof("new ip created: %+v", name)
	return nil
}

func (c *Client) DetachIp(name string) error {
	in := lightsail.DetachStaticIpInput{StaticIpName: &name}

	_, err := c.Client.DetachStaticIp(context.TODO(), &in)
	if err != nil {
		return errors.Wrap(err, "detach ip")
	}
	log.Infof("ip %s detach success", name)
	return nil
}

func (c *Client) AttachIp(name string) error {
	in := lightsail.AttachStaticIpInput{
		InstanceName: &c.instanceName,
		StaticIpName: &name,
	}

	_, err := c.Client.AttachStaticIp(context.TODO(), &in)
	if err != nil {
		return errors.Wrap(err, "attach ip")
	}
	log.Infof("attach ip %s to instance %s success", name, c.instanceName)
	return nil
}

func (c *Client) GetAttachedIp() (string, string, error) {
	in := lightsail.GetStaticIpsInput{}
	out, err := c.Client.GetStaticIps(context.TODO(), &in)
	if err != nil {
		return "", "", errors.Wrap(err, "get static ips")
	}
	for _, ip := range out.StaticIps {
		if ip.IsAttached != nil && *ip.IsAttached {
			return *ip.Name, *ip.IpAddress, err
		}
	}
	return "", "", fmt.Errorf("no attached ip")
}

func (c *Client) DeleteIp(ipName string) error {
	in := lightsail.ReleaseStaticIpInput{StaticIpName: &ipName}
	_, err := c.Client.ReleaseStaticIp(context.TODO(), &in)
	if err != nil {
		return errors.Wrap(err, "release ip")
	}
	log.Infof("release ip %s success", ipName)
	return nil
}

func (c *Client) AttachNewIP() error {
	name, ip, err := c.GetAttachedIp()
	if err != nil {
		return err
	}
	log.Infof("current attached ip: %s(%v)", name, ip)
	err = c.DetachIp(name)
	if err != nil {
		return err
	}
	err = c.DeleteIp(name)
	if err != nil {
		return err
	}

	newName := "ip_" + RandStringRunes(6)
	err = c.CreateIp(newName)
	if err != nil {
		return err
	}

	err = c.AttachIp(newName)
	if err != nil {
		return err
	}
	log.Infof("attach %v a new ip %s", c.instanceName, newName)
	return err
}
