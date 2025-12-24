package helpers

import (
	"regexp"
	"strings"
)

// GetInt safely extracts an int from a map
func GetInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return 0
}

// GetFloat safely extracts a float64 from a map
func GetFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

// GetString safely extracts a string from a map
func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// CleanName sanitizes a string for SQL
func CleanName(name string) string {
	return strings.ReplaceAll(name, "'", "''")
}

// CleanItemName removes color codes and extra formatting from item names
func CleanItemName(name string) string {
	// Remove WoW color codes like |cff123456
	re := regexp.MustCompile(`\|c[0-9a-fA-F]{8}`)
	name = re.ReplaceAllString(name, "")
	// Remove closing tag |r
	name = strings.ReplaceAll(name, "|r", "")
	// Remove item link codes like |Hitem:...
	reLink := regexp.MustCompile(`\|H[^|]+\|h`)
	name = reLink.ReplaceAllString(name, "")
	name = strings.ReplaceAll(name, "|h", "")
	return strings.TrimSpace(name)
}

// FormatSpellDesc formats a spell description with base point placeholders
func FormatSpellDesc(desc string, bps []int) string {
	if desc == "" {
		return ""
	}

	// Replace $sN with actual values
	for i, bp := range bps {
		placeholder := ""
		switch i {
		case 0:
			placeholder = "$s1"
		case 1:
			placeholder = "$s2"
		case 2:
			placeholder = "$s3"
		}
		if placeholder != "" {
			// Handle negative values (usually shown as positive)
			val := bp
			if val < 0 {
				val = -val
			}
			desc = strings.ReplaceAll(desc, placeholder, formatNumber(val))
		}
	}

	// Clean up any remaining placeholders
	cleanups := []string{"$s1", "$s2", "$s3", "$d", "$D", "$o1", "$o2", "$o3"}
	for _, c := range cleanups {
		desc = strings.ReplaceAll(desc, c, "X")
	}

	return desc
}

func formatNumber(n int) string {
	if n < 0 {
		return "-" + formatNumber(-n)
	}
	if n < 1000 {
		return strings.TrimLeft(strings.Repeat("0", 3)+itoa(n), "0")
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
