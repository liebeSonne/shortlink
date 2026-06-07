package model

// ShortLink - модель сокращенных ссылок
type ShortLink struct {
	// ID - идентификатор - уникальное сокращенное имя ссылки.
	ID string
	// URL - сохраняемая полная ссылка.
	URL string
}
