package utils

import (
	"fmt"
	"strconv"
)

func ConvertStringToFloat32(str string) (float32, error) {
	var result float64
	result, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to float32: %w", err)
	}
	return float32(result), nil
}

func ConvertFloat32ToString(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}
