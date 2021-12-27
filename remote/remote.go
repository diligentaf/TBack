package remote

import (
	"TBack/util"
	"log"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	var err error
	mlog, err = util.InitLog("rmot")
	if err != nil {
		log.Fatalf("Error InitLog module[rmot] err[%s]", err.Error())
	}
}
