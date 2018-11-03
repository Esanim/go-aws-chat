package chatsess

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	// DateFmt is a date format.
	DateFmt = "02-01-2006"
)

// TimetoDB converts time to a string pointer.
func TimetoDB(t time.Time) *string {
	tn := t.Unix()
	return aws.String(strconv.FormatInt(tn, 10))
}

// DBtoTime converts a string pointer to time.
func DBtoTime(s *string) time.Time {
	n, _ := strconv.ParseInt(*s, 10, 64)
	return time.Unix(n, 0)
}
