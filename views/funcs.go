package views

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

var funcs = template.FuncMap{
	"asset":     assetURLFn,
	"url":       relURLFn,
	"lower":     lower,
	"upper":     upper,
	"date":      dateFn,
	"humanDate": readableDateFn,
	"t":         translateFn,
}

func lower(input string) string {
	return strings.ToLower(input)
}

func upper(input string) string {
	return strings.ToUpper(input)
}

func assetURLFn(input string) string {
	url := "/static/" //os.Getenv("ASSET_URL")
	return fmt.Sprintf("%s%s", url, input)
}

func relURLFn(input string) string {
	url := "/" //os.Getenv("ASSET_URL")
	return fmt.Sprintf("%s%s", url, input)
}

func dateFn(format string, input interface{}) string {
	var t time.Time
	switch date := input.(type) {
	default:
		t = time.Now()
	case time.Time:
		t = date
	}

	return t.Format(format)
}

func translateFn(language string, text string) string {
	return text
}

func readableDateFn(t time.Time) string {
	if time.Now().Before(t) {
		return "in the future"
	}
	diff := time.Now().Sub(t)
	day := 24 * time.Hour
	month := 30 * day
	year := 12 * month

	switch {
	case diff < time.Second:
		return "just now"
	case diff < 5*time.Minute:
		return "a few minutes ago"
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", diff/time.Minute)
	case diff < day:
		return fmt.Sprintf("%d hours ago", diff/time.Hour)
	case diff < month:
		return fmt.Sprintf("%d days ago", diff/day)
	case diff < year:
		return fmt.Sprintf("%d months ago", diff/month)
	default:
		return fmt.Sprintf("%d years ago", diff/year)
	}
}
