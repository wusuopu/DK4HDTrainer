package trainer

import (
	"fmt"
	"math"
)

// 解析经度
func parseLongitude(integral uint8, fractional uint8) float32 {
	// 0x4E00 差不多是经度0，往东数值增加，0x1 差不多是 2.3度
	var value float32 = (float32(int16(integral) - int16(0x4E)) * 2.3) + (float32(fractional) / 256.0)
	return value
}
// 解析纬度
func parseLatitude(integral uint8, fractional uint8) float32 {
	// 0x2700 差不多是纬度0，往南数值增加，0x1 差不多是 2.25度
	var value float32 = (float32(int16(integral) - int16(0x27)) * 2.25) + (float32(fractional) / 256.0)
	return value
}

func formatLongitude(value float32) string {
	pos := ""
	if value > 0 {
		pos = "E"
	} else {
		pos = "W"
	}
	return fmt.Sprintf("%s%.2f", pos, math.Abs(float64(value)))
}

func formatLatitude(value float32) string {
	pos := ""
	if value > 0 {
		pos = "S"
	} else {
		pos = "N"
	}

	return fmt.Sprintf("%s%.2f", pos, math.Abs(float64(value)))
}