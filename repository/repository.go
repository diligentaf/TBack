package repository

import (
	"TBack/conf"
	"TBack/model"
	"TBack/util"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("repo")
	if err != nil {
		log.Fatalf("Error InitLog module[repo] err[%s]", err.Error())
	}
}

// InitDB ...
func InitDB(tBack conf.ViperConfig) *gorm.DB {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		tBack.GetString("db_user"),
		tBack.GetString("db_pass"),
		tBack.GetString("db_host"),
		tBack.GetInt("db_port"),
		tBack.GetString("db_name"),
	)

	mlog.Infow("InitDB DB Connect", "dbHost", tBack.GetString("db_host"), "dbPort", tBack.GetString("db_port"), "dbName", tBack.GetString("db_name"))
	dbConn, err := gorm.Open("mysql", dbURI)
	if err != nil {
		mlog.Errorw("InitDB gorm.Open", "dbHost", tBack.GetString("db_host"), "dbPort", tBack.GetString("db_port"), "dbName", tBack.GetString("db_name"), "err", err.Error())
		os.Exit(1)
	}

	dbConn.DB().SetMaxIdleConns(100)
	dbConn.DB().SetConnMaxLifetime(time.Second)
	// MOOG
	// dbConn.LogMode(true)
	mlog.Infow("InitDB DB Connect Done")

	return dbConn
}

// TBackRepo ...
type TBackRepo interface {
	GetLogEvents(ctx context.Context, trID string) ([]*cloudwatchlogs.OutputLogEvent, error)
}

// CmcRepo ...
type CmcRepo interface {
	Register(ctx context.Context) error
	GetCmcListBySymbol(ctx context.Context, trID string, symbol *model.GetBySymbolHTTPRequest) ([]*model.Cmc, error)
	GetCmcList(ctx context.Context, trID string) ([]*model.Cmc, error)
	GetCmcByOrder(ctx context.Context) (uint, error)
}

// innerUnitRepo ...
type innerUnitRepo interface {
	GetLogEvents(trID string) ([]*cloudwatchlogs.OutputLogEvent, error)
}
