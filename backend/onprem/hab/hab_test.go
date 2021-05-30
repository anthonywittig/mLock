package hab

import (
	"context"
	"fmt"
	"mlock/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	skipIntegTests = true
)

func TestProcessCommand_listThings_integration(t *testing.T) {
	if skipIntegTests {
		return
	}

	ctx := context.Background()

	resp, err := ProcessCommand(ctx, shared.HabCommandListThings("hi"))
	assert.Nil(t, err)

	fmt.Printf("%s\n", resp.Description)
	assert.Equal(t, resp.Description, "fail!")

	/*
		expected += `[
			{
			  "statusInfo": {
				"status": "ONLINE",
				"statusDetail": "NONE"
			  },
			  "editable": true,
			  "label": "Murset\\u0027s schlage lock",
			  "bridgeUID": "zwave:serial_zstick:94acf23fa3",
			  "configuration": {
				"usercode_code_27": "",
				"config_10_1": 3,
				"usercode_code_28": "",
				"usercode_code_29": "",
				"doorlock_timeout": 0,
				"group_1": [
				  "controller"
				],
				"config_12_4": 4,
				"config_16_1": 4,
				"usercode_code_23": "",
				"config_14_4": -1,
				"usercode_code_24": "",
				"config_18_1": 6,
				"usercode_code_25": "",
				"usercode_code_26": "",
				"usercode_code_20": "",
				"usercode_code_21": "",
				"usercode_code_22": "",
				"config_7_1": 0,
				"config_9_1": 3,
				"config_3_1": -1,
				"usercode_code_30": "",
				"config_5_1": 0,
				"config_11_1": 0,
				"config_17_4": -1,
				"config_13_4": -1,
				"config_15_1": 0,
				"usercode_code_16": "",
				"usercode_code_17": "",
				"usercode_code_18": "",
				"usercode_code_19": "",
				"usercode_code_7": "",
				"usercode_code_8": "",
				"usercode_code_9": "",
				"usercode_code_3": "",
				"usercode_code_4": "",
				"usercode_code_5": "",
				"usercode_code_6": "",
				"usercode_code_12": "",
				"config_8_1": 3,
				"usercode_code_13": "",
				"config_6_4": 50331648,
				"usercode_code_1": "37 35 33 34 0A 0D",
				"usercode_code_14": "",
				"usercode_code_2": "33 38 31 32 0A 0D",
				"usercode_code_15": "",
				"config_4_1": 0,
				"usercode_code_10": "",
				"node_id": 2,
				"usercode_code_11": ""
			  },
			  "properties": {
				"zwave_class_basic": "BASIC_TYPE_ROUTING_SLAVE",
				"zwave_class_generic": "GENERIC_TYPE_ENTRY_CONTROL",
				"zwave_frequent": "true",
				"zwave_neighbours": "1",
				"modelId": "BE468ZP",
				"zwave_version": "3.3",
				"zwave_listening": "false",
				"zwave_plus_devicetype": "NODE_TYPE_ZWAVEPLUS_NODE",
				"manufacturerId": "003B",
				"manufacturerRef": "0001:0468",
				"dbReference": "1223",
				"zwave_deviceid": "1128",
				"zwave_nodeid": "2",
				"zwave_lastheal": "2021-05-29T08:35:32Z",
				"vendor": "Allegion",
				"defaultAssociations": "1",
				"zwave_routing": "true",
				"zwave_plus_roletype": "ROLE_TYPE_SLAVE_SLEEPING_LISTENING",
				"zwave_beaming": "true",
				"zwave_secure": "true",
				"zwave_class_specific": "SPECIFIC_TYPE_SECURE_KEYPAD_DOOR_LOCK",
				"zwave_manufacturer": "59",
				"zwave_devicetype": "1"
			  },
			  "UID": "zwave:device:94acf23fa3:node2",
			  "thingTypeUID": "zwave:schlage_be468zp_00_000",
			  "channels": [
				{
				  "linkedItems": [],
				  "uid": "zwave:device:94acf23fa3:node2:lock_door",
				  "id": "lock_door",
				  "channelTypeUID": "zwave:lock_door",
				  "itemType": "Switch",
				  "kind": "STATE",
				  "label": "Door Lock",
				  "description": "Lock and unlock the door",
				  "defaultTags": [],
				  "properties": {
					"binding:*:OnOffType": "COMMAND_CLASS_DOOR_LOCK"
				  },
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:device:94acf23fa3:node2:alarm_access",
				  "id": "alarm_access",
				  "channelTypeUID": "zwave:alarm_access",
				  "itemType": "Switch",
				  "kind": "STATE",
				  "label": "Alarm (access)",
				  "description": "Indicates if the access control alarm is triggered",
				  "defaultTags": [],
				  "properties": {
					"binding:*:OnOffType": "COMMAND_CLASS_ALARM;type\\u003dACCESS_CONTROL"
				  },
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:device:94acf23fa3:node2:alarm_power",
				  "id": "alarm_power",
				  "channelTypeUID": "zwave:alarm_power",
				  "itemType": "Switch",
				  "kind": "STATE",
				  "label": "Alarm (power)",
				  "description": "Indicates if a power alarm is triggered",
				  "defaultTags": [],
				  "properties": {
					"binding:*:OnOffType": "COMMAND_CLASS_ALARM;type\\u003dPOWER_MANAGEMENT"
				  },
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:device:94acf23fa3:node2:alarm_system",
				  "id": "alarm_system",
				  "channelTypeUID": "zwave:alarm_system",
				  "itemType": "Switch",
				  "kind": "STATE",
				  "label": "Alarm (system)",
				  "description": "Indicates if a system alarm is triggered",
				  "defaultTags": [],
				  "properties": {
					"binding:*:OnOffType": "COMMAND_CLASS_ALARM;type\\u003dSYSTEM"
				  },
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:device:94acf23fa3:node2:alarm_raw",
				  "id": "alarm_raw",
				  "channelTypeUID": "zwave:alarm_raw",
				  "itemType": "String",
				  "kind": "STATE",
				  "label": "Alarm (raw)",
				  "description": "Provides alarm parameters as json string",
				  "defaultTags": [],
				  "properties": {
					"binding:*:StringType": "COMMAND_CLASS_ALARM"
				  },
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:device:94acf23fa3:node2:battery-level",
				  "id": "battery-level",
				  "channelTypeUID": "system:battery-level",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "Battery Level",
				  "defaultTags": [],
				  "properties": {
					"binding:*:PercentType": "COMMAND_CLASS_BATTERY"
				  },
				  "configuration": {}
				}
			  ]
			},
			{
			  "statusInfo": {
				"status": "ONLINE",
				"statusDetail": "NONE"
			  },
			  "editable": true,
			  "label": "Z-Wave Serial Controller",
			  "configuration": {
				"controller_softreset": false,
				"security_networkkey": "EC 41 B4 15 33 8E B1 A5 C2 DA 7C 80 96 D6 FB 5B",
				"security_inclusionmode": 0,
				"controller_sisnode": 1,
				"controller_sync": false,
				"port": "/dev/ttyACM0",
				"controller_master": true,
				"inclusion_mode": 2,
				"controller_wakeupperiod": 3600,
				"heal_time": 2,
				"controller_exclude": false,
				"controller_inclusiontimeout": 30,
				"controller_hardreset": false
			  },
			  "properties": {
				"zwave_nodeid": "1",
				"zwave_neighbours": "2"
			  },
			  "UID": "zwave:serial_zstick:94acf23fa3",
			  "thingTypeUID": "zwave:serial_zstick",
			  "channels": [
				{
				  "linkedItems": [],
				  "uid": "zwave:serial_zstick:94acf23fa3:serial_sof",
				  "id": "serial_sof",
				  "channelTypeUID": "zwave:serial_sof",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "Start Frames",
				  "description": "Counter tracking the number of SOF bytes received",
				  "defaultTags": [],
				  "properties": {},
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:serial_zstick:94acf23fa3:serial_ack",
				  "id": "serial_ack",
				  "channelTypeUID": "zwave:serial_ack",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "Frames Acknowledged",
				  "description": "Counter tracking the number of frames acknowledged by the controller",
				  "defaultTags": [],
				  "properties": {},
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:serial_zstick:94acf23fa3:serial_nak",
				  "id": "serial_nak",
				  "channelTypeUID": "zwave:serial_nak",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "Frames Rejected",
				  "description": "Counter tracking the number of frames rejected by the controller",
				  "defaultTags": [],
				  "properties": {},
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:serial_zstick:94acf23fa3:serial_can",
				  "id": "serial_can",
				  "channelTypeUID": "zwave:serial_can",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "Frames Cancelled",
				  "description": "Counter tracking the number of frames cancelled by the controller",
				  "defaultTags": [],
				  "properties": {},
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:serial_zstick:94acf23fa3:serial_oof",
				  "id": "serial_oof",
				  "channelTypeUID": "zwave:serial_oof",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "OOF Bytes Received",
				  "description": "Counter tracking the number of out of flow bytes received",
				  "defaultTags": [],
				  "properties": {},
				  "configuration": {}
				},
				{
				  "linkedItems": [],
				  "uid": "zwave:serial_zstick:94acf23fa3:serial_cse",
				  "id": "serial_cse",
				  "channelTypeUID": "zwave:serial_cse",
				  "itemType": "Number",
				  "kind": "STATE",
				  "label": "Received Checksum Errors",
				  "description": "Counter tracking the number of frames received with checksum errors",
				  "defaultTags": [],
				  "properties": {},
				  "configuration": {}
				}
			  ]
			}
		  ]`
	*/
}
