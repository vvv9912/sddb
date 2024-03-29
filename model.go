package sddb

import (
	"time"
)

type ExampleRequest struct {
	Login    string `json:"login_login" form:"login_login" query:"login_login"`
	Password string `json:"password_password" form:"password_password" query:"password_password"`
}

type ShopCart struct {
	ID        int       `json:"id,omitempty"`
	TgId      int64     `json:"tg_id,omitempty"`
	Article   int       `json:"article,omitempty"` //В наличии
	Quantity  int       `json:"quantity,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Orders struct {
	ID            int       `json:"id,omitempty" db:"id"`
	TgID          int64     `json:"tg_id,omitempty" db:"tg_id"`
	UserName      string    `json:"user_name,omitempty" db:"user_name"`
	FirstName     string    `json:"first_name,omitempty" db:"first_name"`
	LastName      string    `json:"last_name,omitempty" db:"last_name"`
	StatusOrder   int       `json:"status_order,omitempty" db:"status_order"`
	Pvz           string    `json:"pvz,omitempty" db:"pvz"`
	Order         string    `json:"order,omitempty" db:"orderr"` // структура из OrderCorz
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdateAt      time.Time `json:"update_at" db:"update_at"`
	TypeDostavka  int       `json:"type_dostavka" db:"type_dostavka"`
	PriceDelivery float64   `json:"Price_Delivery" db:"price_delivery"`
	PriceFull     float64   `json:"Price_Full" db:"price_full"`
}
type OrderCorz struct {
	ID       int     `json:"id,omitempty"`
	TgId     int64   `json:"tg_id,omitempty"`
	Article  int     `json:"article,omitempty"` //В наличии
	Quantity int     `json:"quantity,omitempty"`
	Price    float64 `json:"Price,omitempty"`
	Name     string  `json:"Name,omitempty"`
	//CreatedAt time.Time `json:"created_at"`
}
type Users struct {
	id         int   `json:"id,omitempty"`
	TgID       int64 `json:"tg_id,omitempty"`
	StatusUser int   `json:"status_user,omitempty"`
	StateUser  int   `json:"state_user,omitempty"`
	//Corzine    []int     `json:"corzine,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Products struct {
	Article      int      `json:"article,omitempty"`
	Catalog      string   `json:"catalog,omitempty"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	PhotoUrl     [][]byte `json:"photo_url"`
	Price        float64  `json:"price,omitempty"`
	Length       int      `json:"length"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Weight       int      `json:"weight"`
	Availability bool     `json:"availability,omitempty"`
}
