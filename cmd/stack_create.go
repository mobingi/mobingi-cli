package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/filetype"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

var (
	usedb        bool
	usecache     bool
	readreplica1 bool
	readreplica2 bool
	readreplica3 bool
	readreplica4 bool
	readreplica5 bool
)

func StackCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a stack",
		Long: `Create a stack. For now, 'aws' is the only supported vendor.

[V3 API]

The only option you need for v3 is '--alm-template=path_to_template_file'.

Examples:

  $ ` + cmdline.Args0() + " stack create --alm-template=`echo $HOME`/tmpl.json" + `
  $ ` + cmdline.Args0() + ` stack create --alm-template=/home/user/tmpl.yaml"

[V2 API]

You can get your credential id using the command:

  $ ` + cmdline.Args0() + ` creds list

If credential id is empty, cli will attempt to get the list using the
command above and use the first one in the list (if more than one).

For --image, omit the domain part when pulling images from hub.docker.com:

  greyltc/lamp

Otherwise, specify the full path:

  registry.mobingi.com/wayland/lamp

As an example for --spot-range, if you have a total of 20 instances running
in the autoscaling group and your spot range is set to 50 (50%), then there
will be a fleet of 10 spot instances and 10 on-demand instances.

Example:

  $ ` + cmdline.Args0() + ` stack create --nickname=sample`,
		Run: createStack,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("alm-template", "", "`path` to alm template file")
	cmd.Flags().String("vendor", "aws", "vendor/provider")
	cmd.Flags().String("cred", "", "credential id")
	cmd.Flags().String("region", "ap-northeast-1", "region code")
	cmd.Flags().String("nickname", "", "stack [nick]name")
	cmd.Flags().String("arch", "art_elb", "single stack: art_single; load balanced: art_elb")
	cmd.Flags().String("type", "m3.medium", "server type")
	cmd.Flags().String("image", "mobingi/ubuntu-apache2-php7:7.1", "docker registry path")
	cmd.Flags().String("dhub-user", "", "docker hub username if private repo")
	cmd.Flags().String("dhub-pass", "", "docker hub password if private repo")
	cmd.Flags().Int("min", 2, "min auto scale group instance when arch is art_elb")
	cmd.Flags().Int("max", 10, "max auto scale group instance when arch is art_elb")
	cmd.Flags().Int("spot-range", 50, "spot instance percentage to deploy")
	cmd.Flags().String("code", "github.com/mobingilabs/default-site-php", "git repository url; can be updated at any time")
	cmd.Flags().String("code-ref", "master", "git repo branch")
	cmd.Flags().String("code-privkey", "", "private key if repo is private")
	cmd.Flags().BoolVar(&usedb, "usedb", false, "if you want to use database")
	cmd.Flags().String("dbengine", "", "valid values: db_mysql, db_postgresql (requires --usedb)")
	cmd.Flags().String("dbtype", "", "db instance class/type (requires --usedb)")
	cmd.Flags().String("dbstorage", "", "db storage in GB, between 5 to 6144 (requires --usedb)")
	cmd.Flags().BoolVar(&readreplica1, "dbread-replica1", false, "read replica 1 (requires --usedb)")
	cmd.Flags().BoolVar(&readreplica2, "dbread-replica2", false, "read replica 2 (requires --usedb)")
	cmd.Flags().BoolVar(&readreplica3, "dbread-replica3", false, "read replica 3 (requires --usedb)")
	cmd.Flags().BoolVar(&readreplica4, "dbread-replica4", false, "read replica 4 (requires --usedb)")
	cmd.Flags().BoolVar(&readreplica5, "dbread-replica5", false, "read replica 5 (requires --usedb)")
	cmd.Flags().BoolVar(&usecache, "use-elasticache", false, "if you want to use elasticache")
	cmd.Flags().String("elasticache-engine", "", "valid values: Redis, Memcached (requires --use-elasticache)")
	cmd.Flags().String("elasticache-nodetype", "", "elasticache node size; ie. cache.r3.large (requires --use-elasticache)")
	cmd.Flags().String("elasticache-nodecount", "", "if redis, 1 to 6; if memcached, 1 to 20 (requires --use-elasticache)")
	return cmd
}

func createStack(cmd *cobra.Command, args []string) {
	almt := cli.GetCliStringFlag(cmd, "alm-template")
	if almt != "" {
		createAlmStack(cmd)
		return
	}

	vendor := cli.GetCliStringFlag(cmd, "vendor")
	region := cli.GetCliStringFlag(cmd, "region")
	cred := cli.GetCliStringFlag(cmd, "cred")
	nickname := cli.GetCliStringFlag(cmd, "nickname")
	arch := cli.GetCliStringFlag(cmd, "arch")
	archtype := cli.GetCliStringFlag(cmd, "type")
	image := cli.GetCliStringFlag(cmd, "image")
	dhubuser := cli.GetCliStringFlag(cmd, "dhub-user")
	dhubpass := cli.GetCliStringFlag(cmd, "dhub-pass")
	min := cli.GetCliIntFlag(cmd, "min")
	max := cli.GetCliIntFlag(cmd, "max")
	spotrange := cli.GetCliIntFlag(cmd, "spot-range")
	code := cli.GetCliStringFlag(cmd, "code")
	coderef := cli.GetCliStringFlag(cmd, "code-ref")
	codepkey := cli.GetCliStringFlag(cmd, "code-privkey")
	dbengine := cli.GetCliStringFlag(cmd, "dbengine")
	dbtype := cli.GetCliStringFlag(cmd, "dbtype")
	dbstore := cli.GetCliStringFlag(cmd, "dbstorage")
	ecengine := cli.GetCliStringFlag(cmd, "elasticache-engine")
	ectype := cli.GetCliStringFlag(cmd, "elasticache-nodetype")
	eccount := cli.GetCliStringFlag(cmd, "elasticache-nodecount")

	cnf := alm.StackCreateConfig{
		Region:            region,
		Architecture:      arch,
		Type:              archtype,
		Image:             image,
		DockerHubUsername: dhubuser,
		DockerHubPassword: dhubpass,
		Min:               min,
		Max:               max,
		SpotRange:         spotrange,
		Nickname:          nickname,
		Code:              code,
		GitReference:      coderef,
		GitPrivateKey:     codepkey,
	}

	if usedb {
		if dbengine == "" {
			cli.ErrorExit("dbengine is required", 1)
		}

		if dbtype == "" {
			cli.ErrorExit("dbtype is required", 1)
		}

		if dbstore == "" {
			cli.ErrorExit("dbstorage is required", 1)
		}

		dbs := make([]alm.StackCreateDb, 0)
		tmp := alm.StackCreateDb{
			Engine:       dbengine,
			Type:         dbtype,
			Storage:      dbstore,
			ReadReplica1: readreplica1,
			ReadReplica2: readreplica2,
			ReadReplica3: readreplica3,
			ReadReplica4: readreplica4,
			ReadReplica5: readreplica5,
		}

		dbs = append(dbs, tmp)
		cnf.Database = dbs
	}

	if usecache {
		if ecengine == "" {
			cli.ErrorExit("elasticache-engine is required", 1)
		}

		if ectype == "" {
			cli.ErrorExit("elasticache-nodetype is required", 1)
		}

		if eccount == "" {
			cli.ErrorExit("elasticache-nodecount is required", 1)
		}

		caches := make([]alm.StackCreateElasticache, 0)
		tmp := alm.StackCreateElasticache{
			Engine:    ecengine,
			NodeType:  ectype,
			NodeCount: eccount,
		}

		caches = append(caches, tmp)
		cnf.ElastiCache = caches
	}

	// for pretty print
	mi, err := json.MarshalIndent(&cnf, "", pretty.Indent(2))
	cli.ErrorExit(err, 1)

	d.Info("[create stack payload]")
	d.Info("vendor:", vendor)
	d.Info("region:", region)
	d.Info("credentials:", cred)
	d.Info("configurations:")
	fmt.Println(string(mi))

	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackCreateInput{
		Vendor:         vendor,
		Region:         region,
		CredId:         cred,
		Configurations: cnf,
	}

	resp, body, err := svc.Create(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	var success bool
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	cli.ErrorExit(err, 1)

	_, ok := m["stack_id"]
	if ok {
		_, ok = m["status"]
		if ok {
			d.Info(fmt.Sprintf("[%s] stack creation started:", resp.Status))
			d.Info("  stack id:", fmt.Sprintf("%s", m["stack_id"]))
			d.Info("  status:", fmt.Sprintf("%s", m["status"]))
			success = true
		}
	}

	if !success {
		d.Info(string(body))
		return
	}
}

func createAlmStack(cmd *cobra.Command) {
	tf := cli.GetCliStringFlag(cmd, "alm-template")
	b, err := ioutil.ReadFile(tf)
	cli.ErrorExit(err, 1)

	if !filetype.IsJSON(string(b)) {
		cli.ErrorExit("invalid json", 1)
	}

	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackCreateInput{AlmTemplate: &alm.AlmTemplate{
		ContentType: "json",
		Contents:    string(b),
	}}

	resp, body, err := svc.Create(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	if strings.Contains(string(body), "success") {
		res := pretty.JSON(string(body), 2)
		d.Info(fmt.Sprintf("[%s] return payload:", resp.Status))
		fmt.Println(res)
		return
	}

	if (resp.StatusCode / 100) == 2 {
		d.Info(fmt.Sprintf("[%s] %s", resp.Status, string(body)))
	} else {
		d.Error(fmt.Sprintf("[%s] %s", resp.Status, string(body)))
	}
}
