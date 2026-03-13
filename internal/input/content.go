package input

import "strings"

func MergeTags(content string, tags []string) string {
	result := strings.TrimSpace(content)
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if !strings.HasPrefix(tag, "#") {
			tag = "#" + tag
		}
		if !strings.Contains(result, tag) {
			if result == "" {
				result = tag
			} else {
				result += "\n " + tag
			}
		}
	}
	return result
}

func RemoveTag(content string, tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return strings.TrimSpace(content)
	}
	if !strings.HasPrefix(tag, "#") {
		tag = "#" + tag
	}
	result := strings.ReplaceAll(content, tag, "")
	result = strings.Join(strings.Fields(result), " ")
	return strings.TrimSpace(result)
}
