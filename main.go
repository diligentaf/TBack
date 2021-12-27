package main

import (
	ct "TBack/api/controller"
	mw "TBack/api/middleware"
	"TBack/conf"
	repo "TBack/repository"
	"TBack/util"
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

const (
	banner = `
|   |  ----   ___       
|---| |____| |___| |   |
|   | |    | |   \ |___|
                                  

%s
 => Starting listen %s
`
)

var (
	// BuildDate for Program BuildDate
	BuildDate string
	// Version for Program Version
	Version string
	svrInfo = fmt.Sprintf("TBack %s(%s)", Version, BuildDate)
)
var mlog *zap.SugaredLogger

func init() {
	// use all cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	mlog, err = util.InitLog("main")
	if err != nil {
		log.Fatalf("Error InitLog module[main] err[%s]", err.Error())
	}
}

func main() {
	// viper configuration
	tBack := conf.TBack
	if tBack.GetBool("version") {
		mlog.Infow(svrInfo)
		os.Exit(0)
	}
	// sets the url and port
	tBack.SetProfile()

	// setting up the log file
	f := apiLogFile("./tback-api.log")
	defer f.Close()
	e := echoInit(tBack, f)

	// connecting to database
	db := repo.InitDB(tBack)
	defer db.Close()

	// connects to controller and loads up service, repository, etc
	if err := ct.InitHandler(tBack, e, db); err != nil {
		mlog.Errorw("main controller InitHandler", "err", err.Error())
		os.Exit(1)
	}
	initColumn(tBack, db)
	startServer(tBack, e)
	// r := <-longRunningTask(tBack, db)
	// fmt.Println(r)
}

func apiLogFile(logfile string) *os.File {
	// API Logging
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		mlog.Errorw("apiLogFile", "logfile", logfile, "err", err.Error())
		os.Exit(1)
	}
	return f
}

func echoInit(tBack conf.ViperConfig, apiLogFile *os.File) (e *echo.Echo) {
	// Echo instance
	e = echo.New()
	e.Debug = true

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(mw.ResponseDump))
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	// Ping Check
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "TBack API Alive!\n") })
	e.POST("/", func(c echo.Context) error { return c.String(http.StatusOK, "TBack API Alive!\n") })

	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Output = apiLogFile

	e.Use(middleware.LoggerWithConfig(loggerConfig))
	e.Logger.SetOutput(bufio.NewWriterSize(apiLogFile, 1024*16))
	e.Logger.SetLevel(tBack.APILogLevel())
	e.HideBanner = true

	sigInit(e)

	return e
}

func sigInit(e *echo.Echo) chan os.Signal {
	// Signal
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		sig := <-sc
		e.Logger.Error("Got signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Error(err)
		}
		signal.Stop(sc)
		close(sc)
	}()

	return sc
}

func startServer(tBack conf.ViperConfig, e *echo.Echo) {
	// Start Server
	apiServer := fmt.Sprintf("0.0.0.0:%d", tBack.GetInt("port"))
	mlog.Infof("%s => Starting Server Listen %s", svrInfo, apiServer)
	fmt.Printf(banner, svrInfo, apiServer)

	if err := e.Start(apiServer); err != nil {
		mlog.Errorw("startServer", "svrAddr", apiServer, "err", err.Error())
		e.Logger.Error(err)
	}
}

func initColumn(tBack conf.ViperConfig, db *gorm.DB) {
	conn := repo.NewGormCmcRepo(db, tBack.GetString("haru_log_group"))
	_, err := repo.CmcRepo.GetCmcByOrder(conn, nil)

	if gorm.IsRecordNotFoundError(errors.Cause(err)) {
		mlog.Errorw("main GetCmcByOrder", "err", err.Error())
	}

	go func() {
		s := gocron.NewScheduler()
		s.Every(1).Hour().Do(func() { repo.CmcRepo.Register(conn, nil) })
		<-s.Start()
	}()
}
