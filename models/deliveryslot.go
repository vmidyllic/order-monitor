package models

type DeliverySlot struct {
	Date  string  `json:"date"`
	Items []Items `json:"items"`
}
type Items struct {
	ID              string      `json:"id"`
	EndOrderingTime interface{} `json:"end_ordering_time"`
	TimeRange       string      `json:"time_range"`
	Price           int         `json:"price"`
	Currency        string      `json:"currency"`
	IsOpen          bool        `json:"is_open"`
	Date            string      `json:"date"`
}
