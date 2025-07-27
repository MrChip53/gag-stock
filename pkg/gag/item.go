package gag

type Item struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	StockTime   int    `json:"stockTime"`
	RestockTime int    `json:"restockTime"`
}
