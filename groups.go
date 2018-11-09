package redmine

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type groupsResult struct {
	Groups []Group `json:"groups"`
}

type groupResult struct {
	Group Group `json:"group"`
}

type Group struct {
	IdName
	Users       []IdName     `json:"users"`
	Memberships []Membership `json:"memberships"`
}

func (c *Client) Groups() ([]Group, error) {
	res, err := c.Get(c.endpoint + "/groups.json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r groupsResult
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
	return r.Groups, nil
}

func (c *Client) Group(id int, includeUsers bool, includeMemberships bool) (*Group, error) {
	var includeParams []string
	if includeUsers {
		includeParams = append(includeParams, "users")
	}
	if includeMemberships {
		includeParams = append(includeParams, "memberships")
	}
	var includeString string
	if len(includeParams) > 0 {
		includeString = "&include=" + strings.Join(includeParams, ",")
	}

	res, err := c.Get(c.endpoint + "/groups/" + strconv.Itoa(id) + ".json?key=" + c.apikey + includeString)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r groupResult
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
	return &r.Group, nil
}
