package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Rules struct {
	birth [9]bool
	survive [9]bool
}


// Recieves the cell state an neighbour number and decides its fate
func (r* Rules) DecideFate(alive bool, neighbours uint8) bool {
	if alive {
		return r.survive[neighbours]
	} else {
		return r.birth[neighbours]
	}
}

// Parses rules string into the struct
func NewRules(rules string) (*Rules, error) {
	r := &Rules{}
	
	m, err := regexp.Match("^B\\d+/S\\d+$", []byte(rules))

	if err != nil {
		return nil, fmt.Errorf("Regex failed: %w", err)
	}
	if m == false {
		return nil, fmt.Errorf("Rules string is not valid: %s", rules)
	}

	// We know it has "/" as it passed the regex
	b, s, _ := strings.Cut(rules, "/")

	if err := set(strings.TrimPrefix(b, "B"), &r.birth); err != nil {
		return nil, fmt.Errorf("Could not set birth rules: %w", err)
	}
	if err := set(strings.TrimPrefix(s, "S"), &r.survive); err != nil {
		return nil, fmt.Errorf("Could not set survive rules: %w", err)
	}

	return r, nil
}

// Sets to true the values of birth or survive
func set(s string, target *[9]bool) error {
	for _, v := range s {
		res, err := strconv.Atoi(string(v))

		if err != nil {
			return fmt.Errorf("Could not convert ASCII to int: %w", err)
		}

		if res > 8 {
			return fmt.Errorf("Invalid digit %d", res)
		}

		target[res] = true
	}

	return nil
}
