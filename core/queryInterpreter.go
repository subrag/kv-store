package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	cmd  string
	args []string
}

func convertToCmdArgs(rawDecodedCmd interface{}) Command {
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
	return Command{cmd: strings.ToUpper(res[0]), args: res[1:]}
}

func HandlerQuery(b []byte, db *DB) (string, error) {
	commandStr, _ := Decode(b)
	if commandStr == nil {
		return "", errors.New("no data")
	}
	command := convertToCmdArgs(commandStr)

	switch command.cmd {
	case "PING":
		return EncodeString("PONG"), nil
	case "COMMAND":
		return EncodeString(""), nil
	case "SET":
		return setQueryHandler(command.args, db)
	case "GET":
		return getQueryHandler(command.args, db)
	case "DEL":
		return delQueryHandler(command.args, db)
	default:
		return RespErrUnknownCmc, fmt.Errorf("unknown command %v", command)
	}
}

func setQueryHandler(args []string, db *DB) (string, error) {
	if len(args) == 2 {
		db.Set(args[0], args[1], DefaultTTL)
		return EncodeString("OK"), nil
	}
	if len(args) == 4 && strings.ToUpper(args[2]) == "EX" {
		ttl, err := strconv.Atoi(args[3])
		if err != nil {
			return RespErrUnknownCmc, fmt.Errorf("unable to process command, error: %v", err.Error())
		}
		db.Set(args[0], args[1], int32(ttl))
		return EncodeString("OK"), nil
	}
	return RespErrUnknownCmc, fmt.Errorf("unable to process command")
}

func getQueryHandler(args []string, db *DB) (string, error) {
	v := db.Get(args[0])
	return EncodeString(v), nil
}

// ToDo: param args needs to be refactored.
func delQueryHandler(args []string, db *DB) (string, error) {
	ct := 0
	for _, k := range args {
		// count only if key exists and deleted by del method
		if db.Del(k) != "" {
			ct++
		}
	}
	return EncodeInt(ct), nil
}
