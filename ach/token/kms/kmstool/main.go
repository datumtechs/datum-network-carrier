package main

import (
	"errors"
	"fmt"
	kms2 "github.com/datumtechs/datum-network-carrier/ach/token/kms"
	"github.com/datumtechs/datum-network-carrier/common/flags"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"io/ioutil"
	"os"
	"sort"
)

var (
	keystoreFlag = &cli.StringFlag{
		Name:  "keystore",
		Usage: "keystore file",
	}

	inputFlag = &cli.StringFlag{
		Name:  "input",
		Usage: "input string",
	}

	appFlags = []cli.Flag{
		flags.ConfigFileFlag,

		altsrc.NewStringFlag(flags.KMSKeyId),
		altsrc.NewStringFlag(flags.KMSRegionId),
		altsrc.NewStringFlag(flags.KMSAccessKeyId),
		altsrc.NewStringFlag(flags.KMSAccessKeySecret),
		altsrc.NewStringFlag(flags.BlockChain),

		keystoreFlag,
		inputFlag,
	}
)

func main() {
	app := &cli.App{
		//这样写，表示参数是整个app的参数, 输入时格式是：kmstool --config-file xxx -src xxx encrypt
		Flags: appFlags,
		Commands: []*cli.Command{
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "encrypt source file by KMS",
				//这样写，表示参数是命令encrypt/decrypt的参数, 输入时格式是：kmstool encrypt --config-file xxx -src xxx
				//Flags:   appFlags,
				Action: func(ctx *cli.Context) error {
					return encrypt(ctx)
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "decrypt source file by KMS",
				//这样写，表示参数是命令encrypt/decrypt的参数, 输入时格式是：kmstool encrypt --config-file xxx -src xxx
				//Flags:   appFlags,
				Action: func(c *cli.Context) error {
					return decrypt(c)
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("error: %v", err)
	}
}

func buildKMS(ctx *cli.Context) (kms2.KmsService, error) {
	if err := flags.LoadFlagsFromConfig(ctx, appFlags); err != nil {
		return nil, err
	}
	if ctx.IsSet(flags.KMSKeyId.Name) && ctx.IsSet(flags.KMSRegionId.Name) && ctx.IsSet(flags.KMSAccessKeyId.Name) && ctx.IsSet(flags.KMSAccessKeySecret.Name) {

		kmsConfig := &kms2.Config{
			KeyId:           ctx.String(flags.KMSKeyId.Name),
			RegionId:        ctx.String(flags.KMSRegionId.Name),
			AccessKeyId:     ctx.String(flags.KMSAccessKeyId.Name),
			AccessKeySecret: ctx.String(flags.KMSAccessKeySecret.Name),
		}
		alikms := &kms2.AliKms{Config: kmsConfig}
		return alikms, nil
	}
	return nil, errors.New("cannot load KMS configuration")
}

func encrypt(ctx *cli.Context) error {
	if kms, err := buildKMS(ctx); err != nil {
		fmt.Printf("build KMS error:%v\n", err)
		return err
	} else {
		var content string
		if ctx.IsSet(inputFlag.Name) {
			content = ctx.String(inputFlag.Name)
		} else if ctx.IsSet(keystoreFlag.Name) {
			contentBytes, err := ioutil.ReadFile(ctx.String(keystoreFlag.Name))
			if err != nil {
				fmt.Printf("read source file error:%v\n", err)
				return err
			}
			content = string(contentBytes)
		} else {
			fmt.Println("missing --input or --keystore")
		}

		encoded, err := kms.Encrypt(content)
		if err != nil {
			fmt.Printf("encrypt plaintext error:%v\n", err)
			return err
		}

		fmt.Printf("encrypt resulst:\n\n%s\n", encoded)

		if ctx.IsSet(keystoreFlag.Name) {
			err = ioutil.WriteFile(ctx.String(keystoreFlag.Name)+".enc", []byte(encoded), 0644)
			if err != nil {
				fmt.Printf("write dest file error:%v\n", err)
				return err
			}
		}

		return nil
	}
}

func decrypt(ctx *cli.Context) error {
	if kms, err := buildKMS(ctx); err != nil {
		fmt.Printf("build KMS error:%v\n", err)
		return err
	} else {
		var content string
		if ctx.IsSet(inputFlag.Name) {
			content = ctx.String(inputFlag.Name)
		} else if ctx.IsSet(keystoreFlag.Name) {
			contentBytes, err := ioutil.ReadFile(ctx.String(keystoreFlag.Name))
			if err != nil {
				fmt.Printf("read source file error:%v\n", err)
				return err
			}
			content = string(contentBytes)
		} else {
			fmt.Println("missing --input or --keystore")
		}

		plaintext, err := kms.Decrypt(content)
		fmt.Printf("decrypt resulst:\n\n%s\n", plaintext)

		if err != nil {
			fmt.Printf("decrypt error:%v\n", err)
			return err
		}
		if ctx.IsSet(keystoreFlag.Name) {
			err = ioutil.WriteFile(ctx.String(keystoreFlag.Name)+".plain", []byte(plaintext), 0644)
			if err != nil {
				fmt.Printf("write dest file error:%v\n", err)
				return err
			}
		}

		return nil
	}
}
