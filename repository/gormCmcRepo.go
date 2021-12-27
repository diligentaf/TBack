package repository

import (
	"TBack/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
)

type gormCmcRepo struct {
	conn *gorm.DB
}

func NewGormCmcRepo(conn *gorm.DB, hAruLogGroup string) CmcRepo {
	conn = conn.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&model.Cmc{})
	conn.Model(&model.Cmc{})
	return &gormCmcRepo{conn}
}

// Register ...
func (t gormCmcRepo) Register(ctx context.Context) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("symbol", "BTC,ETH,LTC,XRP,ATOM")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "787030a0-6c5d-40b5-b83e-e690655f71f7")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}

	coinData := new(model.CoinData)
	json.NewDecoder(resp.Body).Decode(&coinData)
	timeNow := time.Now()

	clist := []model.Cmc{
		model.Cmc{
			CmcID:           uuid.New().String(),
			Symbol:          coinData.Data.Btc.Symbol,
			Price:           coinData.Data.Btc.Quote.Usd.Price,
			LastUpdated:     coinData.Data.Btc.LastUpdated,
			CommonTableCols: model.CommonTableCols{CreatedAt: timeNow},
		},
		model.Cmc{
			CmcID:           uuid.New().String(),
			Symbol:          coinData.Data.Eth.Symbol,
			Price:           coinData.Data.Eth.Quote.Usd.Price,
			LastUpdated:     coinData.Data.Eth.LastUpdated,
			CommonTableCols: model.CommonTableCols{CreatedAt: timeNow},
		},
		model.Cmc{
			CmcID:           uuid.New().String(),
			Symbol:          coinData.Data.Atom.Symbol,
			Price:           coinData.Data.Atom.Quote.Usd.Price,
			LastUpdated:     coinData.Data.Atom.LastUpdated,
			CommonTableCols: model.CommonTableCols{CreatedAt: timeNow},
		},
		model.Cmc{
			CmcID:           uuid.New().String(),
			Symbol:          coinData.Data.Ltc.Symbol,
			Price:           coinData.Data.Ltc.Quote.Usd.Price,
			LastUpdated:     coinData.Data.Ltc.LastUpdated,
			CommonTableCols: model.CommonTableCols{CreatedAt: timeNow},
		},
		model.Cmc{
			CmcID:           uuid.New().String(),
			Symbol:          coinData.Data.Xrp.Symbol,
			Price:           coinData.Data.Xrp.Quote.Usd.Price,
			LastUpdated:     coinData.Data.Xrp.LastUpdated,
			CommonTableCols: model.CommonTableCols{CreatedAt: timeNow},
		},
	}

	for i := 0; i < len(clist); i++ {
		result := t.conn.Create(clist[i])
		if err := result.Error; err != nil {
			mlog.Infow("Error Register", "err", err.Error())
			return errors.Annotate(err, "haruCmcRepo Register")
		}
		mlog.Infow("Register Cmc Query Done", "nAffectedRow", result.RowsAffected)
	}
	return nil
}

func (a gormCmcRepo) GetCmcList(ctx context.Context, trID string) ([]*model.Cmc, error) {
	card := []*model.Cmc{}
	offset := 0
	limit := 1000
	orderby := "desc"
	query := a.conn.Offset(offset).Limit(limit).
		Select("*").
		Where("deleted_at IS NULL")

	var od string
	if orderby != "" {
		od = "seq " + orderby
	} else {
		od = "seq asc"
	}

	query = query.Order(od)
	query = query.Find(&card).Debug()
	mlog.Infow("GetCmcList CmcCode Query Done", "trID", trID, "nAffectedRow", query.RowsAffected)

	if err := query.Error; err != nil {
		return nil, errors.Annotate(err, "cardRepo GetCmcList")
	}

	return card, nil
}

func (a gormCmcRepo) GetCmcListBySymbol(ctx context.Context, trID string, symbol *model.GetBySymbolHTTPRequest) ([]*model.Cmc, error) {
	card := []*model.Cmc{}
	offset := 0
	limit := 1000
	orderby := "desc"
	query := a.conn.Offset(offset).Limit(limit).
		Select("*").
		Where("deleted_at IS NULL")

	if symbol.Symbol != "" {
		query = query.Where("symbol = ?", symbol.Symbol)
	}

	if symbol.FromDate != "" && symbol.ToDate != "" {
		query = query.Where("last_updated between ? and ?", symbol.FromDate, symbol.ToDate)
	} else if symbol.FromDate != "" && symbol.ToDate == "" {
		query = query.Where("last_updated >= ?", symbol.FromDate)
	} else if symbol.FromDate == "" && symbol.ToDate != "" {
		query = query.Where("last_updated <= ?", symbol.ToDate)
	}

	var od string
	if orderby != "" {
		od = "seq " + orderby
	} else {
		od = "seq asc"
	}

	query = query.Order(od)
	query = query.Find(&card)
	mlog.Infow("GetCmcList CmcCode Query Done", "trID", trID, "nAffectedRow", query.RowsAffected)

	if err := query.Error; err != nil {
		return nil, errors.Annotate(err, "cardRepo GetCmcList")
	}

	return card, nil
}

func (a gormCmcRepo) GetCmcByOrder(ctx context.Context) (uint, error) {
	var count uint
	result := a.conn.Set("gorm:auto_preload", true).
		Select("*").Count(&count)
	if err := result.Error; err != nil {
		return 0, errors.Annotate(err, "Column GetCmcByOrder")
	}

	if count == 0 {
		err := a.Register(ctx)

		if err != nil {
			mlog.Errorw("Register cmcRepo Register")
			return 0, errors.Annotate(err, "cmcSvc Update")
		}
	}

	mlog.Infow("GetCmcByUID Column Query Done", "nAffectedRow", result.RowsAffected)
	return count, nil
}
