package main

import (
	"errors"
	"unicode"
)

const (
	startRune rune = 'a'
	endRune   rune = 'z'
)

func convertRuneToIndex(r rune) (int, error) {
	lr := unicode.ToLower(r)
	if lr < startRune || lr > endRune {
		return -1, errors.New("rune is not between (a -> Z)")
	}

	return int(lr - startRune), nil
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

	return abs(c - d), nil
}
