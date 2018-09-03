package prometheus

import (
	"regexp"

	"github.com/iancoleman/strcase"
)

func Sanitize(s string) string {
	// This function takes metric and tag names and munges them into prometheus
	// compatible format: [a-zA-Z_:][a-zA-Z0-9_:]*
	beginsWithNumberReg := regexp.MustCompile("^[0-9]+")
	characterReg := regexp.MustCompile("[^a-zA-Z0-9_:]")

	snake := strcase.ToSnake(s)
	snake = characterReg.ReplaceAllString(snake, "_")

	if beginsWithNumberReg.MatchString(snake) {
		snake = "_" + snake
	}

	return snake
}
