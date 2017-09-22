package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func RegistryDeleteTag() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete an image",
		Long: `Delete an image.

Example:

  $ ` + cmdline.Args0() + ` registry delete --image=hello`,
		Run: deleteTag,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("image", "", "image name")
	return cmd
}

func deleteTag(cmd *cobra.Command, args []string) {
	image := cli.GetCliStringFlag(cmd, "image")
	if image == "" {
		d.ErrorExit("image name cannot be empty", 1)
	}

	sess, err := clisession()
	d.ErrorExit(err, 1)

	ensureUserPass(cmd, sess)
	svc := registry.New(sess)
	in := &registry.DeleteImageInput{
		Image: image,
	}

	resp, body, err := svc.DeleteImage(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	if resp.StatusCode/100 != 2 {
		d.Error(resp.Status)
	} else {
		d.Info(resp.Status)
	}

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)
	}
}
