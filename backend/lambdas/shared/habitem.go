package shared

type HABItem struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

/*
	List response:
	[
		{
			"link": "http://192.168.86.218:8080/rest/items/Lock000C4C63_BatteryLevel",
			-- "state": "82",
			"stateDescription": {
				"minimum": 0,
				"maximum": 100,
				"step": 1,
				"pattern": "%.0f %%",
				"readOnly": true,
				"options": []
			},
			"editable": true,
			"type": "Number",
			-- "name": "Lock000C4C63_BatteryLevel",
			"label": "Battery Level",
			"category": "Battery",
			"tags": [
				"Point"
			],
			"groupNames": []
		}
	]
*/
