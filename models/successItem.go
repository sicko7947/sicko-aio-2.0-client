package models

type SuccessItem struct {
	Category    CATEGORY `json:"category"`
	Region      string   `json:"region"`
	Sku         string   `json:"sku"`
	ProductName string   `json:"productName"`
	OrderNumber string   `json:"orderNumber"`
	Email       string   `json:"email"`
	Size        string   `json:"size"`
	Timstamp    string   `json:"timestamp"`
	ImageUrl    string   `json:"imageUrl"`
}
