package jira

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// Populate Issue.CustomFields with a map of `field_name` to a string value.
func populateCustomFields(b []byte, iss *Issue) error {
	var tempMap map[string]interface{}
	if err := json.Unmarshal(b, &tempMap); err != nil {
		return err
	}

	reg := regexp.MustCompile("^customfield_[0-9]+$")
	fields, ok := tempMap["fields"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse fields for %s", iss.Key)
	}

	iss.Fields.CustomFields = make(map[string]string)
	for k, v := range fields {
		if reg.MatchString(k) {
			switch v := v.(type) {
			case string:
				iss.Fields.CustomFields[k] = v
			case int:
				iss.Fields.CustomFields[k] = strconv.Itoa(v)
			case map[string]interface{}:
				if val, ok := v["value"]; ok {
					if val, ok := val.(string); ok {
						iss.Fields.CustomFields[k] = val
					}
				}
			case []map[string]interface{}:
				for _, i := range v {
					for key, val := range i {
						if val, ok := val.(string); ok && key == "value" {
							iss.Fields.CustomFields[k] = val
						}
					}
				}
			default:
				// Skip custom field if cannot convert to string match
			}
		}
	}

	return nil
}
