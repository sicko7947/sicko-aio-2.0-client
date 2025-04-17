package utils

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var sensors []string

func init() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Open("sensors.txt")
	if err != nil {
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	rawData, _ := ioutil.ReadAll(file)
	raw := strings.Split(strings.Replace(string(rawData), "\r\n", "", -1), "2,i,")
	for _, v := range raw {
		sensors = append(sensors, "2,i,"+v)
	}
}
