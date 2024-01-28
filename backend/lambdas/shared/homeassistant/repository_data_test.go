package homeassistant_test

import "encoding/json"

type entity struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged string                 `json:"last_changed"`
	LastUpdated string                 `json:"last_updated"`
	Context     struct {
		ID       string  `json:"id"`
		ParentID *string `json:"parent_id"`
		UserID   *string `json:"user_id"`
	} `json:"context"`
}

func getMockEntityData(id string) []byte {
	var entities []entity

	// Unmarshal the JSON data into the slice
	err := json.Unmarshal([]byte(getMockStatesData()), &entities)
	if err != nil {
		panic(err)
	}

	for _, entity := range entities {
		if entity.EntityID == id {
			result, err := json.Marshal(entity)
			if err != nil {
				panic(err)
			}
			return result
		}
	}
	panic("Entity not found")
}

func getMockStatesData() []byte {
	return []byte(`
		[
			{
				"entity_id": "binary_sensor.remote_ui",
				"state": "on",
				"attributes": {
					"device_class": "connectivity",
					"friendly_name": "Remote UI"
				},
				"last_changed": "2024-01-20T22:52:29.446048+00:00",
				"last_updated": "2024-01-20T22:52:29.446048+00:00",
				"context": {
					"id": "01HMMH6NC60HG87TGPTZNWC2MC",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "update.home_assistant_supervisor_update",
				"state": "off",
				"attributes": {
					"auto_update": true,
					"installed_version": "2023.12.1",
					"in_progress": false,
					"latest_version": "2023.12.1",
					"release_summary": null,
					"release_url": "https://github.com/home-assistant/supervisor/releases/tag/2023.12.1",
					"skipped_version": null,
					"title": "Home Assistant Supervisor",
					"entity_picture": "https://brands.home-assistant.io/hassio/icon.png",
					"friendly_name": "Home Assistant Supervisor Update",
					"supported_features": 1
				},
				"last_changed": "2024-01-20T20:22:24.801024+00:00",
				"last_updated": "2024-01-20T20:22:24.801024+00:00",
				"context": {
					"id": "01HMM8KVS19P7CFTY7M9807KWZ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "update.home_assistant_core_update",
				"state": "on",
				"attributes": {
					"auto_update": false,
					"installed_version": "2024.1.4",
					"in_progress": false,
					"latest_version": "2024.1.5",
					"release_summary": null,
					"release_url": "https://www.home-assistant.io/latest-release-notes/",
					"skipped_version": null,
					"title": "Home Assistant Core",
					"entity_picture": "https://brands.home-assistant.io/homeassistant/icon.png",
					"friendly_name": "Home Assistant Core Update",
					"supported_features": 11
				},
				"last_changed": "2024-01-20T22:17:24.243476+00:00",
				"last_updated": "2024-01-20T22:17:24.243476+00:00",
				"context": {
					"id": "01HMMF6DGKDFTHV0ANH0HDPFNT",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "update.z_wave_js_update",
				"state": "off",
				"attributes": {
					"auto_update": false,
					"installed_version": "0.4.3",
					"in_progress": false,
					"latest_version": "0.4.3",
					"release_summary": "# Changelog\n\n## 0.4.3\n\n### Features\n\n- Z-Wave JS Server: Enable server to listen on IPv6 interfaces\n\n### Bug fixes\n\n- Handle more cases of unexpected Serial API restarts\n\n### Config file changes\n\n- Add wakeup instructions for Nexia ZSENS930\n- Correct para",
					"release_url": null,
					"skipped_version": null,
					"title": "Z-Wave JS",
					"entity_picture": "/api/hassio/addons/core_zwave_js/icon",
					"friendly_name": "Z-Wave JS Update",
					"supported_features": 25
				},
				"last_changed": "2024-01-20T20:22:24.803203+00:00",
				"last_updated": "2024-01-20T20:22:24.803203+00:00",
				"context": {
					"id": "01HMM8KVS3P1EC7GBB0EJF6S0K",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "update.advanced_ssh_web_terminal_update",
				"state": "on",
				"attributes": {
					"auto_update": false,
					"installed_version": "17.0.3",
					"in_progress": false,
					"latest_version": "17.0.4",
					"release_summary": "## What’s changed\n\n## ⬆️ Dependency updates\n\n- ⬆️ Update OpenSSL to v3.1.4-r4 @renovate ([#676](https://github.com/hassio-addons/addon-ssh/pull/676))\n- ⬆️ Update ghcr.io/hassio-addons/base Docker tag to v15.0.5 @renovate ([#677](https://github.com/hassio-",
					"release_url": null,
					"skipped_version": null,
					"title": "Advanced SSH & Web Terminal",
					"entity_picture": "/api/hassio/addons/a0d7b954_ssh/icon",
					"friendly_name": "Advanced SSH & Web Terminal Update",
					"supported_features": 25
				},
				"last_changed": "2024-01-20T20:22:24.805641+00:00",
				"last_updated": "2024-01-20T20:22:24.805641+00:00",
				"context": {
					"id": "01HMM8KVS529G1WC4XXZR41VBF",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "update.home_assistant_operating_system_update",
				"state": "off",
				"attributes": {
					"auto_update": false,
					"installed_version": "11.4",
					"in_progress": false,
					"latest_version": "11.4",
					"release_summary": null,
					"release_url": "https://github.com/home-assistant/operating-system/releases/tag/11.4",
					"skipped_version": null,
					"title": "Home Assistant Operating System",
					"entity_picture": "https://brands.home-assistant.io/homeassistant/icon.png",
					"friendly_name": "Home Assistant Operating System Update",
					"supported_features": 3
				},
				"last_changed": "2024-01-20T20:22:24.806667+00:00",
				"last_updated": "2024-01-20T20:22:24.806667+00:00",
				"context": {
					"id": "01HMM8KVS66BHWP5WBRGQT5QGM",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sun.sun",
				"state": "above_horizon",
				"attributes": {
					"next_dawn": "2024-01-28T14:11:20.795629+00:00",
					"next_dusk": "2024-01-28T01:19:59.531548+00:00",
					"next_midnight": "2024-01-28T07:45:59+00:00",
					"next_noon": "2024-01-27T19:45:37+00:00",
					"next_rising": "2024-01-28T14:39:32.903950+00:00",
					"next_setting": "2024-01-28T00:51:42.929110+00:00",
					"elevation": 23.75,
					"azimuth": 140.13,
					"rising": true,
					"friendly_name": "Sun"
				},
				"last_changed": "2024-01-27T14:40:14.663568+00:00",
				"last_updated": "2024-01-27T17:12:53.155029+00:00",
				"context": {
					"id": "01HN5YHVN307FSAHR1VDRWFRJ2",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "zone.home",
				"state": "0",
				"attributes": {
					"latitude": 37.2019399,
					"longitude": -113.270673,
					"radius": 100,
					"passive": false,
					"persons": [],
					"editable": true,
					"icon": "mdi:home",
					"friendly_name": "Home"
				},
				"last_changed": "2024-01-20T20:22:25.609878+00:00",
				"last_updated": "2024-01-20T20:22:25.609878+00:00",
				"context": {
					"id": "01HMM8KWJ9S5B586JJ7VFH10SX",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.sun_next_dawn",
				"state": "2024-01-28T14:11:20+00:00",
				"attributes": {
					"device_class": "timestamp",
					"icon": "mdi:sun-clock",
					"friendly_name": "Sun Next dawn"
				},
				"last_changed": "2024-01-27T14:11:58.834255+00:00",
				"last_updated": "2024-01-27T14:11:58.834255+00:00",
				"context": {
					"id": "01HN5M6KQJRB8GBE8AJYXJXWYG",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.sun_next_dusk",
				"state": "2024-01-28T01:19:59+00:00",
				"attributes": {
					"device_class": "timestamp",
					"icon": "mdi:sun-clock",
					"friendly_name": "Sun Next dusk"
				},
				"last_changed": "2024-01-27T01:18:57.561696+00:00",
				"last_updated": "2024-01-27T01:18:57.561696+00:00",
				"context": {
					"id": "01HN47Z5GSHFQJSJCQ2EQR7TS8",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.sun_next_midnight",
				"state": "2024-01-28T07:45:59+00:00",
				"attributes": {
					"device_class": "timestamp",
					"icon": "mdi:sun-clock",
					"friendly_name": "Sun Next midnight"
				},
				"last_changed": "2024-01-27T07:45:47.007462+00:00",
				"last_updated": "2024-01-27T07:45:47.007462+00:00",
				"context": {
					"id": "01HN4Y3EZZQFERPHGTMD35ZR9Z",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.sun_next_noon",
				"state": "2024-01-27T19:45:37+00:00",
				"attributes": {
					"device_class": "timestamp",
					"icon": "mdi:sun-clock",
					"friendly_name": "Sun Next noon"
				},
				"last_changed": "2024-01-26T19:45:24.009011+00:00",
				"last_updated": "2024-01-26T19:45:24.009011+00:00",
				"context": {
					"id": "01HN3MWD193ZVBXWX6NH0N0WMY",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.sun_next_rising",
				"state": "2024-01-28T14:39:32+00:00",
				"attributes": {
					"device_class": "timestamp",
					"icon": "mdi:sun-clock",
					"friendly_name": "Sun Next rising"
				},
				"last_changed": "2024-01-27T14:40:14.664565+00:00",
				"last_updated": "2024-01-27T14:40:14.664565+00:00",
				"context": {
					"id": "01HN5NTBT8P052SF43NJRBZTY4",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.sun_next_setting",
				"state": "2024-01-28T00:51:42+00:00",
				"attributes": {
					"device_class": "timestamp",
					"icon": "mdi:sun-clock",
					"friendly_name": "Sun Next setting"
				},
				"last_changed": "2024-01-27T00:50:37.275201+00:00",
				"last_updated": "2024-01-27T00:50:37.275201+00:00",
				"context": {
					"id": "01HN46B92V6ANP245YP49WZ60C",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "tts.google_en_com",
				"state": "unknown",
				"attributes": {
					"friendly_name": "Google en com"
				},
				"last_changed": "2024-01-20T20:22:29.304656+00:00",
				"last_updated": "2024-01-20T20:22:29.304656+00:00",
				"context": {
					"id": "01HMM8M05RC5NZQ0109CDRPSYD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "todo.shopping_list",
				"state": "0",
				"attributes": {
					"icon": "mdi:cart",
					"friendly_name": "Shopping List",
					"supported_features": 15
				},
				"last_changed": "2024-01-20T20:22:29.382528+00:00",
				"last_updated": "2024-01-20T20:22:29.382528+00:00",
				"context": {
					"id": "01HMM8M086E4V0TPBVTSFTD8QX",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "weather.forecast_home",
				"state": "sunny",
				"attributes": {
					"temperature": 48,
					"dew_point": 34,
					"temperature_unit": "°F",
					"humidity": 59,
					"cloud_coverage": 6.2,
					"pressure": 30.5,
					"pressure_unit": "inHg",
					"wind_bearing": 83.6,
					"wind_speed": 7.39,
					"wind_speed_unit": "mph",
					"visibility_unit": "mi",
					"precipitation_unit": "in",
					"forecast": [
						{
							"condition": "sunny",
							"datetime": "2024-01-27T19:00:00+00:00",
							"wind_bearing": 217.4,
							"temperature": 65,
							"templow": 48,
							"wind_speed": 8.08,
							"precipitation": 0.0,
							"humidity": 51
						},
						{
							"condition": "partlycloudy",
							"datetime": "2024-01-28T19:00:00+00:00",
							"wind_bearing": 215.4,
							"temperature": 71,
							"templow": 47,
							"wind_speed": 8.26,
							"precipitation": 0.0,
							"humidity": 51
						},
						{
							"condition": "partlycloudy",
							"datetime": "2024-01-29T19:00:00+00:00",
							"wind_bearing": 222.4,
							"temperature": 71,
							"templow": 48,
							"wind_speed": 7.83,
							"precipitation": 0.0,
							"humidity": 54
						},
						{
							"condition": "sunny",
							"datetime": "2024-01-30T19:00:00+00:00",
							"wind_bearing": 71.0,
							"temperature": 69,
							"templow": 49,
							"wind_speed": 7.83,
							"precipitation": 0.0,
							"humidity": 52
						},
						{
							"condition": "cloudy",
							"datetime": "2024-01-31T19:00:00+00:00",
							"wind_bearing": 72.4,
							"temperature": 71,
							"templow": 49,
							"wind_speed": 6.96,
							"precipitation": 0.0,
							"humidity": 53
						},
						{
							"condition": "rainy",
							"datetime": "2024-02-01T19:00:00+00:00",
							"wind_bearing": 54.4,
							"temperature": 61,
							"templow": 54,
							"wind_speed": 6.46,
							"precipitation": 0.59,
							"humidity": 81
						}
					],
					"attribution": "Weather forecast from met.no, delivered by the Norwegian Meteorological Institute.",
					"friendly_name": "Forecast Home",
					"supported_features": 3
				},
				"last_changed": "2024-01-27T15:41:07.029988+00:00",
				"last_updated": "2024-01-27T16:43:08.027124+00:00",
				"context": {
					"id": "01HN5WVCBVTBMK0Q94TN2PSZ2N",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "binary_sensor.rpi_power_status",
				"state": "off",
				"attributes": {
					"device_class": "problem",
					"icon": "mdi:raspberry-pi",
					"friendly_name": "RPi Power status"
				},
				"last_changed": "2024-01-20T20:22:29.412323+00:00",
				"last_updated": "2024-01-20T20:22:29.412323+00:00",
				"context": {
					"id": "01HMM8M094JNEQBV0F7AQJH3HE",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47e243",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 56,
					"temperature": 72,
					"fan_mode": "low",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "09B Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T12:21:13.157181+00:00",
				"last_updated": "2024-01-27T12:21:13.157181+00:00",
				"context": {
					"id": "01HN5DVST5BVBQJMWA6W6GRM3J",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e243_panel_light",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "09B Gree Panel light"
				},
				"last_changed": "2024-01-27T12:21:13.157511+00:00",
				"last_updated": "2024-01-27T12:21:13.157511+00:00",
				"context": {
					"id": "01HN5DVST5Y96XNK2NSXVX9ZPE",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e243_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "09B Gree Quiet"
				},
				"last_changed": "2024-01-27T12:21:13.157696+00:00",
				"last_updated": "2024-01-27T12:21:13.157696+00:00",
				"context": {
					"id": "01HN5DVST50BN6J8RDTWX03QGQ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e243_fresh_air",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "09B Gree Fresh air"
				},
				"last_changed": "2024-01-27T12:21:13.157868+00:00",
				"last_updated": "2024-01-27T12:21:13.157868+00:00",
				"context": {
					"id": "01HN5DVST5XY8JM8JWCMSFEWBF",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e243_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "09B Gree XFan"
				},
				"last_changed": "2024-01-27T12:21:13.158029+00:00",
				"last_updated": "2024-01-27T12:21:13.158029+00:00",
				"context": {
					"id": "01HN5DVST67S3AGKTV3TEQSZDH",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47f179",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 49,
					"temperature": 70,
					"fan_mode": "auto",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "06A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T15:25:16.144619+00:00",
				"last_updated": "2024-01-27T17:00:16.040851+00:00",
				"context": {
					"id": "01HN5XTR98H80SDTC93FQYKJ9H",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f179_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "06A Gree Panel light"
				},
				"last_changed": "2024-01-27T15:25:16.144896+00:00",
				"last_updated": "2024-01-27T15:25:16.144896+00:00",
				"context": {
					"id": "01HN5RCSZGGP6TG7KE2A1YE9TF",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f179_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "06A Gree Quiet"
				},
				"last_changed": "2024-01-27T15:25:16.145029+00:00",
				"last_updated": "2024-01-27T15:25:16.145029+00:00",
				"context": {
					"id": "01HN5RCSZH0Z2BK8YVE1DTCZXZ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f179_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "06A Gree Fresh air"
				},
				"last_changed": "2024-01-27T15:25:16.145216+00:00",
				"last_updated": "2024-01-27T15:25:16.145216+00:00",
				"context": {
					"id": "01HN5RCSZHZPHXDQBDMAACQ33D",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f179_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "06A Gree XFan"
				},
				"last_changed": "2024-01-27T15:25:16.145345+00:00",
				"last_updated": "2024-01-27T15:25:16.145345+00:00",
				"context": {
					"id": "01HN5RCSZHG99AAH27AS56G1TE",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47f10a",
				"state": "unavailable",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"friendly_name": "06B Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:02:26.219714+00:00",
				"last_updated": "2024-01-27T17:02:26.219714+00:00",
				"context": {
					"id": "01HN5XYQDBTDYNFMZ2QN72YN9F",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f10a_panel_light",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "06B Gree Panel light"
				},
				"last_changed": "2024-01-27T17:02:26.220238+00:00",
				"last_updated": "2024-01-27T17:02:26.220238+00:00",
				"context": {
					"id": "01HN5XYQDC8AQSG4D5FEBQC1FC",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f10a_quiet",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "06B Gree Quiet"
				},
				"last_changed": "2024-01-27T17:02:26.220479+00:00",
				"last_updated": "2024-01-27T17:02:26.220479+00:00",
				"context": {
					"id": "01HN5XYQDCZVG089S0W8TV589S",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f10a_fresh_air",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "06B Gree Fresh air"
				},
				"last_changed": "2024-01-27T17:02:26.220694+00:00",
				"last_updated": "2024-01-27T17:02:26.220694+00:00",
				"context": {
					"id": "01HN5XYQDCGB1D84WGJ9S1KKDT",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f10a_xfan",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "06B Gree XFan"
				},
				"last_changed": "2024-01-27T17:02:26.220904+00:00",
				"last_updated": "2024-01-27T17:02:26.220904+00:00",
				"context": {
					"id": "01HN5XYQDC8R49GEB2KYW85C6P",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e461b67",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 54,
					"temperature": 67,
					"fan_mode": "high",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "10A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-26T22:02:44.259249+00:00",
				"last_updated": "2024-01-27T17:07:44.167467+00:00",
				"context": {
					"id": "01HN5Y8DX7FF9RQSDSCR7Y7FB3",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b67_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "10A Gree Panel light"
				},
				"last_changed": "2024-01-26T22:02:44.259582+00:00",
				"last_updated": "2024-01-26T22:02:44.259582+00:00",
				"context": {
					"id": "01HN3WQW537HAWWCQVC3591Q01",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b67_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "10A Gree Quiet"
				},
				"last_changed": "2024-01-26T22:02:44.259736+00:00",
				"last_updated": "2024-01-26T22:02:44.259736+00:00",
				"context": {
					"id": "01HN3WQW538T4CBQ7AZ988ANFR",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b67_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "10A Gree Fresh air"
				},
				"last_changed": "2024-01-26T22:02:44.259878+00:00",
				"last_updated": "2024-01-26T22:02:44.259878+00:00",
				"context": {
					"id": "01HN3WQW53V34Y7Q17AG0V032R",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b67_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "10A Gree XFan"
				},
				"last_changed": "2024-01-26T22:02:44.260017+00:00",
				"last_updated": "2024-01-26T22:02:44.260017+00:00",
				"context": {
					"id": "01HN3WQW54QSHR3WXBZ038B5G7",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47e1fd",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 54,
					"temperature": 70,
					"fan_mode": "medium",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "09A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-26T20:40:23.045316+00:00",
				"last_updated": "2024-01-27T17:11:22.943704+00:00",
				"context": {
					"id": "01HN5YF3HZKMBQBN6RNS7A31RM",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1fd_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "09A Gree Panel light"
				},
				"last_changed": "2024-01-26T20:40:23.045626+00:00",
				"last_updated": "2024-01-26T20:40:23.045626+00:00",
				"context": {
					"id": "01HN3R12R55EVHGJQCBBW5F8GA",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1fd_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "09A Gree Quiet"
				},
				"last_changed": "2024-01-26T20:40:23.046243+00:00",
				"last_updated": "2024-01-26T20:40:23.046243+00:00",
				"context": {
					"id": "01HN3R12R6A1V5QVAEVQ8DNGNQ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1fd_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "09A Gree Fresh air"
				},
				"last_changed": "2024-01-26T20:40:23.046411+00:00",
				"last_updated": "2024-01-26T20:40:23.046411+00:00",
				"context": {
					"id": "01HN3R12R6WX0MSRMS3EDZF0SG",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1fd_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "09A Gree XFan"
				},
				"last_changed": "2024-01-26T20:40:23.046558+00:00",
				"last_updated": "2024-01-26T20:40:23.046558+00:00",
				"context": {
					"id": "01HN3R12R6J126C2MM3W3Y4G1F",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47d0ae",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 46,
					"temperature": 71,
					"fan_mode": "auto",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "05A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T12:02:05.191209+00:00",
				"last_updated": "2024-01-27T17:01:05.192471+00:00",
				"context": {
					"id": "01HN5XW8983BVRMWTXRCJYNV7A",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ae_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "05A Gree Panel light"
				},
				"last_changed": "2024-01-27T12:02:05.191483+00:00",
				"last_updated": "2024-01-27T12:02:05.191483+00:00",
				"context": {
					"id": "01HN5CRRR71JPAB04WXR7T5PR1",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ae_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "05A Gree Quiet"
				},
				"last_changed": "2024-01-27T12:02:05.191617+00:00",
				"last_updated": "2024-01-27T12:02:05.191617+00:00",
				"context": {
					"id": "01HN5CRRR70SXSWK8MW5JXQZB4",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ae_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "05A Gree Fresh air"
				},
				"last_changed": "2024-01-27T12:02:05.191745+00:00",
				"last_updated": "2024-01-27T12:02:05.191745+00:00",
				"context": {
					"id": "01HN5CRRR79Q1HAWK6N94FDCPQ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ae_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "05A Gree XFan"
				},
				"last_changed": "2024-01-27T12:02:05.191867+00:00",
				"last_updated": "2024-01-27T12:02:05.191867+00:00",
				"context": {
					"id": "01HN5CRRR7EJ2WC968X7EVPJAA",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47f056",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 52,
					"temperature": 70,
					"fan_mode": "high",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "11A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T05:01:49.908322+00:00",
				"last_updated": "2024-01-27T13:02:54.911473+00:00",
				"context": {
					"id": "01HN5G84XZ84YGFNJM16E1VDMS",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f056_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "11A Gree Panel light"
				},
				"last_changed": "2024-01-27T05:01:49.908589+00:00",
				"last_updated": "2024-01-27T05:01:49.908589+00:00",
				"context": {
					"id": "01HN4MQ8EMF9CQEVVQ52JT0QW6",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f056_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "11A Gree Quiet"
				},
				"last_changed": "2024-01-27T05:01:49.908709+00:00",
				"last_updated": "2024-01-27T05:01:49.908709+00:00",
				"context": {
					"id": "01HN4MQ8EM9B0XHEP4AD64ZH42",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f056_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "11A Gree Fresh air"
				},
				"last_changed": "2024-01-27T05:01:49.908822+00:00",
				"last_updated": "2024-01-27T05:01:49.908822+00:00",
				"context": {
					"id": "01HN4MQ8EM8BSKDJE7DZ5K8XAF",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f056_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "11A Gree XFan"
				},
				"last_changed": "2024-01-27T05:01:49.908933+00:00",
				"last_updated": "2024-01-27T05:01:49.908933+00:00",
				"context": {
					"id": "01HN4MQ8EMK8C1JY8SBD3H1EXC",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47e1b4",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 51,
					"temperature": 68,
					"fan_mode": "auto",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "12A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:16:05.078825+00:00",
				"last_updated": "2024-01-27T17:16:05.078825+00:00",
				"context": {
					"id": "01HN5YQQ2PE86E7FN8QT8XY0B4",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1b4_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "12A Gree Panel light"
				},
				"last_changed": "2024-01-27T17:16:05.079430+00:00",
				"last_updated": "2024-01-27T17:16:05.079430+00:00",
				"context": {
					"id": "01HN5YQQ2Q6HPVXZ6T0VHC8EFH",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1b4_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "12A Gree Quiet"
				},
				"last_changed": "2024-01-27T17:16:05.079659+00:00",
				"last_updated": "2024-01-27T17:16:05.079659+00:00",
				"context": {
					"id": "01HN5YQQ2Q5MWAF725114NPTAA",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1b4_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "12A Gree Fresh air"
				},
				"last_changed": "2024-01-27T17:16:05.079869+00:00",
				"last_updated": "2024-01-27T17:16:05.079869+00:00",
				"context": {
					"id": "01HN5YQQ2Q2EC6AD41XJYBEGHQ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1b4_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "12A Gree XFan"
				},
				"last_changed": "2024-01-27T17:16:05.080080+00:00",
				"last_updated": "2024-01-27T17:16:05.080080+00:00",
				"context": {
					"id": "01HN5YQQ2RS20WBNB03FCVKVAH",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01b_inside_temperature",
				"state": "51.8",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "01B Daikin Inside temperature"
				},
				"last_changed": "2024-01-27T17:14:10.699956+00:00",
				"last_updated": "2024-01-27T17:14:10.699956+00:00",
				"context": {
					"id": "01HN5YM7CBVSRNAR1E0W74YDRC",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01b_outside_temperature",
				"state": "39.2",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "01B Daikin Outside temperature"
				},
				"last_changed": "2024-01-27T16:46:40.518915+00:00",
				"last_updated": "2024-01-27T16:46:40.518915+00:00",
				"context": {
					"id": "01HN5X1VW6FVQQXW8QZD4YMZ19",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01b_compressor_estimated_power_consumption",
				"state": "0.0",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "kW",
					"device_class": "power",
					"friendly_name": "01B Daikin Compressor estimated power consumption"
				},
				"last_changed": "2024-01-26T07:14:06.692975+00:00",
				"last_updated": "2024-01-26T07:14:06.692975+00:00",
				"context": {
					"id": "01HN29WR74B657EJYK3P15HYX7",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01b_energy_consumption",
				"state": "0",
				"attributes": {
					"state_class": "total_increasing",
					"unit_of_measurement": "kWh",
					"device_class": "energy",
					"friendly_name": "01B Daikin Energy consumption"
				},
				"last_changed": "2024-01-20T20:39:52.249650+00:00",
				"last_updated": "2024-01-25T16:16:27.858546+00:00",
				"context": {
					"id": "01HN0PH3PJJ7PEK43P1MCRNZAV",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.01b_streamer",
				"state": "off",
				"attributes": {
					"icon": "mdi:air-filter",
					"friendly_name": "01B Daikin Streamer"
				},
				"last_changed": "2024-01-20T20:39:52.253349+00:00",
				"last_updated": "2024-01-25T16:16:27.859084+00:00",
				"context": {
					"id": "01HN0PH3PKDZS664WKVK5BXFM9",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.01b_none",
				"state": "off",
				"attributes": {
					"icon": "mdi:power",
					"friendly_name": "01B Daikin None"
				},
				"last_changed": "2024-01-20T20:53:22.312603+00:00",
				"last_updated": "2024-01-25T16:16:27.859494+00:00",
				"context": {
					"id": "01HN0PH3PKQCSCXZFDV446PTN1",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.01b",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"fan_only",
						"dry",
						"cool",
						"heat",
						"heat_cool",
						"off"
					],
					"min_temp": 45,
					"max_temp": 95,
					"target_temp_step": 1,
					"preset_modes": [
						"none",
						"away",
						"eco",
						"boost"
					],
					"current_temperature": 52,
					"temperature": 70,
					"hvac_action": "off",
					"preset_mode": "none",
					"friendly_name": "01B Daikin",
					"supported_features": 17
				},
				"last_changed": "2024-01-20T20:53:19.842370+00:00",
				"last_updated": "2024-01-27T17:15:01.552958+00:00",
				"context": {
					"id": "01HN5YNS1GBYBHN8QW5CQQR16G",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01a_inside_temperature",
				"state": "77.0",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "01A Daikin Inside temperature"
				},
				"last_changed": "2024-01-27T16:56:24.978063+00:00",
				"last_updated": "2024-01-27T16:56:24.978063+00:00",
				"context": {
					"id": "01HN5XKPMJG8B0B2K9JQ2MS42Y",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01a_outside_temperature",
				"state": "42.8",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "01A Daikin Outside temperature"
				},
				"last_changed": "2024-01-27T17:08:25.006550+00:00",
				"last_updated": "2024-01-27T17:08:25.006550+00:00",
				"context": {
					"id": "01HN5Y9NSE0SQMX6VD33SEXNS5",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01a_compressor_estimated_power_consumption",
				"state": "0",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "kW",
					"device_class": "power",
					"friendly_name": "01A Daikin Compressor estimated power consumption"
				},
				"last_changed": "2024-01-27T16:59:24.982998+00:00",
				"last_updated": "2024-01-27T16:59:24.982998+00:00",
				"context": {
					"id": "01HN5XS6DP6ZJZH506GJ8S3PA9",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.01a_energy_consumption",
				"state": "0",
				"attributes": {
					"state_class": "total_increasing",
					"unit_of_measurement": "kWh",
					"device_class": "energy",
					"friendly_name": "01A Daikin Energy consumption"
				},
				"last_changed": "2024-01-20T20:41:06.868840+00:00",
				"last_updated": "2024-01-25T16:16:05.179884+00:00",
				"context": {
					"id": "01HN0PGDHVY6TY3ABBNYHDTK81",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.01a_streamer",
				"state": "off",
				"attributes": {
					"icon": "mdi:air-filter",
					"friendly_name": "01A Daikin Streamer"
				},
				"last_changed": "2024-01-20T20:41:06.871489+00:00",
				"last_updated": "2024-01-25T16:16:05.180176+00:00",
				"context": {
					"id": "01HN0PGDHWQGN10DYM30TNN5K6",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.01a_none",
				"state": "on",
				"attributes": {
					"icon": "mdi:power",
					"friendly_name": "01A Daikin None"
				},
				"last_changed": "2024-01-27T14:45:24.709536+00:00",
				"last_updated": "2024-01-27T14:45:24.709536+00:00",
				"context": {
					"id": "01HN5P3TK5EAQKXDX4PM4GF10E",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.01a",
				"state": "heat_cool",
				"attributes": {
					"hvac_modes": [
						"fan_only",
						"dry",
						"cool",
						"heat",
						"heat_cool",
						"off"
					],
					"min_temp": 45,
					"max_temp": 95,
					"target_temp_step": 1,
					"preset_modes": [
						"none",
						"away",
						"eco",
						"boost"
					],
					"current_temperature": 77,
					"temperature": 77,
					"preset_mode": "none",
					"friendly_name": "01A Daikin",
					"supported_features": 17
				},
				"last_changed": "2024-01-27T14:45:16.292853+00:00",
				"last_updated": "2024-01-27T16:56:16.423417+00:00",
				"context": {
					"id": "01HN5XKE97A3PDFB8NVZMM8Q8Y",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02a_inside_temperature",
				"state": "69.8",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "02A Daikin Inside temperature"
				},
				"last_changed": "2024-01-27T17:15:49.125098+00:00",
				"last_updated": "2024-01-27T17:15:49.125098+00:00",
				"context": {
					"id": "01HN5YQ7G5GB3P22N50H498MAS",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02a_outside_temperature",
				"state": "50.0",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "02A Daikin Outside temperature"
				},
				"last_changed": "2024-01-27T17:09:18.967228+00:00",
				"last_updated": "2024-01-27T17:09:18.967228+00:00",
				"context": {
					"id": "01HN5YBAFQAHKYDJZXJ133Y776",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02a_compressor_estimated_power_consumption",
				"state": "0.75",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "kW",
					"device_class": "power",
					"friendly_name": "02A Daikin Compressor estimated power consumption"
				},
				"last_changed": "2024-01-27T16:34:18.910292+00:00",
				"last_updated": "2024-01-27T16:34:18.910292+00:00",
				"context": {
					"id": "01HN5WB7MYX6QN9X8Q98CGAN8H",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02a_energy_consumption",
				"state": "0",
				"attributes": {
					"state_class": "total_increasing",
					"unit_of_measurement": "kWh",
					"device_class": "energy",
					"friendly_name": "02A Daikin Energy consumption"
				},
				"last_changed": "2024-01-20T20:41:30.354384+00:00",
				"last_updated": "2024-01-25T16:17:31.436417+00:00",
				"context": {
					"id": "01HN0PK1SCA1GXD3D0MWDX9E3M",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.02a_streamer",
				"state": "off",
				"attributes": {
					"icon": "mdi:air-filter",
					"friendly_name": "02A Daikin Streamer"
				},
				"last_changed": "2024-01-20T20:41:30.357554+00:00",
				"last_updated": "2024-01-25T16:17:31.436721+00:00",
				"context": {
					"id": "01HN0PK1SC26RKNKH9Z0RNZFW4",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.02a_none",
				"state": "on",
				"attributes": {
					"icon": "mdi:power",
					"friendly_name": "02A Daikin None"
				},
				"last_changed": "2024-01-27T15:49:18.837125+00:00",
				"last_updated": "2024-01-27T15:49:18.837125+00:00",
				"context": {
					"id": "01HN5SRTVNH4P4YMJV8AB6VXBA",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.02a",
				"state": "heat",
				"attributes": {
					"hvac_modes": [
						"fan_only",
						"dry",
						"cool",
						"heat",
						"heat_cool",
						"off"
					],
					"min_temp": 45,
					"max_temp": 95,
					"target_temp_step": 1,
					"preset_modes": [
						"none",
						"away",
						"eco",
						"boost"
					],
					"current_temperature": 68,
					"temperature": 75,
					"hvac_action": "heating",
					"preset_mode": "none",
					"friendly_name": "02A Daikin",
					"supported_features": 17
				},
				"last_changed": "2024-01-27T15:49:39.256581+00:00",
				"last_updated": "2024-01-27T17:01:39.317541+00:00",
				"context": {
					"id": "01HN5XX9KNXP3SXVG1SYY25X43",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02b_inside_temperature",
				"state": "77.0",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "02B Daikin Inside temperature"
				},
				"last_changed": "2024-01-27T17:03:14.491851+00:00",
				"last_updated": "2024-01-27T17:03:14.491851+00:00",
				"context": {
					"id": "01HN5Y06HVKKT6BMRVTC21HJKD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02b_outside_temperature",
				"state": "51.8",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "02B Daikin Outside temperature"
				},
				"last_changed": "2024-01-27T17:15:44.381141+00:00",
				"last_updated": "2024-01-27T17:15:44.381141+00:00",
				"context": {
					"id": "01HN5YQ2VXMEYZGNNS9WBVJSZN",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02b_compressor_estimated_power_consumption",
				"state": "1.12",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "kW",
					"device_class": "power",
					"friendly_name": "02B Daikin Compressor estimated power consumption"
				},
				"last_changed": "2024-01-27T17:15:44.382008+00:00",
				"last_updated": "2024-01-27T17:15:44.382008+00:00",
				"context": {
					"id": "01HN5YQ2VYA8Q0J7JJDMJGRM65",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.02b_energy_consumption",
				"state": "0",
				"attributes": {
					"state_class": "total_increasing",
					"unit_of_measurement": "kWh",
					"device_class": "energy",
					"friendly_name": "02B Daikin Energy consumption"
				},
				"last_changed": "2024-01-20T20:41:56.084797+00:00",
				"last_updated": "2024-01-25T16:18:00.837923+00:00",
				"context": {
					"id": "01HN0PKYG569S2XCQ81NKP72WW",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.02b_streamer",
				"state": "off",
				"attributes": {
					"icon": "mdi:air-filter",
					"friendly_name": "02B Daikin Streamer"
				},
				"last_changed": "2024-01-20T20:41:56.090342+00:00",
				"last_updated": "2024-01-25T16:18:00.838259+00:00",
				"context": {
					"id": "01HN0PKYG6610TEEZR3D23PESN",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.02b_none",
				"state": "on",
				"attributes": {
					"icon": "mdi:power",
					"friendly_name": "02B Daikin None"
				},
				"last_changed": "2024-01-27T16:07:44.268452+00:00",
				"last_updated": "2024-01-27T16:07:44.268452+00:00",
				"context": {
					"id": "01HN5TTJCCTDN9BW5T2PVP2N7V",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.02b",
				"state": "heat",
				"attributes": {
					"hvac_modes": [
						"fan_only",
						"dry",
						"cool",
						"heat",
						"heat_cool",
						"off"
					],
					"min_temp": 45,
					"max_temp": 95,
					"target_temp_step": 1,
					"preset_modes": [
						"none",
						"away",
						"eco",
						"boost"
					],
					"current_temperature": 77,
					"temperature": 77,
					"hvac_action": "heating",
					"preset_mode": "none",
					"friendly_name": "02B Daikin",
					"supported_features": 17
				},
				"last_changed": "2024-01-27T16:08:05.626064+00:00",
				"last_updated": "2024-01-27T17:04:05.675824+00:00",
				"context": {
					"id": "01HN5Y1RHB1TM6JKKK94R6QDYV",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.03a_inside_temperature",
				"state": "55.4",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "03A Daikin Inside temperature"
				},
				"last_changed": "2024-01-27T12:04:37.485554+00:00",
				"last_updated": "2024-01-27T12:04:37.485554+00:00",
				"context": {
					"id": "01HN5CXDFDGF64C6ZNSTY614CK",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.03a_outside_temperature",
				"state": "39.2",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "°F",
					"device_class": "temperature",
					"friendly_name": "03A Daikin Outside temperature"
				},
				"last_changed": "2024-01-27T17:04:08.090909+00:00",
				"last_updated": "2024-01-27T17:04:08.090909+00:00",
				"context": {
					"id": "01HN5Y1TWTEZ08PC10CZR5KWMD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.03a_compressor_estimated_power_consumption",
				"state": "0.0",
				"attributes": {
					"state_class": "measurement",
					"unit_of_measurement": "kW",
					"device_class": "power",
					"friendly_name": "03A Daikin Compressor estimated power consumption"
				},
				"last_changed": "2024-01-27T07:02:06.504813+00:00",
				"last_updated": "2024-01-27T07:02:06.504813+00:00",
				"context": {
					"id": "01HN4VKFX8WZSYWA8SFCS7P85K",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "sensor.03a_energy_consumption",
				"state": "0",
				"attributes": {
					"state_class": "total_increasing",
					"unit_of_measurement": "kWh",
					"device_class": "energy",
					"friendly_name": "03A Daikin Energy consumption"
				},
				"last_changed": "2024-01-20T20:42:19.200864+00:00",
				"last_updated": "2024-01-25T16:18:30.953860+00:00",
				"context": {
					"id": "01HN0PMVX9W8GTKJYZ0Y8M67DE",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.03a_streamer",
				"state": "off",
				"attributes": {
					"icon": "mdi:air-filter",
					"friendly_name": "03A Daikin Streamer"
				},
				"last_changed": "2024-01-20T20:42:19.206195+00:00",
				"last_updated": "2024-01-25T16:18:30.954148+00:00",
				"context": {
					"id": "01HN0PMVXAB95DB4QR3T1SH4DQ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.03a_none",
				"state": "off",
				"attributes": {
					"icon": "mdi:power",
					"friendly_name": "03A Daikin None"
				},
				"last_changed": "2024-01-20T20:42:19.207365+00:00",
				"last_updated": "2024-01-25T16:18:30.954525+00:00",
				"context": {
					"id": "01HN0PMVXAAD1DZBVK59HAM0WD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.03a",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"fan_only",
						"dry",
						"cool",
						"heat",
						"heat_cool",
						"off"
					],
					"min_temp": 45,
					"max_temp": 95,
					"target_temp_step": 1,
					"preset_modes": [
						"none",
						"away",
						"eco",
						"boost"
					],
					"current_temperature": 55,
					"temperature": 60,
					"hvac_action": "off",
					"preset_mode": "none",
					"friendly_name": "03A Daikin",
					"supported_features": 17
				},
				"last_changed": "2024-01-20T20:42:19.324233+00:00",
				"last_updated": "2024-01-27T12:05:28.127905+00:00",
				"context": {
					"id": "01HN5CYYXZ29D2RNB1Z03TKEED",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e461b53",
				"state": "unavailable",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"friendly_name": "03B Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:15:43.064557+00:00",
				"last_updated": "2024-01-27T17:15:43.064557+00:00",
				"context": {
					"id": "01HN5YQ1JRV2M1EF7K3SVFRD3D",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b53_panel_light",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "03B Gree Panel light"
				},
				"last_changed": "2024-01-27T17:15:43.065077+00:00",
				"last_updated": "2024-01-27T17:15:43.065077+00:00",
				"context": {
					"id": "01HN5YQ1JSXK2SZKA791QQ4FQG",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b53_quiet",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "03B Gree Quiet"
				},
				"last_changed": "2024-01-27T17:15:43.065295+00:00",
				"last_updated": "2024-01-27T17:15:43.065295+00:00",
				"context": {
					"id": "01HN5YQ1JSFMJQ2M6BHW937D8X",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b53_fresh_air",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "03B Gree Fresh air"
				},
				"last_changed": "2024-01-27T17:15:43.065490+00:00",
				"last_updated": "2024-01-27T17:15:43.065490+00:00",
				"context": {
					"id": "01HN5YQ1JSGK26DWVW8BR37VP4",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e461b53_xfan",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "03B Gree XFan"
				},
				"last_changed": "2024-01-27T17:15:43.065672+00:00",
				"last_updated": "2024-01-27T17:15:43.065672+00:00",
				"context": {
					"id": "01HN5YQ1JSFSP34ZQW99S7GVYK",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47f0dd",
				"state": "heat",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 59,
					"temperature": 60,
					"fan_mode": "high",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "04B Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:15:55.947634+00:00",
				"last_updated": "2024-01-27T17:15:55.947634+00:00",
				"context": {
					"id": "01HN5YQE5BSGPT2JNKCXT946YD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f0dd_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "04B Gree Panel light"
				},
				"last_changed": "2024-01-27T17:15:55.948195+00:00",
				"last_updated": "2024-01-27T17:15:55.948195+00:00",
				"context": {
					"id": "01HN5YQE5C88T0A3SEKJVBA3ZG",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f0dd_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "04B Gree Quiet"
				},
				"last_changed": "2024-01-27T17:15:55.948378+00:00",
				"last_updated": "2024-01-27T17:15:55.948378+00:00",
				"context": {
					"id": "01HN5YQE5CPHFES7HN07D5H7R3",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f0dd_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "04B Gree Fresh air"
				},
				"last_changed": "2024-01-27T17:15:55.948543+00:00",
				"last_updated": "2024-01-27T17:15:55.948543+00:00",
				"context": {
					"id": "01HN5YQE5C9F9QKMHKWAN8BSZD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f0dd_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "04B Gree XFan"
				},
				"last_changed": "2024-01-27T17:15:55.948704+00:00",
				"last_updated": "2024-01-27T17:15:55.948704+00:00",
				"context": {
					"id": "01HN5YQE5CN7RP0C2DVJR1F7CM",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47f104",
				"state": "heat",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 50,
					"temperature": 68,
					"fan_mode": "auto",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "04A Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:13:43.274573+00:00",
				"last_updated": "2024-01-27T17:13:43.274573+00:00",
				"context": {
					"id": "01HN5YKCKANYAR9C0K40GT7G54",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f104_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "04A Gree Panel light"
				},
				"last_changed": "2024-01-27T15:33:43.335212+00:00",
				"last_updated": "2024-01-27T15:33:43.335212+00:00",
				"context": {
					"id": "01HN5RW997M7EZABV4D7JP3WX0",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f104_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "04A Gree Quiet"
				},
				"last_changed": "2024-01-27T15:33:43.335369+00:00",
				"last_updated": "2024-01-27T15:33:43.335369+00:00",
				"context": {
					"id": "01HN5RW997WP0F67E7J8JG0RQA",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f104_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "04A Gree Fresh air"
				},
				"last_changed": "2024-01-27T15:33:43.335515+00:00",
				"last_updated": "2024-01-27T15:33:43.335515+00:00",
				"context": {
					"id": "01HN5RW99732DDT5PEW6FFWE06",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47f104_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "04A Gree XFan"
				},
				"last_changed": "2024-01-27T15:33:43.335657+00:00",
				"last_updated": "2024-01-27T15:33:43.335657+00:00",
				"context": {
					"id": "01HN5RW997XFSPB98QET4MJGME",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ac_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "08B Gree Panel light"
				},
				"last_changed": "2024-01-27T17:15:18.992929+00:00",
				"last_updated": "2024-01-27T17:15:18.992929+00:00",
				"context": {
					"id": "01HN5YPA2GC2XE0J6WR5S55YN8",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ac_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "08B Gree Quiet"
				},
				"last_changed": "2024-01-27T17:15:18.993140+00:00",
				"last_updated": "2024-01-27T17:15:18.993140+00:00",
				"context": {
					"id": "01HN5YPA2H5CKT9Z4NT04TQ9KJ",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ac_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "08B Gree Fresh air"
				},
				"last_changed": "2024-01-27T17:15:18.993325+00:00",
				"last_updated": "2024-01-27T17:15:18.993325+00:00",
				"context": {
					"id": "01HN5YPA2HAWP17H3Q2XTNC4M7",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47d0ac_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "08B Gree XFan"
				},
				"last_changed": "2024-01-27T17:15:18.993512+00:00",
				"last_updated": "2024-01-27T17:15:18.993512+00:00",
				"context": {
					"id": "01HN5YPA2HSGW8EHM6EF2EYPPD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47d0ac",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 52,
					"temperature": 72,
					"fan_mode": "high",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "08B Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:15:18.992438+00:00",
				"last_updated": "2024-01-27T17:15:18.992438+00:00",
				"context": {
					"id": "01HN5YPA2GKB3R3TCFYJQAGH8J",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47e22b",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 51,
					"temperature": 71,
					"fan_mode": "low",
					"preset_mode": "none",
					"swing_mode": "off",
					"friendly_name": "05B Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-26T22:00:55.092827+00:00",
				"last_updated": "2024-01-27T17:07:54.995195+00:00",
				"context": {
					"id": "01HN5Y8RFKPTDSH7JGNBPFD720",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e22b_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "05B Gree Panel light"
				},
				"last_changed": "2024-01-26T22:00:55.093247+00:00",
				"last_updated": "2024-01-26T22:00:55.093247+00:00",
				"context": {
					"id": "01HN3WMHHNY8HY8D9MXRPA4FN9",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e22b_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "05B Gree Quiet"
				},
				"last_changed": "2024-01-26T22:00:55.093424+00:00",
				"last_updated": "2024-01-26T22:00:55.093424+00:00",
				"context": {
					"id": "01HN3WMHHN240JDYKJF1DCABAP",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e22b_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "05B Gree Fresh air"
				},
				"last_changed": "2024-01-26T22:00:55.093590+00:00",
				"last_updated": "2024-01-26T22:00:55.093590+00:00",
				"context": {
					"id": "01HN3WMHHNG65AH1C08QM2AWAF",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e22b_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "05B Gree XFan"
				},
				"last_changed": "2024-01-26T22:00:55.093752+00:00",
				"last_updated": "2024-01-26T22:00:55.093752+00:00",
				"context": {
					"id": "01HN3WMHHNZEC3TJT1RXAW9GEW",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47e1f8",
				"state": "off",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"current_temperature": 54,
					"temperature": 70,
					"fan_mode": "medium",
					"preset_mode": "none",
					"swing_mode": "vertical",
					"friendly_name": "07A? Gree",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:15:33.906823+00:00",
				"last_updated": "2024-01-27T17:15:33.906823+00:00",
				"context": {
					"id": "01HN5YPRMJ4YJ0666PBD88VTBM",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1f8_panel_light",
				"state": "on",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "07A? Gree Panel light"
				},
				"last_changed": "2024-01-27T17:15:33.907448+00:00",
				"last_updated": "2024-01-27T17:15:33.907448+00:00",
				"context": {
					"id": "01HN5YPRMKG65B8VWJYCB6CQAG",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1f8_quiet",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "07A? Gree Quiet"
				},
				"last_changed": "2024-01-27T17:15:33.907717+00:00",
				"last_updated": "2024-01-27T17:15:33.907717+00:00",
				"context": {
					"id": "01HN5YPRMKQR7WY6C1ZPW71K8S",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1f8_fresh_air",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "07A? Gree Fresh air"
				},
				"last_changed": "2024-01-27T17:15:33.907964+00:00",
				"last_updated": "2024-01-27T17:15:33.907964+00:00",
				"context": {
					"id": "01HN5YPRMKXR6VF68WJA397CPF",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e1f8_xfan",
				"state": "off",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "07A? Gree XFan"
				},
				"last_changed": "2024-01-27T17:15:33.908208+00:00",
				"last_updated": "2024-01-27T17:15:33.908208+00:00",
				"context": {
					"id": "01HN5YPRMMFAEJEA5ZMDT3SZDA",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "climate.1e47e21c",
				"state": "unavailable",
				"attributes": {
					"hvac_modes": [
						"auto",
						"cool",
						"dry",
						"fan_only",
						"heat",
						"off"
					],
					"min_temp": 46,
					"max_temp": 86,
					"target_temp_step": 1,
					"fan_modes": [
						"auto",
						"low",
						"medium low",
						"medium",
						"medium high",
						"high"
					],
					"preset_modes": [
						"eco",
						"away",
						"boost",
						"none",
						"sleep"
					],
					"swing_modes": [
						"off",
						"vertical",
						"horizontal",
						"both"
					],
					"friendly_name": "Gree 47:E2:1C",
					"supported_features": 57
				},
				"last_changed": "2024-01-27T17:15:13.328031+00:00",
				"last_updated": "2024-01-27T17:15:13.328031+00:00",
				"context": {
					"id": "01HN5YP4HGN6NG8FYZZ18B7FQD",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e21c_panel_light",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"icon": "mdi:lightbulb",
					"friendly_name": "Gree 47:E2:1C Panel light"
				},
				"last_changed": "2024-01-27T17:15:13.328555+00:00",
				"last_updated": "2024-01-27T17:15:13.328555+00:00",
				"context": {
					"id": "01HN5YP4HGKT6JWVJJRHNPWYY5",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e21c_quiet",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "Gree 47:E2:1C Quiet"
				},
				"last_changed": "2024-01-27T17:15:13.328794+00:00",
				"last_updated": "2024-01-27T17:15:13.328794+00:00",
				"context": {
					"id": "01HN5YP4HGZ105DGMX2G6S13CH",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e21c_fresh_air",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "Gree 47:E2:1C Fresh air"
				},
				"last_changed": "2024-01-27T17:15:13.329007+00:00",
				"last_updated": "2024-01-27T17:15:13.329007+00:00",
				"context": {
					"id": "01HN5YP4HH0YXF9JAC8N2R996Z",
					"parent_id": null,
					"user_id": null
				}
			},
			{
				"entity_id": "switch.1e47e21c_xfan",
				"state": "unavailable",
				"attributes": {
					"device_class": "switch",
					"friendly_name": "Gree 47:E2:1C XFan"
				},
				"last_changed": "2024-01-27T17:15:13.329214+00:00",
				"last_updated": "2024-01-27T17:15:13.329214+00:00",
				"context": {
					"id": "01HN5YP4HHKBG6JPX812ZE2JSE",
					"parent_id": null,
					"user_id": null
				}
			}
		]
	`)
}
