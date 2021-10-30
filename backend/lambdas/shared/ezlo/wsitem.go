package ezlo

import (
	"encoding/json"
	"fmt"
	"mlock/lambdas/shared"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Most of this was taken from https://stackoverflow.com/questions/33436730/unmarshal-json-with-some-known-and-some-unknown-field-names

type wsItem struct {
	DeviceID string                 `json:"deviceId"`
	Name     string                 `json:"name"`
	Extra    map[string]interface{} `json:"-"`
	// "value": 1,
	/* "value": {
			"1": {
	            "duration": 5,
	            "name": "Default Siren"
	          }
	        },
	*/
	//value:map[1:map[code:********** mode:enabled name:Code 1]
}

/*
type userCodeItem struct {
	Value map[string]shared.UserCodeItemValue `json:"value"`
}
*/

type _wsItem wsItem

func (t wsItem) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	// Take everything in Extra
	for k, v := range t.Extra {
		data[k] = v
	}

	// Take all the struct values with a json tag
	val := reflect.ValueOf(t)
	typ := reflect.TypeOf(t)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldv := val.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			data[jsonTag] = fieldv.Interface()
		}
	}
	return json.Marshal(data)
}

func (t *wsItem) UnmarshalJSON(b []byte) error {
	t2 := _wsItem{}
	err := json.Unmarshal(b, &t2)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &(t2.Extra))
	if err != nil {
		return err
	}

	typ := reflect.TypeOf(t2)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			delete(t2.Extra, jsonTag)
		}
	}

	*t = wsItem(t2)

	return nil
}

func (i *wsItem) getBatteryLevel() (int, error) {
	if i.Name != "battery" {
		return -1, fmt.Errorf("wrong item type, name is: %s", i.Name)
	}

	j, err := json.Marshal(i.Extra["value"])
	if err != nil {
		return -1, fmt.Errorf("error marshalling value: %s", err.Error())
	}

	var battery int
	if err := json.Unmarshal(j, &battery); err != nil {
		return -1, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	return battery, nil
}

func (i *wsItem) getLockCodes() ([]shared.RawDeviceLockCode, error) {
	lcs := []shared.RawDeviceLockCode{}

	if i.Name != "user_codes" {
		return lcs, fmt.Errorf("wrong item type, name is: %s", i.Name)
	}

	j, err := json.Marshal(i.Extra["value"])
	if err != nil {
		return lcs, fmt.Errorf("error marshalling value: %s", err.Error())
	}

	lcm := map[string]shared.RawDeviceLockCode{}
	if err := json.Unmarshal(j, &lcm); err != nil {
		return lcs, fmt.Errorf("error unmarshalling: %s", err.Error())
	}

	for key, lc := range lcm {
		slot, err := strconv.Atoi(key)
		if err != nil {
			return lcs, fmt.Errorf("error getting slot: %s", err.Error())
		}

		lc.Slot = slot
		lcs = append(lcs, lc)
	}

	sort.Slice(lcs, func(i, j int) bool {
		return lcs[i].Slot < lcs[j].Slot
	})

	return lcs, nil
}
