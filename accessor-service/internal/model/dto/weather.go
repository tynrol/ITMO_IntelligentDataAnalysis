package dto

import (
	"math/rand"
)

type Weather int

const (
	Sunny Weather = iota
	Sun
	Cloudy
	Sunrise
	Sunset
	Rain
	Rainy
)

func (w Weather) String() string {
	return [...]string{"Sunny", "Sun", "Cloudy", "Sunrise", "Sunset", "Rain", "Rainy"}[w]
}

func RandomWeather() string {
	k := rand.Intn(int(Rainy + 1))
	return Weather(k).String()
}
