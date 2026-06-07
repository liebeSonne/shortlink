package audit

// Observer - интерфейс наблюдателя.
type Observer interface {
	// Update - получение события.
	Update(event Event)
}
