package main

import (
	"errors"
	"unicode"
)

const (
	StartRune rune = 'a'
	EndRune   rune = 'z'
)

func convertRuneToIndex(r rune) (int, error) {
	lr := unicode.ToLower(r)
	if lr < StartRune || lr > EndRune {
		return -1, errors.New("rune is not between (a -> Z)")
	}

	return int(lr - StartRune), nil
}

func getRuneDifference(a rune, b rune) (int, error) {
	c, err := convertRuneToIndex(a)
	if err != nil {
		return -1, err
	}

	d, err := convertRuneToIndex(b)
	if err != nil {
		return -1, err
	}

	result := c - d
	if result < 0 {
		result = -result
	}

	return result, nil
}
