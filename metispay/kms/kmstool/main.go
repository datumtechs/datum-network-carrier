package main

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/metispay/kms"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"io/ioutil"
	"os"
	"sort"
)

var (
	srcFlag = &cli.StringFlag{
		Name:  "src",
		Usage: "source file",
	}

	appFlags = []cli.Flag{
		flags.ConfigFileFlag,

		altsrc.NewStringFlag(flags.KMS_KeyId),
		altsrc.NewStringFlag(flags.KMS_RegionId),
		altsrc.NewStringFlag(flags.KMS_AccessKeyId),
		altsrc.NewStringFlag(flags.KMS_AccessKeySecret),
		altsrc.NewStringFlag(flags.Chain),

		srcFlag,
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

func buildKMS(ctx *cli.Context) (kms.KmsService, error) {
	if err := flags.LoadFlagsFromConfig(ctx, appFlags); err != nil {
		return nil, err
	}
	if ctx.IsSet(flags.KMS_KeyId.Name) && ctx.IsSet(flags.KMS_RegionId.Name) && ctx.IsSet(flags.KMS_AccessKeyId.Name) && ctx.IsSet(flags.KMS_AccessKeySecret.Name) {

		kmsConfig := &kms.Config{
			KeyId:           ctx.String(flags.KMS_KeyId.Name),
			RegionId:        ctx.String(flags.KMS_RegionId.Name),
			AccessKeyId:     ctx.String(flags.KMS_AccessKeyId.Name),
			AccessKeySecret: ctx.String(flags.KMS_AccessKeySecret.Name),
		}
		alikms := &kms.AliKms{Config: kmsConfig}
		return alikms, nil
	}
	return nil, errors.New("cannot load KMS configuration")
}
func encrypt(ctx *cli.Context) error {
	if kms, err := buildKMS(ctx); err != nil {
		fmt.Printf("build KMS error:%v\n", err)
		return err
	} else {
		content, err := ioutil.ReadFile(ctx.String(srcFlag.Name))
		if err != nil {
			fmt.Printf("read source file error:%v\n", err)
			return err
		}

		encoded, err := kms.Encrypt(content)
		if err != nil {
			fmt.Printf("encrypt source file error:%v\n", err)
			return err
		}

		err = ioutil.WriteFile(ctx.String(srcFlag.Name)+".enc", encoded, 0644)
		if err != nil {
			fmt.Printf("write dest file error:%v\n", err)
			return err
		}
		return err
	}
}

func decrypt(ctx *cli.Context) error {
	if kms, err := buildKMS(ctx); err != nil {
		fmt.Printf("build KMS error:%v\n", err)
		return err
	} else {
		content, err := ioutil.ReadFile(ctx.String(srcFlag.Name))
		if err != nil {
			fmt.Printf("read source file error:%v\n", err)
			return err
		}
		plaintext, err := kms.Decrypt(content)
		if err != nil {
			fmt.Printf("decrypt source file error:%v\n", err)
			return err
		}
		err = ioutil.WriteFile(ctx.String(srcFlag.Name)+".plain", plaintext, 0644)
		if err != nil {
			fmt.Printf("write dest file error:%v\n", err)
			return err
		}
		return err
	}
}
