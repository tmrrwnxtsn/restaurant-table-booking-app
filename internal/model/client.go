package model

// Client представляет клиента ресторана, бронирующего столики.
type Client struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
