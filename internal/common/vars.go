package common

import "flag"

var (
	ServiceName *string
	NodeName    *string
	IPPort      *string
	CfgPath     *string
)

func SetEnvVars() {
	ServiceName = flag.String("servicename", "", "service-name")
	NodeName = flag.String("nodename", "", "node-name")
	IPPort = flag.String("ipport", "", "ip-address-and-port-where-this-service-is-deployed")
	CfgPath = flag.String("configpath", "./config.yaml", "path-to-config-file")
	flag.Parse()
}
