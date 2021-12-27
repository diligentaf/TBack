package controller

import (
	mw "TBack/api/middleware"
	"TBack/api/service"
	"TBack/conf"
	"TBack/remote"
	repo "TBack/repository"
	"TBack/util"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("ctlr")
	if err != nil {
		log.Fatalf("Error InitLog module[ctlr] err[%s]", err.Error())
	}
}

// InitHandler ...
func InitHandler(tBack conf.ViperConfig, e *echo.Echo, db *gorm.DB) (err error) {
	timeout := time.Duration(tBack.GetInt("timeout")) * time.Second

	//
	cmcRepo := repo.NewGormCmcRepo(db, tBack.GetString("haru_card_log_group"))

	//
	cmcSvc := service.NewCmcSvc(cmcRepo, timeout)

	// Default Group
	api := e.Group("/api")
	ver := api.Group("/v1")
	ver.Use(mw.TransID())

	sys := ver.Group("/tback")
	if tBack.GetBool("jwt_mode") {
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(strings.Replace(tBack.GetString("jwt_public_key"), "\\n", "\n", -1)))
		if err != nil {
			return errors.Annotatef(err, "controller InitHandler")
		}

		sys.Use(middleware.JWTWithConfig(
			middleware.JWTConfig{
				SigningMethod: "RS256",
				SigningKey:    publicKey,
				ContextKey:    "token",
			},
		))
	}
	cmc := sys.Group("/cmc")

	/*
	 * IMPORTANT
	 * Private APIs can be accessed by public because these APIs are open to public.
	 * Therefore, validation check is strict
	 */
	privateSys := ver.Group("/private/TBack")
	privateSys.Use(middleware.BasicAuth(remote.ValidateRemote))

	newCmcHTTPHandler(cmc, cmcSvc, tBack.GetBool("jwt_mode"))

	mlog.Infow("InitHandler HTTP Handler Initailize Done")

	return nil
}

type cmcHTTPHandler struct {
	cmcSvc  service.CmcSvc
	jwtMode bool
}

func newCmcHTTPHandler(eg *echo.Group, scd service.CmcSvc, jm bool) {
	handler := cmcHTTPHandler{
		cmcSvc:  scd,
		jwtMode: jm,
	}

	eg.GET("/symbol", handler.getCmcListBySymbol)
	eg.GET("/conversion", handler.getConversion)
}

type tbackHTTPHandler struct {
	tbackSvc service.TBackSvc
	jwtMode  bool
}
