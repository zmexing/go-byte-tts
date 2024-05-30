package internal

import "errors"

// CheckParams 检查参数传递
func CheckParams(params map[string]map[string]any) error {
	if params["audio"] == nil {
		params["audio"] = make(map[string]any)
	}
	if _, ok := params["audio"]["voice_type"]; !ok {
		return errors.New("audio.voice_type cannot be empty")
	}
	if params["request"] == nil {
		params["request"] = make(map[string]any)
	}
	_, ok := params["request"]["text"]
	if !ok {
		return errors.New("request.text cannot be empty")
	}
	return nil
}

// DeepCopyParams 参数map值深拷贝
func DeepCopyParams(original map[string]map[string]any) map[string]map[string]any {
	copyData := make(map[string]map[string]any)
	for k, v := range original {
		innerCopy := make(map[string]any)
		for innerK, innerV := range v {
			innerCopy[innerK] = innerV
		}
		copyData[k] = innerCopy
	}
	return copyData
}
