package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/urfave/cli"
)

func init() {
	cmdList := cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Action:  listParameters,
		Usage:   "List Parameters",
	}

	cmdGet := cli.Command{
		Name:    "get",
		Aliases: []string{"read"},
		Action:  getParameter,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "decrypt",
				Usage: "(default: false)",
			},
		},
		ArgsUsage: "<name>",
		Usage:     "Read Parameter",
	}

	cmdPut := cli.Command{
		Name:    "put",
		Aliases: []string{"write"},
		Action:  putParameter,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "allowed-pattern",
			},
			cli.StringFlag{
				Name: "description",
			},
			cli.BoolFlag{
				Name:  "encrypt",
				Usage: "(default: false)",
			},
			cli.StringFlag{
				Name:  "encryption-key",
				Usage: "implies --encrypt",
				Value: "alias/aws/ssm",
			},
			cli.BoolFlag{
				Name:  "overwrite",
				Usage: "(default: false)",
			},
		},
		ArgsUsage: "<name> <value>",
		Usage:     "Write Parameter",
	}

	app.Commands = append(app.Commands, cli.Command{
		Name:  "parameter",
		Usage: "Parameter Store",
		Subcommands: []cli.Command{
			cmdList,
			cmdGet,
			cmdPut,
		},
		Category: "RESOURCES",
	})
}

func listParameters(c *cli.Context) error {
	req := svc.DescribeParametersRequest(&ssm.DescribeParametersInput{})
	res, err := req.Send()
	if err != nil {
		return err
	}

	json.NewEncoder(c.App.Writer).Encode(res.Parameters)

	return nil
}

func getParameter(c *cli.Context) error {
	name := c.Args().Get(0)
	decrypt := c.Bool("decrypt")

	req := svc.GetParameterRequest(&ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &decrypt,
	})
	res, err := req.Send()
	if err != nil {
		return err
	}

	json.NewEncoder(c.App.Writer).Encode(res.Parameter)

	return nil
}

func putParameter(c *cli.Context) error {
	name := c.Args().Get(0)
	value := c.Args().Get(1)
	encrypt := c.Bool("encrypt")
	overwrite := c.Bool("overwrite")

	req := svc.PutParameterRequest(&ssm.PutParameterInput{
		Name:      &name,
		Value:     &value,
		Overwrite: &overwrite,
	})

	key := c.String("encryption-key")

	if encrypt || key != "alias/aws/ssm" {
		req.Input.SetKeyId(key)
		req.Input.SetType(ssm.ParameterTypeSecureString)
	} else {
		req.Input.SetType(ssm.ParameterTypeString)
	}

	req.Input.SetDescription(c.String("description"))
	req.Input.SetAllowedPattern(c.String("allowed-pattern"))

	res, err := req.Send()
	if err != nil {
		return err
	}

	json.NewEncoder(c.App.Writer).Encode(map[string]interface{}{
		"Name":    req.Input.Name,
		"Version": res.Version,
	})

	return nil
}
