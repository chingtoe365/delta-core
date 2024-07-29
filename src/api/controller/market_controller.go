package controller

import (
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal"
	"delta-core/repository"
	"delta-core/usecase"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketController struct {
	MarketRepository   *repository.MarketRepository
	SignalSetupUsecase *usecase.SignalSetupUsecase
	ProfileUsecase     domain.ProfileUsecase
	Env                *bootstrap.Env
}

type QuoteParams struct {
	Symbol   string `form:"symbol" binding:"required"`
	Interval int    `form:"interval" default:"60"`
	Start    string `form:"start" default:"2023-01-02T15:04:05Z"`
	End      string `form:"end" default:"2023-01-03T15:04:05Z"`
}
type PriceParams struct {
	Key   string `form:"key" binding:"required"`
	Start string `form:"start" default:"2024-07-21T20:00:00Z"`
	End   string `form:"end" default:"2024-07-21T20:00:30Z"`
}
type Quotes struct {
	Open      decimal.Decimal `json:"open"`
	Low       decimal.Decimal `json:"low"`
	High      decimal.Decimal `json:"high"`
	Close     decimal.Decimal `json:"close"`
	AdjClose  decimal.Decimal `json:"adjClose"`
	Volume    int             `json:"volume"`
	Timestamp int             `json:"timestamp"`
}
type Price struct {
	Key   string
	Time  string
	Value decimal.Decimal
}

// // GetQuote godoc
// // @Summary get quote
// // @Schemes
// // @Description get quote
// // @Tags Market
// // @Param symbol query string true "Symbol"
// // @Param interval query string false "Interval"
// // @Param start query string false "Start time"
// // @Param end query string false "End time"
// // @Accept json
// // @Produce json
// // @Success 200
// // @Router /get-quote [get]
// func (mc *MarketController) GetQuote(c *gin.Context) {
// 	fmt.Print("sdf")
// 	c.JSON(http.StatusOK, "HEllo WORLD")
// }

// GetQuote godoc
// @Summary get series
// @Schemes
// @Description get series
// @Tags Market
// @Param key query string true "Key symbols"
// @Param start query string false "Start time"
// @Param end query string false "End time"
// @Accept json
// @Produce json
// @Success 200
// @Router /get-series [get]
func (mc *MarketController) GetSeries(c *gin.Context) {
	var params = PriceParams{
		Key:   c.Query("key"),
		Start: c.Query("start"),
		End:   c.Query("end"),
	}
	tstart, err := time.Parse(time.RFC3339, params.Start)
	if err != nil {
		panic(err)
	}
	tend, err := time.Parse(time.RFC3339, params.End)
	if err != nil {
		panic(err)
	}
	series := mc.MarketRepository.FetchSeries(c, params.Key, tstart, tend)
	log.Println(series)
	c.JSON(http.StatusOK, series)
}

// PingExample godoc
// @Summary List singal reminder
// @Schemes
// @Description List customized market signals
// @Tags Market
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200
// @Router /list-signals [get]
func (mc *MarketController) ListSubscribedTradeSignals(c *gin.Context) {
	userID := c.GetString("x-user-id")
	signals, err := mc.MarketRepository.FetchByUserID(c, userID)
	// tasks, err := u.TaskUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	// tasksOut := internal.FilterTasks(tradeItem, tradeItemCategory, tradeSignal, tasks)
	c.JSON(http.StatusOK, signals)
}

// PingExample godoc
// @Summary Setup singal reminder
// @Schemes
// @Description Customize market signals for email/msg reminder
// @Tags Market
// @Security ApiKeyAuth
// @Param config body domain.MarketSignalMeta true "Signal setup request".
// @Accept json
// @Produce json
// @Success 200
// @Router /setup-signal [post]
func (mc *MarketController) SetupTradeSignals(c *gin.Context) {
	userID := c.GetString("x-user-id")
	var request domain.MarketSignalMeta
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	}
	userProfile, err := mc.ProfileUsecase.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	signalId := primitive.NewObjectID()
	signaler, err := mc.SignalSetupUsecase.MakeMarketSignaler(
		signalId, request.Key, string(request.Type), request.Config, mc.MarketRepository, c, mc.Env, userProfile)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	go signaler.Roll(signalId)
	err = mc.MarketRepository.CreateSignaler(c, domain.MarketSignalDto{
		Id:     signalId,
		UserId: userId,
		SignalMeta: domain.MarketSignalMeta{
			Key:    request.Key,
			Type:   request.Type,
			Config: request.Config,
		},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	}
	// request.
	c.JSON(http.StatusAccepted, "Trade signal was setup")
}

// PingExample godoc
// @Summary Delete singal reminder
// @Schemes
// @Description Remove customized market signals
// @Tags Market
// @Security ApiKeyAuth
// @Param signalId query string true "ID of signal to be deleted".
// @Accept json
// @Produce json
// @Success 200
// @Router /delete-signal [delete]
func (mc *MarketController) DeleteTradeSignals(c *gin.Context) {
	signalId := c.Query("signalId")
	marketSignal, err := mc.MarketRepository.FetchSignalerById(c, signalId)
	log.Print(marketSignal)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	}
	mc.SignalSetupUsecase.RemoveMarketSignaler(signalId)
	err = mc.MarketRepository.Delete(c, &marketSignal)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
	}
	// request.
	c.JSON(http.StatusAccepted, "Trade signal was deleted")
}

// PingExample godoc
// @Summary Fetch all available trade signal categories
// @Schemes
// @Description Fetch trade signal categories
// @Tags Market
// @Accept json
// @Produce json
// @Success 200
// @Router /get-all-trade-signal-categories [get]
func (u *MarketController) GetAllTradeSignalCategories(c *gin.Context) {
	// userID := c.GetString("x-user-id")
	tradeSignalCategories := internal.ReadTradeSignalCategoriesFromJsonFile()
	c.JSON(http.StatusOK, tradeSignalCategories)
}

// PingExample godoc
// @Summary Fetch all available trade items
// @Schemes
// @Description Fetch trade items
// @Tags Market
// @Accept json
// @Produce json
// @Success 200
// @Router /get-all-trade-items [get]
func (u *MarketController) GetAllTradeItems(c *gin.Context) {
	// userID := c.GetString("x-user-id")
	tradeItems := internal.ReadTradeItemsFromJsonFile()
	c.JSON(http.StatusOK, tradeItems)
}

// PingExample godoc
// @Summary Fetch all available trade signals
// @Schemes
// @Description Fetch trade signals
// @Tags Market
// @Accept json
// @Produce json
// @Success 200
// @Router /get-all-trade-signals [get]
func (u *MarketController) GetAllTradeSignals(c *gin.Context) {
	// userID := c.GetString("x-user-id")
	tradeSignals := internal.ReadTradeSignalsFromJsonFile()
	c.JSON(http.StatusOK, tradeSignals)
}
