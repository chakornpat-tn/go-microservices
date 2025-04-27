package utils

import (
	"encoding/json"
	"fmt"
)

func Debug(object any) {
	raw, _ := json.MarshalIndent(object, "", "\t")
	fmt.Println(raw)
}
