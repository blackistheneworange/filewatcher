package logger

import (
	"fmt"
	"os"
)

var prefix string = "[filewatcher]"

func prepare(items ...interface{}) []interface{} {
	itemsToPrint := make([]interface{}, 0, len(items)+1)
	itemsToPrint = append(itemsToPrint, prefix)
	for _, item := range items {
		itemsToPrint = append(itemsToPrint, item);
	}

	return itemsToPrint
}

func Log(items ...interface{}) {
	fmt.Println(prepare(items...)...);
}

func Fatal(items ...interface{}) {
	fmt.Println(prepare(items...)...);
	os.Exit(1)
}

func Format(separator rune, items ...interface{}) string{
	formatters := ""
	finalItems := make([]interface{}, 0, len(items)*2-1)

	for idx := range items {
		formatters += "%v"
		finalItems = append(finalItems, items[idx])
		if idx < len(items)-1 {
			formatters += "%c"
			finalItems = append(finalItems, separator)
		}
	}
	return fmt.Sprintf(formatters, finalItems...)
}