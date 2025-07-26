package gag

type Item struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	RestockTime int    `json:"restockTime"`
}
