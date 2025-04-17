package utils

import "github.com/brianvoe/gofakeit/v6"

func GetSensor() string {
	return gofakeit.RandomString(sensors)
}
