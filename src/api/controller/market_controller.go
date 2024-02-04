package controller

import (
	"delta-core/internal"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/shopspring/decimal"
)

type MarketController struct {
}

type QuoteParams struct {
	Symbol   string `form:"symbol" binding:"required"`
	Interval int    `form:"interval" default:"60"`
	Start    string `form:"start" default:"2023-01-02T15:04:05Z"`
	End      string `form:"end" default:"2023-01-03T15:04:05Z"`
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

func (qp *QuoteParams) formatSymbol() (datetime.Interval, datetime.Datetime, datetime.Datetime) {
	tradeItems := internal.ReadTradeItemsFromJsonFile()
	// treatment for symbol
	for _, item := range tradeItems.Items {
		if item.Name == qp.Symbol {
			if item.Category == "Forex" {
				qp.Symbol = qp.Symbol + "=X"
				break
			}
		}
	}
	var interval datetime.Interval
	// treatment for interval
	switch i := qp.Interval; i {
	case 60:
		interval = datetime.OneHour
	case 1:
		interval = datetime.OneMin
	case 1440:
		interval = datetime.OneDay
	case 5:
		interval = datetime.FiveMins
	case 15:
		interval = datetime.FifteenMins
	}
	var start datetime.Datetime
	var end datetime.Datetime
	// dateForm := "2022/05/12"
	tstart, _ := time.Parse(time.RFC3339, qp.Start)
	// if err == nil {
	start = datetime.Datetime{
		Day: tstart.Day(), Month: int(tstart.Month()), Year: tstart.Year(),
	}
	// } else {
	// 	start = datetime.Datetime{
	// 		Day: time.Now().Day(), Month: int(time.Now().Month()), Year: time.Now().Year(),
	// 	}
	// }
	tend, _ := time.Parse(time.RFC3339, qp.End)
	// if err == nil {
	end = datetime.Datetime{
		Day: tend.Day(), Month: int(tend.Month()), Year: tend.Year(),
	}
	// } else {
	// 	end = datetime.Datetime{
	// 		Day: time.Now().Day(), Month: int(time.Now().Month()), Year: time.Now().Year(),
	// 	}
	// }
	// end := time.Parse("2022/05/12", qp.End)
	return interval, start, end
}

// GetQuote godoc
// @Summary get quote
// @Schemes
// @Description get quote
// @Tags Market
// @Param symbol query string true "Symbol"
// @Param interval query string false "Interval"
// @Param start query string false "Start time"
// @Param end query string false "End time"
// @Accept json
// @Produce json
// @Success 200
// @Router /get-quote [get]
func (mc *MarketController) GetQuote(c *gin.Context) {
	// q, err := quote.Get("AAPL")
	// if err != nil {
	// 	// Uh-oh.
	// 	panic(err)
	// }
	var params = QuoteParams{}
	if c.ShouldBindQuery(&params) == nil {
		log.Println(params.Symbol)
		log.Println(params.Interval)
		log.Println(params.Start)
		log.Println(params.End)
	}
	interval, start, end := params.formatSymbol()
	log.Println(interval)
	log.Println(start)
	log.Println(end)
	log.Println(params.Symbol)
	quoteParams := &chart.Params{
		Symbol: params.Symbol, //"GBPUSD=X",
		Start:  &start,        //  &datetime.Datetime{
		// 	Day:   25,
		// 	Month: 1,
		// 	Year:  2024,
		// },
		End: &end, // &datetime.Datetime{
		// 	Day:   27,
		// 	Month: 1,
		// 	Year:  2024,
		// },
		Interval: interval, // datetime.OneHour
	}
	iter := chart.Get(quoteParams)
	var results []Quotes
	for iter.Next() {
		point := iter.Bar()
		q := Quotes{
			Open:      point.Open,
			Close:     point.Close,
			High:      point.High,
			Low:       point.Low,
			AdjClose:  point.AdjClose,
			Volume:    point.Volume,
			Timestamp: point.Timestamp,
		}
		// fmt.Println()
		results = append(results, q)
	}
	c.JSON(http.StatusOK, results)
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
