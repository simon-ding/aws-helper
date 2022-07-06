package main

import (
	"aws-helper/aws"
	"aws-helper/cfg"
	"aws-helper/cloudflare"
	"aws-helper/log"

	"github.com/pkg/errors"
)

func main() {
	cfg, err := cfg.LoadConfig()
	if err != nil {
		log.Errorf("load config error: %v", err)
		return
	}

	if err := doSwapIp(cfg); err != nil {
		log.Errorf("swap ip error: %v", err)
		return
	}
}

func doSwapIp(cfg *cfg.Config) error {
	cf, err := cloudflare.NewClient(cfg.Cloudflare)
	if err != nil {
		return errors.Wrap(err, "new cloudflare")
	}

	awsCfgDir := "."
	aws.WriteAwsConfig(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, awsCfgDir)
	lightsail, err := aws.CreateLightsailClient(cfg.AWS.NodeName, cfg.AWS.Region, awsCfgDir)
	if err != nil {
		return errors.Wrap(err, "create aws client")
	}

	err = lightsail.AttachNewIP()
	if err != nil {
		return errors.Wrap(err, "attach new ip")
	}

	_, addr, err := lightsail.GetAttachedIp()
	if err != nil {
		return errors.Wrap(err, "get attach ip")
	}
	err = cf.UpdateDnsRecord(cfg.AWS.DnsName, addr)
	if err != nil {
		return errors.Wrap(err, "update dns record")
	}

	return nil
}
