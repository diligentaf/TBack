package controller

import (
	"TBack/model"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/juju/errors"
	"github.com/labstack/echo"
)

func (a cmcHTTPHandler) getConversion(c echo.Context) error {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cmc := new(model.Conversion)

	if err := c.Bind(cmc); err != nil {
		mlog.Warnw("cardCtrl getConversion Bind error", "trID", trID, "err", err.Error())
		return c.JSON(http.StatusBadRequest, model.CommonHTTPResponse{
			ResultCode: "1001",
			ResultMsg:  fmt.Sprintf("Invalid Parameter err[%s]", err.Error()),
			TrID:       trID,
		})
	}

	if cmc.From == "" || cmc.To == "" {
		mlog.Infow("cardCtrl columnHTTPHandler getConversion (textfield missing)", "trID", trID, "From", cmc.From, "To", cmc.To)
		return c.JSON(http.StatusUnauthorized, model.CommonHTTPResponse{
			ResultCode: "1010",
			ResultMsg:  "symbol missing",
			TrID:       trID,
		})
	}

	result, err := a.cmcSvc.GetConversion(ctx, trID, cmc)
	if err != nil {
		mlog.Errorw("getConversion cmcSvc Create", "trID", trID, "From", cmc.From, "err", err.Error())
		if errors.Cause(err).Error() == "Invalid cmc ID" {
			mlog.Warnw("getConversion Parameter Validation", "trID", trID)
			return c.JSON(http.StatusBadRequest, model.CommonHTTPResponse{
				ResultCode: "1001",
				ResultMsg:  fmt.Sprintf("Invalid parameter[%s]", err.Error()),
				TrID:       trID,
			})
		}

		if strings.Contains(errors.Cause(err).Error(), "Duplicate entry") {
			mlog.Warnw("getConversion Duplicate entry", "trID", trID)
			return c.JSON(http.StatusBadRequest, model.CommonHTTPResponse{
				ResultCode: "1002",
				ResultMsg:  "Duplicate entry",
				TrID:       trID,
			})
		}
	}

	return c.JSON(http.StatusOK, model.ConversionResponse{
		CommonHTTPResponse: model.CommonHTTPResponse{
			ResultCode: "0000",
			ResultMsg:  "Conversion Done",
			TrID:       trID,
		},
		Conversion: result,
	})
}

func (a cmcHTTPHandler) getCmcListBySymbol(c echo.Context) error {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	symbol := new(model.GetBySymbolHTTPRequest)

	if err := c.Bind(symbol); err != nil {
		mlog.Warnw("cardCtrl getCmcListBySymbol Bind error", "trID", trID, "err", err.Error())
		return c.JSON(http.StatusBadRequest, model.CommonHTTPResponse{
			ResultCode: "1001",
			ResultMsg:  fmt.Sprintf("Invalid Parameter err[%s]", err.Error()),
			TrID:       trID,
		})
	}

	result, err := a.cmcSvc.GetCmcListBySymbol(ctx, trID, symbol)
	if err != nil {
		mlog.Errorw("cmcCtrl getCmcListBySymbol cmcSvc GetCmcList", "trID", trID, "err", err.Error())
		return c.JSON(http.StatusInternalServerError, model.CommonHTTPResponse{
			ResultCode: "1101",
			ResultMsg:  "Internal Server Error",
			TrID:       trID,
		})
	}

	return c.JSON(http.StatusOK, model.CmcListResponse{
		CommonHTTPResponse: model.CommonHTTPResponse{
			ResultCode: "0000",
			ResultMsg:  "Get CmcListBySymbol Fetched",
			TrID:       trID,
		},
		CmcList: result,
	})
}
