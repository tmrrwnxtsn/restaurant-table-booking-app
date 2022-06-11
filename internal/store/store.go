package store

// Store представляет слой хранения данных (база данных).
type Store interface {
	// Restaurants позволяет обратиться к таблице с информацией о ресторанах.
	Restaurants() RestaurantRepository
	// Tables позволяет обратиться к таблице с информацией о столиках в ресторане.
	Tables() TableRepository
	// Clients позволяет обратиться к таблице с информацией о клиентах ресторанов.
	Clients() ClientRepository
	// Bookings позволяет обратиться к таблице с информацией о совершённых клиентами бронях.
	Bookings() BookingRepository
}
