package ws

import "encoding/json"

type Command struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func ParseCommand(cmdBytes []byte) (Command, error) {
	cmd := Command{}
	err := json.Unmarshal(cmdBytes, &cmd)
	if err != nil {
		return cmd, err
	}

	return cmd, nil
}
