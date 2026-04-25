package domain

import (
	"fmt"
	"strings"
)

type Category string

const (
	CategoryGame  Category = "game"
	CategoryMovie Category = "movie"
	CategoryTV    Category = "tv"
)

func ParseCategory(raw string) (Category, error) {
	switch Category(strings.TrimSpace(strings.ToLower(raw))) {
	case CategoryGame:
		return CategoryGame, nil
	case CategoryMovie:
		return CategoryMovie, nil
	case CategoryTV:
		return CategoryTV, nil
	default:
		return "", fmt.Errorf("invalid category %q: must be one of game|movie|tv", raw)
	}
}
