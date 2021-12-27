package service

import (
	"TBack/model"
	"TBack/util"
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("srvc")
	if err != nil {
		log.Fatalf("Error InitLog module[middleware] err[%s]", err.Error())
	}
}

type TBackSvc interface {
	GetLogEvents(ctx context.Context, trID string) ([]*cloudwatchlogs.OutputLogEvent, error)
}

type CmcSvc interface {
	GetCmcList(ctx context.Context, trID string) ([]*model.Cmc, error)
	GetCmcListBySymbol(ctx context.Context, trID string, symbol *model.GetBySymbolHTTPRequest) ([]*model.Cmc, error)
	GetConversion(ctx context.Context, trID string, cmc *model.Conversion) (*model.Conversion, error)
	Register(ctx context.Context) error
}
