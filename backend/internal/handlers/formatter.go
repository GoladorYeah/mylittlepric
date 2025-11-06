package handlers

import "strings"

func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getFloatValue(data map[string]interface{}, key string) float64 {
	switch v := data[key].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func getIntValue(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	if val, ok := data[key].(int); ok {
		return val
	}
	return 0
}

func getBoolValue(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

func extractPageTokenFromLink(serpAPILink string) string {
	if serpAPILink == "" {
		return ""
	}

	tokenStart := strings.Index(serpAPILink, "page_token=")
	if tokenStart == -1 {
		return ""
	}

	tokenStart += len("page_token=")
	tokenEnd := strings.Index(serpAPILink[tokenStart:], "&")

	if tokenEnd == -1 {
		return serpAPILink[tokenStart:]
	}

	return serpAPILink[tokenStart : tokenStart+tokenEnd]
}
