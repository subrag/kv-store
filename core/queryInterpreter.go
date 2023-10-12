package core

import (
	"errors"
	"fmt"
	"strings"
)

type CmdArgs struct {
	cmd  string
	args []string
}

func convertToCmdArgs(rawDecodedCmd interface{}) CmdArgs {
	arr := rawDecodedCmd.([]interface{})
	res := make([]string, len(arr))
	for i, v := range arr {
		if str, ok := v.(string); ok {
			res[i] = str
		} else {
			// Handle the case where the element is not a string (e.g., convert to a default value or ignore)
			res[i] = "---"
		}
	}
	return CmdArgs{cmd: strings.ToUpper(res[0]), args: res[1:]}
}

func HandlerQuery(b []byte, db *KV) (string, error) {
	command, _ := Decode(b)
	if command == nil {
		return "", errors.New("no data")
	}
	cmdArgs := convertToCmdArgs(command)

	switch cmdArgs.cmd {
	case "PING":
		return EncodeString("PONG"), nil
	case "COMMAND":
		return EncodeString(""), nil
	case "SET":
		if len(cmdArgs.args) != 2 {
			return "-\r\n", fmt.Errorf("unable to process command %v", command)
		}
		db.set(cmdArgs.args[0], cmdArgs.args[1])
		return EncodeString("OK"), nil
	case "GET":
		v := db.get(cmdArgs.args[0])
		return EncodeString(v), nil
	default:
		return "-\r\n", fmt.Errorf("unknown command %v", command)
	}
}
