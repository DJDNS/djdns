package app

import "github.com/docopt/docopt-go"

var version = "djdns 0.0.12"
var usage = `Usage: djdns [options]

Options:
    --addr=<address>      host:port to expose DNS on [default: 0.0.0.0:9953]
    --display-name=<name> Hostname to provide in network log messages.
    --root=<root>         Target URL to serve as <ROOT> [default: deje://localhost:8080/root]
    -h --help             Show this message.
    --version             Print the version number.
`

func argToString(arguments map[string]interface{}, key string) string {
	arg := arguments[key]
	if arg == nil {
		return ""
	} else {
		return arg.(string)
	}
}

type ServerConfig struct {
	HostAddress string
	DisplayName string
	RootAlias   string
}

// Parse a list of arguments into a config struct.
//
// To use ARGV flags, pass nil instead of a real []string.
func Parse(argv []string, exit bool) (ServerConfig, error) {
	arguments, err := docopt.Parse(usage, argv, true, version, false, exit)
	if err != nil {
		return ServerConfig{}, err
	}

	var conf ServerConfig
	conf.HostAddress = argToString(arguments, "--addr")
	conf.DisplayName = argToString(arguments, "--display-name")
	conf.RootAlias = argToString(arguments, "--root")

	return conf, nil
}
