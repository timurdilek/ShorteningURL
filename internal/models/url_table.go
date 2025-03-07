package models

import "time"

type UrlMap struct {
	Url, ShortenedUrl string
	Date              time.Time
}
