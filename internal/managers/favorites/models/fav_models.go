package models

// FavoriteRequest — тело запроса для добавления в избранное
type FavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

// FavoriteItem — один элемент в списке избранного
type FavoriteItem struct {
	AddedAt       string        `json:"added_at"`
	Advertisement Advertisement `json:"advertisement"`
}

// Advertisement — вложенный объект объявления (упрощённая версия)
type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
}

// FavoritesResponse — ответ от GET /favorites
type FavoritesResponse struct {
	Items []FavoriteItem `json:"items"`
}

// FavoriteAddResponse — ответ от POST /favorites
type FavoriteAddResponse struct {
	AddedAt       string        `json:"added_at"`
	Advertisement Advertisement `json:"advertisement"`
}
