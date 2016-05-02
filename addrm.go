package main

import (
	"github.com/dylanmei/go-nsnitro/nsnitro"
	"gopkg.in/alecthomas/kingpin.v2"
)

func doAddLBMonitor(client *nsnitro.Client) {
	lbmonitor := nsnitro.LBMonitor{
		Name: *add_lb_monitor_name,
		Type: *add_lb_monitor_type,
		Send: *add_lb_monitor_send,
		Recv: *add_lb_monitor_recv,
	}
	if add_lb_monitor_port != nil {
		lbmonitor.Port = *add_lb_monitor_port
	}
	if add_lb_monitor_interval != nil {
		lbmonitor.Interval = *add_lb_monitor_interval
	}

	err := client.AddLBMonitor(lbmonitor)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doRemoveLBMonitor(client *nsnitro.Client) {
	err := client.RemoveLBMonitor(
		*rm_lb_monitor_name,
		*rm_lb_monitor_type)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doAddLBVServer(client *nsnitro.Client) {
	lbvserver := nsnitro.LBVServer{
		Name: *add_lb_vserver_name,
		Type: *add_lb_vserver_type,
		IP:   add_lb_vserver_ipv4.String(),
		Port: *add_lb_vserver_port,
	}

	err := client.AddLBVServer(lbvserver)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doRemoveLBVServer(client *nsnitro.Client) {
	err := client.RemoveLBVServer(*rm_lb_vserver_name)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doAddServer(client *nsnitro.Client) {
	server := nsnitro.Server{
		Name: *add_server_name,
		IP:   add_server_ipv4.String(),
	}

	err := client.AddServer(server)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doRemoveServer(client *nsnitro.Client) {
	err := client.RemoveServer(*rm_server_name)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doAddServiceGroup(client *nsnitro.Client) {
	servicegroup := nsnitro.ServiceGroup{
		Name: *add_servicegroup_name,
		Type: *add_servicegroup_type,
	}

	err := client.AddServiceGroup(servicegroup)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}

func doRemoveServiceGroup(client *nsnitro.Client) {
	err := client.RemoveServiceGroup(*rm_servicegroup_name)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}
}
