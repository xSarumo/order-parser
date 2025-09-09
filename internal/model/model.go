package model

import "time"

type Delivery struct {
	Name    string `json:"name" fake:"{name}" validate:"required"`
	Phone   string `json:"phone" fake:"{phone}" validate:"required"`
	Zip     string `json:"zip" fake:"{zip}" validate:"required"`
	City    string `json:"city" fake:"{city}" validate:"required"`
	Address string `json:"address" fake:"{street}" validate:"required"`
	Region  string `json:"region" fake:"{state}" validate:"required"`
	Email   string `json:"email" fake:"{email}" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" fake:"{regex:[a-z1-9]{10}}" validate:"required"`
	RequestID    string `json:"request_id,omitempty" fake:"{regex:[a-z1-9]{9}}" validate:"required"`
	Currency     string `json:"currency" fake:"{regex:[A-Z]{3}}" validate:"required"`
	Provider     string `json:"provider" fake:"{company}" validate:"required"`
	Amount       int    `json:"amount" fake:"{number:100,3000}" validate:"required,gt=0"`
	PaymentDT    int64  `json:"payment_dt" fake:"{number:10000000,99999999}" validate:"required,gt=0"`
	Bank         string `json:"bank" fake:"{company}" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:100,1000}" validate:"gte=0"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:50,500}" validate:"gte=0"`
	CustomFee    int    `json:"custom_fee" fake:"{number:0,50}" validate:"gte=0"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" fake:"{number:1000000,9999999}" validate:"required,gt=0"`
	TrackNumber string `json:"track_number" fake:"{regex:[A-Z]{14}}" validate:"required"`
	Price       int    `json:"price" fake:"{number:100,1000}" validate:"required,gt=0"`
	RID         string `json:"rid" fake:"{regex:[a-z1-9]{19}}" validate:"required"`
	Name        string `json:"name" fake:"{productname}" validate:"required"`
	Sale        int    `json:"sale" fake:"{number:0,100}" validate:"gte=0,lte=100"`
	Size        string `json:"size" fake:"{regex:[0-9]{2}}" validate:"required"`
	TotalPrice  int    `json:"total_price" fake:"{number:100,1200}" validate:"gte=0"`
	NmID        int64  `json:"nm_id" fake:"{number:1000000,9999999}" validate:"required,gt=0"`
	Brand       string `json:"brand" fake:"{company}" validate:"required"`
	Status      int    `json:"status" fake:"{number:100,299}" validate:"gte=0"`
}

type Order struct {
	OrderUID        string    `json:"order_uid" fake:"{regex:[a-z1-9]{19}}" validate:"required"`
	TrackNumber     string    `json:"track_number" fake:"{regex:[A-Z]{14}}" validate:"required"`
	Entry           string    `json:"entry" fake:"{regex:[A-Z]{4}}" validate:"required"`
	Delivery        Delivery  `json:"delivery" validate:"required"`
	Payment         Payment   `json:"payment"  validate:"required"`
	Items           []Item    `json:"items" fakesize:"1,3" validate:"required,dive"`
	Locale          string    `json:"locale" fake:"{regex:[A-Z]{3}}" validate:"required"`
	InternalSign    string    `json:"internal_signature" fake:"{regex:[a-z1-9]{10}}" validate:"required"`
	CustomerID      string    `json:"customer_id" fake:"{username}" validate:"required"`
	DeliveryService string    `json:"delivery_service" fake:"{word}" validate:"required"`
	ShardKey        string    `json:"shardkey" fake:"{digit}{digit}" validate:"required"`
	SmID            int       `json:"sm_id" fake:"{number:1,100}" validate:"required,gt=0"`
	DateCreated     time.Time `json:"date_created" fake:"{date}" validate:"required"`
	OofShard        string    `json:"oof_shard" fake:"{digit}" validate:"required"`
}
