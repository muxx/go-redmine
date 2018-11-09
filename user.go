package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type userResult struct {
	User User `json:"user"`
}

type usersResult struct {
	Users []User `json:"users"`
}

type User struct {
	Id           int            `json:"id"`
	Login        string         `json:"login"`
	Firstname    string         `json:"firstname"`
	Lastname     string         `json:"lastname"`
	Mail         string         `json:"mail"`
	CreatedOn    string         `json:"created_on"`
	LatLoginOn   string         `json:"last_login_on"`
	Memberships  []Membership   `json:"memberships"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

type UserFilter struct {
	Status  int8
	Name    string
	GroupId int
}

func (c *Client) Users() ([]User, error) {
	return c.getUsers("/users.json?key=" + c.apikey + c.getPaginationClause())
}

func (c *Client) UsersByFilter(f *UserFilter) ([]User, error) {
	filterString := ""
	if f != nil {
		if f.Status > 0 {
			filterString = filterString + fmt.Sprintf("&status=%d", f.Status)
		}
		if f.Name != "" {
			filterString = filterString + fmt.Sprintf("&name=%s", f.Name)
		}
		if f.GroupId > 0 {
			filterString = filterString + fmt.Sprintf("&group_id=%d", f.GroupId)
		}
	}

	return c.getUsers("/users.json?key=" + c.apikey + c.getPaginationClause() + filterString)
}

func (c *Client) User(id int) (*User, error) {
	res, err := c.Get(c.endpoint + "/users/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r userResult
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
	return &r.User, nil
}

func (c *Client) getUsers(url string) ([]User, error) {
	res, err := c.Get(c.endpoint + url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r usersResult
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
	return r.Users, nil

}
