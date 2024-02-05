package domain

type TradeItems struct {
	Items []TradeItem `json:"items"`
}

type TradeSignals struct {
	Items []TradeSignal `json:"items"`
}

type TradeItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type TradeSignal struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}