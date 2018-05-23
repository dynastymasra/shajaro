package helper

import "encoding/json"

func Stringify(str interface{}) string {
	out, err := json.Marshal(str)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
