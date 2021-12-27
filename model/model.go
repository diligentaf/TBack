package model

import (
	"TBack/util"
	"log"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("modl")
	if err != nil {
		log.Fatalf("Error InitLog module[modl] err[%s]", err.Error())
	}
}

// HTTPResponse ...
type HTTPResponse interface {
	SetResult(resultCode, resultMsg, trID, jobType string)
}
