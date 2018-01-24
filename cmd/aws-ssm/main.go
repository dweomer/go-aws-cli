package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/urfave/cli"
)

var (
	app = cli.NewApp()
	cfg aws.Config
	svc *ssm.SSM
)

func init() {
	app.Usage = "AWS (Simple) Systems Manager CLI"
	app.Version = Version
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "profile",
			EnvVar: external.AWSDefaultProfileEnvVar,
			Usage:  fmt.Sprintf("$%s", external.AWSProfileEnvVar),
			Value: func() string {
				if val, found := os.LookupEnv(external.AWSProfileEnvVar); found {
					return val
				}
				if val, found := os.LookupEnv(external.AWSDefaultProfileEnvVar); found {
					return val
				}
				return external.DefaultSharedConfigProfile
			}(),
		},
		cli.StringFlag{
			Name:   "region",
			EnvVar: external.AWSDefaultRegionEnvVar,
			Usage:  fmt.Sprintf("$%s", external.AWSRegionEnvVar),
			Value: func() string {
				if val, found := os.LookupEnv(external.AWSRegionEnvVar); found {
					return val
				}
				return os.Getenv(external.AWSDefaultRegionEnvVar)
			}(),
		},
	}
	app.Before = func(*cli.Context) (err error) {
		if cfg, err = external.LoadDefaultAWSConfig(); err == nil {
			svc = ssm.New(cfg)
		}
		return err
	}
}

func main() {
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	app.RunAndExitOnError()
}
