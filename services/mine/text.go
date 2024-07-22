package mine

import (
	"fmt"
	"strconv"
	"strings"
)

func intToStringWithSign(value int) string {
	if value < 0 {
		return strconv.Itoa(value)
	} else {
		return fmt.Sprintf("+%d", value)
	}
}

func (c *content) locationTemplate(title string, text string, items map[string]int) string {
	stringBuilder := strings.Builder{}

	stringBuilder.WriteString(fmt.Sprintf("## %s\n%s\n\n", title, text))

	for item, count := range items {
		stringBuilder.WriteString(fmt.Sprintf("%s %s\n", intToStringWithSign(count), item))
	}

	return stringBuilder.String()
}

func (c *content) endTemplate(bills map[string]int, title string, text string) string {
	stringBuilder := strings.Builder{}

	stringBuilder.WriteString(fmt.Sprintf("## %s\n%s\n\n", title, text))

	billsTotal := 0
	for message, bill := range bills {
		billsTotal += bill
		stringBuilder.WriteString(fmt.Sprintf(" %s (%s)\n", intToStringWithSign(bill), message))
	}

	stringBuilder.WriteString(fmt.Sprintf("**Итог:** %s", intToStringWithSign(billsTotal)))

	return stringBuilder.String()
}
