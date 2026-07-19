package model

// ShortLink - модель сокращенных ссылок
type ShortLink struct {
	// ID - идентификатор - уникальное сокращенное имя ссылки.
	ID string
	// URL - сохраняемая полная ссылка.
	URL string
}

// StatsData - данные статистики сокращенных ссылок
type StatsData struct {
	CountShortLinks int
	CountUsers      int
}
