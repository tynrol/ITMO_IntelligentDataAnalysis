package domain

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
)

func (w Weather) String() string {
	return [...]string{"Sunny", "Sun", "Cloudy", "Sunrise", "Sunset", "Rain"}[w]
}

func RandomWeather() string {
	k := rand.Intn(int(Rain))
	return Weather(k).String()
}
