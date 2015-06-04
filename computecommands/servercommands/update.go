package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var update = cli.Command{
	Name:        "update",
	Usage:       fmt.Sprintf("%s %s update <serverID> [flags]", util.Name, commandPrefix),
	Description: "Updates an existing server",
	Action:      commandUpdate,
	Flags:       util.CommandFlags(flagsUpdate),
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Update the server's name",
		},
		cli.StringFlag{
			Name:  "accessIPv4",
			Usage: "Update the server's IPv4 address",
		},
		cli.StringFlag{
			Name:  "accessIPv6",
			Usage: "Update the server's IPv6 address",
		},
	}
}

func commandUpdate(c *cli.Context) {
	util.CheckArgNum(c, 1)
	serverID := c.Args()[0]
	opts := &osServers.UpdateOpts{
		Name:       c.String("name"),
		AccessIPv4: c.String("accessIPv4"),
		AccessIPv6: c.String("accessIPv6"),
	}
	client := auth.NewClient("compute")
	o, err := servers.Update(client, serverID, opts).Extract()
	if err != nil {
		fmt.Printf("Error updating server: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableUpdate)
}

func tableUpdate(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"ID", "Name", "Public IPv4", "Public IPv6"}
	for _, key := range keys {
		tmp := ""
		switch key {
		case "Public IPv4":
			tmp = fmt.Sprint(m["AccessIPv4"])
		case "Public IPv6":
			tmp = fmt.Sprint(m["AccessIPv6"])
		default:
			tmp = fmt.Sprint(m[key])
		}
		if tmp == "<nil>" {
			tmp = ""
		}
		t.Append([]string{key, tmp})
	}
	t.Render()
}
