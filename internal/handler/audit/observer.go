package audit

type Observer interface {
	Update(event Event)
}
