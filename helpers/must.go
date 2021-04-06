package helpers

import "time"

func TimeMust(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}
