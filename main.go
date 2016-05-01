package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dylanmei/go-nsnitro/nsnitro"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("nsnitro", "A NetScaler 10+ Nitro API cli")
	ns_server   = app.Flag("server", "URL of the NetScalar server").Envar("NSNITRO_SERVER").Required().URL()
	ns_username = app.Flag("username", "NetScaler Nitro API user name").Envar("NSNITRO_USERNAME").String()
	ns_password = app.Flag("password", "NetScaler Nitro API password").Envar("NSNITRO_PASSWORD").String()

	show                   = app.Command("show", "")
	show_server            = show.Command("server", "Print one or all servers")
	show_server_name       = show_server.Arg("name", "").String()
	show_lb                = show.Command("lb", "")
	show_lb_vserver        = show_lb.Command("vserver", "Print one or all lb vservers")
	show_lb_vserver_name   = show_lb_vserver.Arg("name", "").String()
	show_lb_monitor        = show_lb.Command("monitor", "Print one or all lb monitors")
	show_lb_monitor_name   = show_lb_monitor.Arg("name", "").String()
	show_servicegroup      = show.Command("servicegroup", "Print one or all servicegroups")
	show_servicegroup_name = show_servicegroup.Arg("name", "").String()
	show_version           = show.Command("version", "Print the NetScalar version")

	//	add                    = app.Command("add", "")
	//	add_server             = add.Command("server", "")
	//	add_server_name        = add_server.Arg("name", "").String()
	//	add_server_ipv4        = add_server.Arg("ip", "").IP()
	//	add_lb                 = add.Command("lb", "")
	//	add_lb_vserver         = add_lb.Command("vserver", "")
	//	add_lb_vserver_name    = add_lb_vserver.Arg("name", "").String()
	//	add_lb_vserver_type    = add_lb_vserver.Arg("service-type", "").String()
	//	add_lb_vserver_ipv4    = add_lb_vserver.Arg("ip", "").IP()
	//	add_lb_vserver_port    = add_lb_vserver.Arg("port", "").Int()
	//	add_lb_monitor         = add_lb.Command("monitor", "")
	//	add_lb_monitor_name    = add_lb_monitor.Arg("name", "").String()
	//	add_lb_monitor_type    = add_lb_monitor.Arg("type", "").String()
	//	add_lb_monitor_send    = add_lb_monitor.Flag("send", "").String()
	//	add_lb_monitor_recv    = add_lb_monitor.Flag("recv", "").String()
	//	add_lb_monitor_port    = add_lb_monitor.Flag("dest-port", "").Int()
	//	add_lb_monitor_headers = add_lb_monitor.Flag("custom-headers", "").String()
	//	add_servicegroup       = add.Command("servicegroup", "")
	//	add_servicegroup_name  = add_servicegroup.Arg("name", "").String()
	//	add_servicegroup_type  = add_servicegroup.Arg("service-type", "").String()

	bind                         = app.Command("bind", "")
	bind_lb                      = bind.Command("lb", "")
	bind_lb_monitor              = bind_lb.Command("monitor", "Bind an lb monitor to a service or service-group")
	bind_lb_monitor_name         = bind_lb_monitor.Arg("name", "Name of an lb monitor").Required().String()
	bind_lb_monitor_servicegroup = bind_lb_monitor.Flag("service-group", "Name of a service-group").String()
	bind_lb_vserver              = bind_lb.Command("vserver", "Bind an lb vserver to a service or service-group")
	bind_lb_vserver_name         = bind_lb_vserver.Arg("name", "Name of an lb vserver").Required().String()
	bind_lb_vserver_servicegroup = bind_lb_vserver.Flag("service-group", "Name of a service-group").String()
	bind_servicegroup            = bind.Command("servicegroup", "Bind a service-group to a service")
	bind_servicegroup_name       = bind_servicegroup.Arg("name", "Name of an service-group").Required().String()
	bind_servicegroup_server     = bind_servicegroup.Arg("server", "Name of a server").Required().String()
	bind_servicegroup_port       = bind_servicegroup.Arg("port", "Port of a server").Required().Int()

	unbind                         = app.Command("unbind", "")
	unbind_lb                      = unbind.Command("lb", "")
	unbind_lb_monitor              = unbind_lb.Command("monitor", "Unbind an lb monitor from a service or service-group")
	unbind_lb_monitor_name         = unbind_lb_monitor.Arg("name", "Name of an lb monitor").Required().String()
	unbind_lb_monitor_servicegroup = unbind_lb_monitor.Flag("service-group", "Name of a service-group").String()
	unbind_lb_vserver              = unbind_lb.Command("vserver", "Unbind an lb vserver from a service or service-group")
	unbind_lb_vserver_name         = unbind_lb_vserver.Arg("name", "Name of an lb vserver").Required().String()
	unbind_lb_vserver_servicegroup = unbind_lb_vserver.Flag("service-group", "Name of a service-group").String()
	unbind_servicegroup            = unbind.Command("servicegroup", "Unbind a service-group to a service")
	unbind_servicegroup_name       = unbind_servicegroup.Arg("name", "Name of an service-group").Required().String()
	unbind_servicegroup_server     = unbind_servicegroup.Arg("server", "Name of a server").Required().String()
	unbind_servicegroup_port       = unbind_servicegroup.Arg("port", "Port of a server").Required().Int()
)

func main() {
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		kingpin.Fatalf("%s, try %s help", err, os.Args[:1])
	}

	client, err := newClient(*ns_server, *ns_username, *ns_password)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}

	switch cmd {
	case "show server":
		doShowServer(client)
		break
	case "show servicegroup":
		doShowServiceGroup(client)
		break
	case "show lb monitor":
		doShowLBMonitor(client)
		break
	case "show lb vserver":
		doShowLBVServer(client)
		break
	case "show version":
		doShowVersion(client)
		break
	case "bind lb monitor":
		doBindLBMonitor(client)
		break
	case "unbind lb monitor":
		doUnbindLBMonitor(client)
		break
	case "bind lb vserver":
		doBindLBVServer(client)
		break
	case "unbind lb vserver":
		doUnbindLBVServer(client)
		break
	case "bind servicegroup":
		doBindServiceGroup(client)
		break
	case "unbind servicegroup":
		doUnbindServiceGroup(client)
		break
	}
}

func newClient(uri *url.URL, username, password string) (*nsnitro.Client, error) {
	if uri.User != nil {
		if username == "" {
			username = uri.User.Username()
		}
		if password == "" {
			password, _ = uri.User.Password()
		}

		uri.User = nil
	}

	config := &nsnitro.Config{
		URL:      uri.String(),
		User:     username,
		Password: password,
	}

	config.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client := nsnitro.NewClient(config)
	_, err := client.Version()
	return client, err
}
