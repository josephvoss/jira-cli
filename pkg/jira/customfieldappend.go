package jira

import (
	"fmt"
	"encoding/json"

	"regexp"
)

// Would love this to be the generic marshaller for issue, not sure how to load
// curr keys / struct tags
func populateCustomFields(b []byte, iss *Issue) error {
	// temp map holding all fields from json
	var tempMap map[string]interface{}
	if err := json.Unmarshal(b, &tempMap); err != nil {
		return err
	}

	// TODO temp map -> issue

	// build custom field type, append to issue
	// need to input / lookup map of keys / ids to custom field names
	// or nvm, just return custom_id : string
	reg := regexp.MustCompile("^customfield_[0-9]+$")
	// iter over every field
	fields, ok := tempMap["fields"].(map[string]interface{})
	if !ok {
		// TODO return error
		return nil
	}
	iss.Fields.CustomFields = make(map[string]string)
	for k, v := range fields {
		// replace w/ regex
		if reg.MatchString(k) {
			switch v := v.(type) {
			case string:
				iss.Fields.CustomFields[k] = v
			case int:
				iss.Fields.CustomFields[k] = string(v)
			case nil:
				//fmt.Printf("field %s is null\n", k)
			case map[string]interface{}:
				// if value key exists and string
				if val, ok := v["value"]; ok {
					if val, ok := val.(string); ok {
						iss.Fields.CustomFields[k] = val
					}
				}
			case []map[string]interface{}:
				// should recurse here.
				// [{"value":"str"}]
				// iter slice, if "value" key and v is a string
				// TODO last value overwrites first. Prob should be slice string we
				// append to. Change type from str:str to str:[str]?
				for _, i := range v {
					for key, val := range i {
						if val, ok := val.(string); ok && key == "value" {
									iss.Fields.CustomFields[k] = val
						}
					}
				}
			default:
				fmt.Printf("field %s type is unknown: %T\n", k, v)
			}
		}
	}

	return nil
}
