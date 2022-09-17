package main

import (
	"flag"
	"strings"

	"github.com/maciejkosiarski/lucy/forwarder"
)

func main() {

	var arg_c string
	flag.StringVar(&arg_c, "c", "", "Cluster names (coma separated). Default all.")
	flag.Parse()
	cluster_names := ParseClusterNames(arg_c)

	f := forwarder.Forwarder{
		ConfigFileName: "forwarder.yaml",
		ClusterNames:   cluster_names,
	}
	f.ForwardPorts(cluster_names)
}

func ParseClusterNames(value string) []string {
	s := strings.Split(value, ",")
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	return s
}
