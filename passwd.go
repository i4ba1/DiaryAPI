package main

import (
	"fmt"
	"unicode"
)

func validPassword(s string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": 	{unicode.Upper, unicode.Title},
		"lower case": 	{unicode.Lower},
		"numeric":    	{unicode.Number, unicode.Digit},
		"special":    	{unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("password must have at least one %s character", name)
	}
	return nil
}

/*func main() {
	for _, s := range []string{
		"bad",
		"testPassword",
		"testPa##word",
		"b3tterPa$$w0rd",
	} {
		fmt.Printf("validPassword(%q) = %v\n", s, validPassword(s))
	}
}*/
