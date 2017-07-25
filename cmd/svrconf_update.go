package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
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

	cmd.Flags().StringP("id", "i", "", "stack id to query")
	cmd.Flags().StringP("env", "e", "", "comma-separated key/val pair(s)")
	cmd.Flags().StringP("filepath", "p", "", "file path")
	return cmd
}

func update(cmd *cobra.Command, args []string) {
	sid := cli.GetCliStringFlag(cmd, "id")
	if sid == "" {
		check.ErrorExit("stack id cannot be empty", 1)
	}

	// each parameter set is sent separately
	opts := []string{"env", "filepath"}
	for _, opt := range opts {
		var payload string
		val := cli.GetCliStringFlag(cmd, opt)

		switch opt {
		case "env":
			in := buildEnvPayload(sid, val)
			if in != "" {
				rm := json.RawMessage(in)
				pl, err := json.Marshal(&rm)
				check.ErrorExit(err, 1)
				payload = string(pl)
			}

		case "filepath":
			if val != "" {
				in := buildFilePathPayload(sid, val)
				rm := json.RawMessage(in)
				pl, err := json.Marshal(&rm)
				check.ErrorExit(err, 1)
				payload = string(pl)
			}
		}

		if payload == "" {
			continue
		}

		d.Info("payload:", payload)
		c := client.NewGrClient(client.NewApiConfig(cmd))
		resp, body, errs := c.Put(`/alm/serverconfig?stack_id=`+sid, payload)
		if errs != nil {
			if len(errs) > 0 {
				continue
			}
		}

		serr := check.ResponseError(resp, body)
		if serr != "" {
			d.Error(serr)
			continue
		}

		// display return status
		var m map[string]interface{}
		_ = json.Unmarshal(body, &m)
		if status, found := m["status"]; found {
			s := fmt.Sprintf("%s", status)
			if s == "success" {
				line := "[" + resp.Status + "] " + s
				d.Info(line)
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
