package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mobingi/mobingi-cli/client"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/private/debug"
	"github.com/spf13/cobra"
)

func ServerConfigUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update current server config",
		Long: `Show current server config. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values:

  json (default), raw

Example on how to input environment variables via --env option:

  --env=KEY1:value1,KEY2:value2,KEYX:valuex

  or (enclose in double quotes when you have whitespaces)

  --env="KEY1:value1, KEY2:value2, KEYX:valuex"`,
		Run: update,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringP("id", "i", "", "stack id to query")
	cmd.Flags().StringP("env", "e", "", "comma-separated key/val pair(s)")
	cmd.Flags().StringP("filepath", "p", "", "file path")
	return cmd
}

func update(cmd *cobra.Command, args []string) {
	sid := cli.GetCliStringFlag(cmd, "id")
	if sid == "" {
		d.ErrorExit("stack id cannot be empty", 1)
	}

	// each parameter set is sent separately
	opts := []string{"env", "filepath"}
	for _, opt := range opts {
		var payload []byte
		val := cli.GetCliStringFlag(cmd, opt)

		switch opt {
		case "env":
			in := buildEnvPayload(sid, val)
			if in != "" {
				rm := json.RawMessage(in)
				p, err := json.Marshal(&rm)
				d.ErrorExit(err, 1)
				payload = p
			}

		case "filepath":
			if val != "" {
				in := buildFilePathPayload(sid, val)
				rm := json.RawMessage(in)
				p, err := json.Marshal(&rm)
				d.ErrorExit(err, 1)
				payload = p
			}
		}

		if len(payload) == 0 {
			continue
		}

		d.Info("payload:", string(payload))
		c := client.NewClient(client.NewApiConfig(cmd))
		_, body, err := c.AuthPut(`/alm/serverconfig?stack_id=`+sid, payload)
		if err != nil {
			continue
		}

		// display return status
		var m map[string]interface{}
		_ = json.Unmarshal(body, &m)
		if status, found := m["status"]; found {
			s := fmt.Sprintf("%s", status)
			if s == "success" {
				d.Info(s)
				continue
			}
		}

		// or just the raw output
		d.Info(string(body))
	}
}

func buildEnvPayload(sid, env string) string {
	cnt := 0
	payload := `{"stack_id":"` + sid + `",`

	// check if delete all
	if env == "null" {
		payload += `"envvars":{}}`
		return payload
	}

	if env != "" {
		line := `"envvars":{`
		envs := strings.Split(env, ",")
		for i, s := range envs {
			kv := strings.Split(s, ":")
			if len(kv) == 2 {
				line += `"` + strings.TrimSpace(kv[0]) + `":"` + strings.TrimSpace(kv[1]) + `"`
				cnt += 1
			}

			if i < len(envs)-1 {
				line += `,`
			}
		}

		line += `}`
		payload += line
	}

	payload += `}`
	if cnt == 0 {
		return ""
	}

	return payload
}

func buildFilePathPayload(sid, fp string) string {
	return `{"stack_id":"` + sid + `","filepath":"` + fp + `"}`
}
