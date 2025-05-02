package uiutil

import "encoding/json"

func JSON(val any) string {
	b, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return string(b)
}
