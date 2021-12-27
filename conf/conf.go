package conf

import (
	"TBack/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	elog "github.com/labstack/gommon/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DefaultConf ...
type DefaultConf struct {
	envDEV   string
	envSTAGE string
	envPROD  string

	port          int
	profile       bool
	profilePort   int
	serverTimeout int
	apiLogLevel   string

	dbHost string
	dbPort int
	dbUser string
	dbPass string
	dbName string

	tBackLogGroup string
	hAruLogGroup  string

	remoteUser string
	remotePass string
}

var defaultConf = DefaultConf{
	envDEV:   ".env.dev",
	envSTAGE: ".env.stage",
	envPROD:  ".env",

	port:          2379,
	profile:       true,
	profilePort:   6071,
	serverTimeout: 30,
	apiLogLevel:   "debug",
}

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
}

// TBack ...
var TBack ViperConfig
var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("conf")
	if err != nil {
		log.Fatalf("Error InitLog module[conf] err[%s]", err.Error())
	}

	pflag.BoolP("version", "v", false, "Show version number and quit")
	pflag.IntP("port", "p", defaultConf.port, "TBack Port")
	pflag.IntP("timeout", "t", defaultConf.serverTimeout, "TBack context timeout (sec)")

	pflag.String("db_host", defaultConf.dbHost, "TBack's DB host")
	pflag.Int("db_port", defaultConf.dbPort, "TBack's DB port")
	pflag.String("db_user", defaultConf.dbUser, "TBack's DB user")
	pflag.String("db_pass", defaultConf.dbPass, "TBack's DB password")
	pflag.String("db_name", defaultConf.dbName, "TBack's DB name")
	pflag.String("tback_log_group", defaultConf.tBackLogGroup, "TBack's logGroup name")
	pflag.String("haru", defaultConf.hAruLogGroup, "HAru's logGroup name")

	pflag.String("remote_user", defaultConf.remoteUser, "Remote user")
	pflag.String("remote_pass", defaultConf.remotePass, "Remote password")

	pflag.Parse()

	TBack, err = readConfig(map[string]interface{}{
		"port":            defaultConf.port,
		"timeout":         defaultConf.serverTimeout,
		"api_log_level":   defaultConf.apiLogLevel,
		"profile":         defaultConf.profile,
		"profile_port":    defaultConf.profilePort,
		"tback_log_group": defaultConf.tBackLogGroup,
		"haru_log_group":  defaultConf.hAruLogGroup,
	})
	if err != nil {
		mlog.Errorw("init readConfig", "err", err.Error())
		os.Exit(1)
	}

	TBack.BindPFlags(pflag.CommandLine)

	TBack.validation()
}

func readConfig(defaults map[string]interface{}) (ViperConfig, error) {
	// Read Sequence (will overloading)
	// defaults -> config file -> env -> cmd flag
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath("./conf")
	v.AutomaticEnv()

	env := strings.ToUpper(v.GetString("env"))
	switch env {
	case "DEVELOPMENT":
		mlog.Infow("Loading Enviroment", "env", env)
		v.SetConfigName(defaultConf.envDEV)
	case "STAGE":
		mlog.Infow("Loading Enviroment", "env", env)
		v.SetConfigName(defaultConf.envSTAGE)
	case "PRODUCTION":
		mlog.Infow("Loading Enviroment", "env", env)
		v.SetConfigName(defaultConf.envPROD)
	default:
		mlog.Infow("Loading Enviroment", "env", "PRODUCTION(Default)")
		v.SetConfigName(defaultConf.envPROD)
	}

	if err := v.ReadInConfig(); err != nil {
		return ViperConfig{}, err
	}

	return ViperConfig{v}, nil
}

// APILogLevel string to log level
func (vp ViperConfig) APILogLevel() elog.Lvl {
	switch strings.ToLower(vp.GetString("api_log_level")) {
	case "off":
		return elog.OFF
	case "error":
		return elog.ERROR
	case "warn", "warning":
		return elog.WARN
	case "info":
		return elog.INFO
	case "debug":
		return elog.DEBUG
	default:
		return elog.DEBUG
	}
}

// SetProfile ...
func (vp ViperConfig) SetProfile() {
	if vp.GetBool("profile") {
		runtime.SetBlockProfileRate(1)
		go func() {
			profileListen := fmt.Sprintf("0.0.0.0:%d", vp.GetInt("profile_port"))
			http.ListenAndServe(profileListen, nil)
		}()
	}
}

func (vp ViperConfig) validation() {
	if vp.GetInt("port") < 1024 || vp.GetInt("port") > 49151 {
		mlog.Errorw("validation Config Validation", "port", vp.GetInt("port"), "err", "Port Out of Range (range: 1024 ~ 49151)")
		os.Exit(1)
	}

	if vp.GetBool("profile") && (vp.GetInt("profile_port") < 1024 || vp.GetInt("profile_port") > 49151) {
		mlog.Errorw("validation Config Validation", "profilePort", vp.GetInt("profile_port"), "err", "Profile Port Out of Range (range: 1024 ~ 49151)")
		os.Exit(1)
	}

	if vp.GetInt("timeout") < 1 {
		mlog.Errorw("validation Config Validation", "timeout", vp.GetInt("timeout"), "err", "Invalid Timeout")
		os.Exit(1)
	}

	if vp.GetString("db_host") == "" {
		mlog.Errorw("validation Config Validation", "err", "Empty DB Host")
		os.Exit(1)
	}

	if vp.GetString("db_port") == "" {
		mlog.Errorw("validation Config Validation", "err", "Empty DB Port")
		os.Exit(1)
	}

	if vp.GetString("db_user") == "" {
		mlog.Errorw("validation Config Validation", "err", "Empty DB User")
		os.Exit(1)
	}

	if vp.GetString("db_pass") == "" {
		mlog.Errorw("validation Config Validation", "err", "Empty DB Password")
		os.Exit(1)
	}

	if vp.GetString("db_name") == "" {
		mlog.Errorw("validation Config Validation", "err", "Empty DB Name")
		os.Exit(1)
	}

	mlog.Infow("validation Config Validataion Done")
}
