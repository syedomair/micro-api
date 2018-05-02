package common

import (
	"errors"
	"strconv"
)

func ValidateQueryString(limit string, defaultLimit string, offset string, defaultOffset string, orderby string, defaultOrderby string, sort string, defaultSort string) (string, string, string, string, error) {

	if limit != "" {
		if _, err := strconv.Atoi(limit); err != nil {
			return "", "", "", "", errors.New("Invalid 'limit' number in query string. Must be a number. ")
		}
	} else {
		limit = defaultLimit
	}
	if offset != "" {
		if _, err := strconv.Atoi(offset); err != nil {
			return "", "", "", "", errors.New("Invalid 'offset' number in query string. Must be a number. ")
		}
	} else {
		offset = defaultOffset
	}

	if orderby != "" {
		if _, err := strconv.Atoi(orderby); err == nil {
			return "", "", "", "", errors.New("Invalid 'orderby' value in query string. Must be a string. ")
		}
	} else {
		orderby = defaultOrderby
	}

	if sort != "" {
		if _, err := strconv.Atoi(sort); err == nil {
			return "", "", "", "", errors.New("Invalid 'sort' value in query string. Must be a string. ")
		}
		if (sort != "asc") && (sort != "desc") {
			return "", "", "", "", errors.New("Invalid 'sort' value in query string. Must be either 'asc' or 'desc'. ")
		}
	} else {
		sort = defaultSort
	}

	return limit, offset, orderby, sort, nil
}
