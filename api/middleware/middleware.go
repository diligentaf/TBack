package middleware

import (
	"TBack/util"
	"log"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("midw")
	if err != nil {
		log.Fatalf("Error InitLog module[midw] err[%s]", err.Error())
	}
}
