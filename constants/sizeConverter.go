package constants

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var (
	sizeChart gjson.Result = gjson.Parse(`[{"US": "6", "EU": "38.5", "UK": "5.5", "CM": "24"}, {"US": "6.5", "EU": "39", "UK": "6", "CM": "24.5"}, {"US": "7", "EU": "40", "UK": "6.5", "CM": "25"}, {"US": "7.5", "EU": "40.5", "UK": "7", "CM": "25.5"}, {"US": "8", "EU": "41", "UK": "7.5", "CM": "26"}, {"US": "8.5", "EU": "42", "UK": "8", "CM": "26.5"}, {"US": "9", "EU": "42.5", "UK": "8.5", "CM": "27"}, {"US": "9.5", "EU": "43", "UK": "9", "CM": "27.5"}, {"US": "10", "EU": "44", "UK": "9.5", "CM": "28"}, {"US": "10.5", "EU": "44.5", "UK": "10", "CM": "28.5"}, {"US": "11", "EU": "45", "UK": "10.5", "CM": "29"}, {"US": "11.5", "EU": "45.5", "UK": "11", "CM": "29.5"}, {"US": "12", "EU": "46", "UK": "11.5", "CM": "30"}, {"US": "12.5", "EU": "46.5", "UK": "12", "CM": "30.5"}, {"US": "13", "EU": "47.5", "UK": "12.5", "CM": "31"}]`)
)

func SizeConverter(size string, from string, to string) string {
	obj := sizeChart.Get(fmt.Sprintf(`#(%s=="%s").%s`, from, size, to))
	if obj.Exists() {
		return obj.String()
	}
	return size
}
