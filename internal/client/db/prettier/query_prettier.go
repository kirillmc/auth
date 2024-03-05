package prettier

import "fmt"

const (
	PlaceholderDollar   = "$"
	PlaceHolderQuestion = "?"
)

func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			valuse := fmt.Sprintf("")
		}
	}
}
