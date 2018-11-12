package redmine

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	CustomizedTypeIssue             = "issue"
	CustomizedTypeIssuePriority     = "issue_priority"
	CustomizedTypeTimeEntry         = "time_entry"
	CustomizedTypeTimeEntryActivity = "time_entry_activity"
	CustomizedTypeProject           = "project"
	CustomizedTypeVersion           = "version"
	CustomizedTypeDocument          = "document"
	CustomizedTypeDocumentCategory  = "document_category"
	CustomizedTypeUser              = "user"
	CustomizedTypeGroup             = "group"

	FieldFormatVersion     = "version"
	FieldFormatDate        = "date"
	FieldFormatText        = "text"
	FieldFormatBool        = "bool"
	FieldFormatUser        = "user"
	FieldFormatFloat       = "float"
	FieldFormatList        = "list"
	FieldFormatEnumeration = "enumeration"
	FieldFormatLink        = "link"
	FieldFormatString      = "string"
	FieldFormatInt         = "int"
)

type customFieldsDictionaryResult struct {
	CustomFields []CustomFieldDictionary `json:"custom_fields"`
}

type CustomFieldDictionary struct {
	Id             int                                  `json:"id"`
	Name           string                               `json:"name"`
	CustomizedType string                               `json:"customized_type"`
	FieldFormat    string                               `json:"field_format"`
	Regexp         string                               `json:"regexp"`
	MinLength      int                                  `json:"min_length"`
	MaxLength      int                                  `json:"max_length"`
	IsRequired     bool                                 `json:"is_required"`
	IsFilter       bool                                 `json:"is_filter"`
	Searchable     bool                                 `json:"searchable"`
	Multiple       bool                                 `json:"multiple"`
	DefaultValue   interface{}                          `json:"default_value"`
	Visible        bool                                 `json:"visible"`
	PossibleValues []CustomFieldDictionaryPossibleValue `json:"possible_values"`
	Trackers       []IdName                             `json:"trackers"`
	Roles          []IdName                             `json:"roles"`
}

type CustomFieldDictionaryPossibleValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (c *Client) CustomFields() ([]CustomFieldDictionary, error) {
	res, err := c.Get(c.endpoint + "/custom_fields.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r customFieldsDictionaryResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.CustomFields, nil
}
