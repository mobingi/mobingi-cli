package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
}

func update(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.ErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	sid := util.GetCliStringFlag(cmd, "id")
	if sid == "" {
		util.ErrorExit("stack id cannot be empty", 1)
	}

	env := util.GetCliStringFlag(cmd, "env")
	in := buildPayload(sid, env)
	if in == "" {
		util.ErrorExit("Something is wrong with the --env input.", 1)
	}

	rm := json.RawMessage(in)
	payload, err := json.Marshal(&rm)
	if err != nil {
		util.ErrorExit(err.Error(), 1)
	}

	log.Println("payload:", string(payload))
	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	resp, body, errs := c.PutSafe(c.RootUrl+`/alm/serverconfig?stack_id=`+sid, fmt.Sprintf("%s", token), string(payload))
	if errs != nil {
		log.Println("error(s):", errs)
		os.Exit(1)
	}

	serr := util.ResponseError(resp, body)
	if serr != "" {
		util.ErrorExit(serr, 1)
	}

	// display return status
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if status, found := m["status"]; found {
		s := fmt.Sprintf("%s", status)
		if s == "success" {
			line := "[" + resp.Status + "] " + s
			log.Println(line)
			os.Exit(0)
		}
	}

	// or just the raw output
	log.Println(string(body))
}

func buildPayload(sid, env string) string {
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
