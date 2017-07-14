package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
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

func init() {
	svrconfCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("id", "i", "", "stack id to query")
	updateCmd.Flags().StringP("env", "e", "", "comma-separated key/val pair(s)")
	updateCmd.Flags().StringP("filepath", "p", "", "file path")
}

func update(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.CheckErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	sid := util.GetCliStringFlag(cmd, "id")
	if sid == "" {
		util.CheckErrorExit("stack id cannot be empty", 1)
	}

	// each parameter set is sent separately
	opts := []string{"env", "filepath"}
	for _, opt := range opts {
		var payload string
		val := util.GetCliStringFlag(cmd, opt)

		switch opt {
		case "env":
			in := buildEnvPayload(sid, val)
			if in != "" {
				rm := json.RawMessage(in)
				pl, err := json.Marshal(&rm)
				util.CheckErrorExit(err, 1)
				payload = string(pl)
			}

		case "filepath":
			if val != "" {
				in := buildFilePathPayload(sid, val)
				rm := json.RawMessage(in)
				pl, err := json.Marshal(&rm)
				util.CheckErrorExit(err, 1)
				payload = string(pl)
			}
		}

		if payload == "" {
			continue
		}

		log.Println("payload:", payload)
		c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
		resp, body, errs := c.PutSafe(c.RootUrl+`/alm/serverconfig?stack_id=`+sid, fmt.Sprintf("%s", token), payload)
		if errs != nil {
			if len(errs) > 0 {
				continue
			}
		}

		serr := util.ResponseError(resp, body)
		if serr != "" {
			log.Println("error:", serr)
			continue
		}

		// display return status
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		if status, found := m["status"]; found {
			s := fmt.Sprintf("%s", status)
			if s == "success" {
				line := "[" + resp.Status + "] " + s
				log.Println(line)
				continue
			}
		}

		// or just the raw output
		log.Println(string(body))
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
