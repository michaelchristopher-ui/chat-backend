package common

import "flag"

// These group of variables variables define the command line flags. Add new command line flags variables here
var (
	ServiceName *string
	NodeName    *string
	IPPort      *string
	CfgPath     *string
)

//SetEnvVars parses Command Line flags for the given names, making them available for use throughout the whole application
func SetEnvVars() {
	ServiceName = flag.String("servicename", "", "service-name")
	NodeName = flag.String("nodename", "", "node-name")
	IPPort = flag.String("ipport", "", "ip-address-and-port-where-this-service-is-deployed")
	CfgPath = flag.String("configpath", "./config-dev.yaml", "path-to-config-file")
	flag.Parse()
}
