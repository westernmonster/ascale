package conf

import (
	"ascale/pkg/cache/redis"
	"ascale/pkg/database/sqalx"
	"ascale/pkg/log"
	"ascale/pkg/net/http/vin"
	"ascale/pkg/tracing"

	"github.com/BurntSushi/toml"
	flag "github.com/spf13/pflag"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	DC     *DC
	Log    *log.Config
	Vin    *vin.ServerConfig
	Tracer *tracing.Config
	DB     *sqalx.Config
	Redis  *redis.Config
}

type DC struct {
	Num  int
	Desc string
}

type SheetItem struct {
	SheetID   string
	SheetName string
}

func init() {
	flag.StringVarP(&confPath, "config", "c", "", "default config path")
	// flag.Parse()
}

// Init init conf
func Init() error {
	return local()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
