package cond

import (
	"errors"
	"upper.io/db.v3"
)

func HandleParams(params map[string]interface{}) (db.Cond, error) {
	cond := make(db.Cond)

	for k, v := range params {
		if _, ok := v.([]interface{}); ok {
			cond[k] = db.In(v)
			continue
		} else if val, ok := v.(map[string]interface{}); ok {
			if val["type"] == nil || val["value"] == nil {
				continue
			}
			cond, err := packageParams(cond, k, val)
			if err != nil {
				return cond, err
			}
		} else {
			cond[k] = db.Eq(v)
		}
	}

	return cond, nil
}

func packageParams(cond db.Cond, key string, value map[string]interface{}) (db.Cond, error) {
	switch value["type"] {
	case "between":
		reqVal, ok := value["value"].([]interface{})
		if !ok {
			return cond, errors.New("The argument must be array in call to 'between'")
		}
		if len(reqVal) < 2 {
			return cond, errors.New("Not enough arguments in call to 'between'")
		}
		cond[key] = db.Between(reqVal[0], reqVal[1])
	case "lt":
		cond[key] = db.Lt(value["value"])
	case "lte":
		cond[key] = db.Lte(value["value"])
	case "gt":
		cond[key] = db.Gt(value["value"])
	case "gte":
		cond[key] = db.Gte(value["value"])
	case "eq":
		cond[key] = db.Eq(value["value"])
	case "like":
		reqVal, ok := value["value"].(string)
		if !ok {
			return cond, errors.New("The argument must be string in call to 'like'")
		}
		cond[key] = db.Like("%" + reqVal + "%")
	case "in":
		reqVal, ok := value["value"].([]interface{})
		if !ok {
			return cond, errors.New("The argument must be array in call to 'in'")
		}
		cond[key] = db.In(reqVal)
	}

	return cond, nil
}
