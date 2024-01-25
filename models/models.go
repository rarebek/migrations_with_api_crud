package models

type User struct {
	Id      int    `json:"id"`
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Content string `json:"content"`
}
