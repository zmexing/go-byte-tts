package internal

import "unicode/utf8"

// SplitText 将文本分割成每个元素最大maxBytes字节的切片
func SplitText(text string, maxBytes int) []string {
	//var chunks []string
	//for len(text) > 0 {
	//	if len(text) > chunkSize {
	//		chunks = append(chunks, text[:chunkSize])
	//		text = text[chunkSize:]
	//	} else {
	//		chunks = append(chunks, text)
	//		break
	//	}
	//}
	//return chunks

	var result []string
	start := 0
	for start < len(text) {
		end := start + maxBytes
		if end > len(text) {
			end = len(text)
		} else {
			for end > start && !utf8.RuneStart(text[end]) {
				end--
			}
		}
		result = append(result, text[start:end])
		start = end
	}
	return result
}
