package controller

import (
	"delta-core/internal"
	"delta-core/repository"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type MarketController struct {
	MarketRepository *repository.RedisMarketRepository
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
func (mc *MarketController) GetQuote(c *gin.Context) {
	fmt.Print("sdf")
	c.JSON(http.StatusOK, "HEllo WORLD")
}

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
	var params = PriceParams{}
	var results []Price
	if c.ShouldBindQuery(&params) == nil {
		log.Println(params)
		log.Println("Error in parsing query parameters")
		c.JSON(http.StatusBadRequest, results)
	}
	tstart, err := time.Parse(time.RFC3339, params.Start)
	if err != nil {
		panic(err)
	}
	tend, err := time.Parse(time.RFC3339, params.End)
	if err != nil {
		panic(err)
	}
	series := mc.MarketRepository.FetchSeries(params.Key, tstart, tend)
	// interval, start, end := params.formatSymbol()
	// quoteParams := &chart.Params{
	// 	Symbol: params.Symbol, //"GBPUSD=X",
	// 	Start:  &start,        //  &datetime.Datetime{
	// 	// 	Day:   25,
	// 	// 	Month: 1,
	// 	// 	Year:  2024,
	// 	// },
	// 	End: &end, // &datetime.Datetime{
	// 	// 	Day:   27,
	// 	// 	Month: 1,
	// 	// 	Year:  2024,
	// 	// },
	// 	Interval: interval, // datetime.OneHour
	// }
	// iter := chart.Get(quoteParams)
	// for iter.Next() {
	// 	point := iter.Bar()
	// 	q := Quotes{
	// 		Open:      point.Open,
	// 		Close:     point.Close,
	// 		High:      point.High,
	// 		Low:       point.Low,
	// 		AdjClose:  point.AdjClose,
	// 		Volume:    point.Volume,
	// 		Timestamp: point.Timestamp,
	// 	}
	// 	// fmt.Println()
	// 	results = append(results, q)
	// }
	c.JSON(http.StatusOK, series)
	// Success!
	// fmt.Println(q)
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
