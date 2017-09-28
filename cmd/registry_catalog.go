package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
)

type catalog struct {
	Repositories []string `json:"repositories"`
}

func RegistryCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "list catalog images",
		Long: `List catalog images. Note that this command will probably
take some time to complete.

Example:

  $ ` + cmdline.Args0() + ` registry catalog --username=foo --password=bar`,
		Run: printCatalog,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	return cmd
}

func printCatalog(cmd *cobra.Command, args []string) {
	service := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	if scope == "" {
		scope = "registry:catalog:*"
	}

	sess, err := clisession()
	cli.ErrorExit(err, 1)

	ensureUserPass(cmd, sess)
	svc := registry.New(sess)
	in := &registry.GetUserCatalogInput{
		Service: service,
		Scope:   scope,
	}

	resp, _, list, err := svc.GetUserCatalog(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	if len(list) > 0 {
		d.Info("Catalog list for user:", sess.Config.Username)
		for _, v := range list {
			fmt.Println(v)
		}
	} else {
		d.Info("Catalog is empty for user", sess.Config.Username+".")
	}
}
