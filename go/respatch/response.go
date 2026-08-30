package respatch

import "encoding/json"

type apiEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func unwrapData(raw json.RawMessage) json.RawMessage {
	var env apiEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return raw
	}
	if env.Data != nil {
		return env.Data
	}
	return raw
}

func decodeData[T any](raw json.RawMessage, out *T) error {
	return json.Unmarshal(unwrapData(raw), out)
}

func decodeList[T any](raw json.RawMessage) ([]T, error) {
	data := unwrapData(raw)

	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err == nil && len(tuple) > 0 {
		var probe []T
		if err := json.Unmarshal(tuple[0], &probe); err == nil {
			return probe, nil
		}
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}
