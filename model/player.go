package model

type Player struct {
	Name      string
	Space     Space
	Direction string
	WasHit    bool
	Score     int
}
