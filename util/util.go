package util

import (
	"log"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = InitLog("util")
	if err != nil {
		log.Fatalf("Error InitLog module[util] err[%s]", err.Error())
	}
}
