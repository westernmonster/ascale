// Package env get env & app config, all the public field must after init()
// finished and flag.Parse().
package env

import (
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// deploy env.
const (
	DeployEnvDev  = "dev"
	DeployEnvFat1 = "fat1"
	DeployEnvUat  = "uat"
	DeployEnvPre  = "pre"
	DeployEnvProd = "prod"
)

// env default value.
const (
	// env
	_region    = "us-west2"
	_zone      = "us-west2-a"
	_deployEnv = "uat"
	_projectID = "done-280702"
)

// env configuration.
var (
	// ProjectID gcloud projectID
	ProjectID string
	// Site
	SiteURL string

	Domain string

	// Pubsub Endpoint
	PubsubEndpoint string
	// Sematic Version
	Version string
	// Region avaliable region where app at.
	Region string
	// Zone avaliable zone where app at.
	Zone string
	// Hostname machine hostname.
	Hostname string
	// DeployEnv deploy env where app at.
	DeployEnv string

	IP = os.Getenv("POD_IP")
	// AppID is global unique application id, register by service tree.
	// such as main.arch.disocvery.
	AppID string
	// Color is the identification of different experimental group in one caster cluster.
	Color string
)

// app default value.
const (
	_httpPort  = "8000"
	_gorpcPort = "8099"
	_grpcPort  = "9000"
)

// app configraution.
var (
	// HTTPPort app listen http port.
	HTTPPort string
	// GORPCPort app listen gorpc port.
	GORPCPort string
	// GRPCPort app listen grpc port.
	GRPCPort string
)

func GetHostname() (hostname string) {
	var err error
	if hostname, err = os.Hostname(); err != nil || hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}
	hostname = strings.ReplaceAll(hostname, ".", "-")
	return
}

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}

	Hostname = strings.ReplaceAll(Hostname, ".", "-")

	addFlag(flag.CommandLine)

	SiteURL = os.Getenv("SITE_URL")
	Domain = os.Getenv("DOMAIN")
	if SiteURL == "" && Domain != "" {
		SiteURL = "https://" + Domain
	}
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(
		&ProjectID,
		"projectid",
		defaultString("PROJECT_ID", _projectID),
		"gcloud project id",
	)
	fs.StringVar(&Version, "app.version", os.Getenv("APP_VERSION"), "app semantic version")
	fs.StringVar(
		&SiteURL,
		"siteurl",
		os.Getenv("SITE_URL"),
		"site url, example: http://dev.donefirst.com",
	)
	fs.StringVar(&Domain, "domain", os.Getenv("DOMAIN"), "domain, example: dev.donefirst.com")
	fs.StringVar(
		&PubsubEndpoint,
		"pubsub_endpoint",
		defaultString("PUBSUB_ENDPOINT", ""),
		"pubsub endpoint",
	)
	fs.StringVar(
		&Region,
		"region",
		defaultString("REGION", _region),
		"avaliable region. or use REGION env variable, value: us-west1 etc.",
	)
	fs.StringVar(
		&Zone,
		"zone",
		defaultString("ZONE", _zone),
		"avaliable zone. or use ZONE env variable, value: us-west1-a/us-west1-b etc.",
	)
	fs.StringVar(
		&DeployEnv,
		"deploy.env",
		defaultString("DEPLOY_ENV", _deployEnv),
		"deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.",
	)
	fs.StringVar(
		&AppID,
		"appid",
		os.Getenv("APP_ID"),
		"appid is global unique application id, register by service tree. or use APP_ID env variable.",
	)
	fs.StringVar(
		&Color,
		"deploy.color",
		os.Getenv("DEPLOY_COLOR"),
		"deploy.color is the identification of different experimental group.",
	)

	// app
	fs.StringVar(
		&HTTPPort,
		"http.port",
		defaultString("DISCOVERY_HTTP_PORT", _httpPort),
		"app listen http port, default: 8000",
	)
	fs.StringVar(
		&GORPCPort,
		"gorpc.port",
		defaultString("DISCOVERY_GORPC_PORT", _gorpcPort),
		"app listen gorpc port, default: 8100",
	)
	fs.StringVar(
		&GRPCPort,
		"grpc.port",
		defaultString("DISCOVERY_GRPC_PORT", _grpcPort),
		"app listen grpc port, default: 9000",
	)
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}
