package requests

import "encoding/json"

type Headers map[string]string

func (h Headers) Clone() Headers {
	if h == nil {
		return nil
	}
	hh := Headers{}
	for k, v := range h {
		hh[k] = v
	}
	return hh
}

func ParseStruct(h Headers, v interface{}) Headers {
	data, err := json.Marshal(v)
	if err != nil {
		return h
	}

	err = json.Unmarshal(data, &h)
	return h
}

func HeaderFromStruct(v interface{}) Headers {
	var header Headers
	header = ParseStruct(header, v)
	return header
}
