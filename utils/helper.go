package utils

import (
	"context"
	"errors"
	"strconv"
	"time"
)

// StringToInt converts a string to an int
func StringToInt(str string) (int, error) {
	if str == "" {
		return 0, errors.New("empty string")
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("invalid string")
	}
	return num, nil
}

func GetContext() (context.Context, context.CancelFunc) {
	var ctx = context.Background()
	return context.WithTimeout(ctx, time.Second*60)
}
