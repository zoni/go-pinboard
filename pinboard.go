package pinboard

import "fmt"
import "time"

type Pinboard struct {
	user  string
	token string
}

type Pin struct {
	url         string
	description string
	hash        string
	tags        []string
	date        time.Time
}

func (p *Pinboard) GetPins() []Pin {
	fmt.Println("hi")
	Pins := make([]Pin, 3)
	return Pins
}

func (p *Pinboard) GetAllPins() []Pin {
	Pins := make([]Pin, 3)
	return Pins
}
