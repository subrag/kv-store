package core

import (
	"fmt"
)

func convertToCmdArgs(rawDecodedCmd interface{}) (string, []string) {
	arr := rawDecodedCmd.([]interface{})
	res := make([]string, len(arr))
	for i, v := range arr {
		if str, ok := v.(string); ok {
			res[i] = str
		} else {
			// Handle the case where the element is not a string (e.g., convert to a default value or ignore)
			res[i] = "default"
		}
	}
	return res[0], res[1:]
}

// func commandParser(b []byte) {
// 	cmd, _, _ := readArray(b[:])

// }

func HandlerQuery(b []byte, db *KV) (string, error) {
	command, _ := Decode(b)
	cmd, args := convertToCmdArgs(command)

	switch cmd {
	case "PING":
		return EncodeString("PONG"), nil
	case "COMMAND":
		return EncodeString(""), nil
	case "SET":
		db.set(args[0], args[1])
		return EncodeString("OK"), nil
	case "GET":
		v := db.get(args[0])
		return EncodeString(v), nil
	default:
		return "-\r\n", fmt.Errorf("unknown command %v", command)
	}

}
