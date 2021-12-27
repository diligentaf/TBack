package service

import (
	"TBack/model"
	"TBack/repository"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

type cmcUsecase struct {
	cmcRepo    repository.CmcRepo
	ctxTimeout time.Duration
}

// NewCmcSvc ...
func NewCmcSvc(cmcR repository.CmcRepo, timeout time.Duration) CmcSvc {
	return &cmcUsecase{
		cmcRepo:    cmcR,
		ctxTimeout: timeout,
	}
}

// GetConversion ...
func (a cmcUsecase) GetConversion(ctx context.Context, trID string, cmc *model.Conversion) (*model.Conversion, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("symbol", cmc.From)
	q.Add("convert", cmc.To)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "787030a0-6c5d-40b5-b83e-e690655f71f7")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	sbody := string(respBody)
	value := gjson.Get(sbody, "data."+cmc.From+".quote."+cmc.To+".price")
	cmc.Price = value.Num

	return cmc, nil
}

// GetCmcList...
func (a cmcUsecase) GetCmcList(ctx context.Context, trID string) ([]*model.Cmc, error) {
	innerCtx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	existRows, err := a.cmcRepo.GetCmcList(innerCtx, trID)
	if err != nil {
		mlog.Errorw("GetCmcList cmcRepo GetCmcList", "trID", trID, "err", err.Error())
		return nil, errors.Annotate(err, "cmcSvc GetCmcCount")
	}

	return existRows, nil
}

// GetCmcListBySymbol ...
func (a cmcUsecase) GetCmcListBySymbol(ctx context.Context, trID string, symbol *model.GetBySymbolHTTPRequest) ([]*model.Cmc, error) {
	innerCtx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	existRows, err := a.cmcRepo.GetCmcListBySymbol(innerCtx, trID, symbol)
	if err != nil {
		mlog.Errorw("GetCmcList cmcRepo GetCmcList", "trID", trID, "err", err.Error())
		return nil, errors.Annotate(err, "cmcSvc GetCmcCount")
	}

	return existRows, nil
}

// Register ...
func (a cmcUsecase) Register(ctx context.Context) error {
	innerCtx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	err := a.cmcRepo.Register(innerCtx)

	if err != nil {
		mlog.Errorw("Register cmcRepo Register")
		return errors.Annotate(err, "cmcSvc Update")
	}
	return nil
}
