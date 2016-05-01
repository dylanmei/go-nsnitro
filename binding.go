package main

import (
	"github.com/dylanmei/go-nsnitro/nsnitro"
	"gopkg.in/alecthomas/kingpin.v2"
)

func doBindLBMonitor(client *nsnitro.Client) {
	err := client.BindServiceGroupToLBMonitor(
		*bind_lb_monitor_servicegroup,
		*bind_lb_monitor_name)

	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doUnbindLBMonitor(client *nsnitro.Client) {
	err := client.UnbindServiceGroupFromLBMonitor(
		*unbind_lb_monitor_servicegroup,
		*unbind_lb_monitor_name)

	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doBindLBVServer(client *nsnitro.Client) {
	err := client.BindLBVServerToServiceGroup(
		*bind_lb_vserver_name,
		*bind_lb_vserver_servicegroup)

	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doUnbindLBVServer(client *nsnitro.Client) {
	err := client.UnbindLBVServerFromServiceGroup(
		*unbind_lb_vserver_name,
		*unbind_lb_vserver_servicegroup)

	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}
