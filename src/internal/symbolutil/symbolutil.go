package symbolutil

import (
	"delta-core/internal"
	"time"

	"github.com/piquette/finance-go/datetime"
)

type QuoteParams struct {
	Symbol   string `form:"symbol" binding:"required"`
	Interval string `form:"interval" default:"60"`
	Start    string `form:"start" default:"2023-01-02T15:04:05Z"`
	End      string `form:"end" default:"2023-01-03T15:04:05Z"`
}

func (qp *QuoteParams) FormatSymbol() (datetime.Interval, datetime.Datetime, datetime.Datetime) {
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
	case "60":
		interval = datetime.OneHour
	case "1":
		interval = datetime.OneMin
	case "1440":
		interval = datetime.OneDay
	case "5":
		interval = datetime.FiveMins
	case "15":
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
