package shared

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	HABThing          HABThing   `json:"habThing"`
	ID                uuid.UUID  `json:"id"`
	LastRefreshedAt   time.Time  `json:"lastRefreshedAt"`
	LastWentOfflineAt *time.Time `json:"lastWentOfflineAt"`
	PropertyID        uuid.UUID  `json:"propertyId"`
	UnitID            *uuid.UUID `json:"unitId"`
}

type HABThing struct {
	Configuration struct {
		Config_10_1     int      `json:"config_10_1"`
		Config_11_1     int      `json:"config_11_1"`
		Config_12_4     int      `json:"config_12_4"`
		Config_13_4     int      `json:"config_13_4"`
		Config_14_4     int      `json:"config_14_4"`
		Config_15_1     int      `json:"config_15_1"`
		Config_16_1     int      `json:"config_16_1"`
		Config_17_4     int      `json:"config_17_4"`
		Config_18_1     int      `json:"config_18_1"`
		Config_3_1      int      `json:"config_3_1"`
		Config_4_1      int      `json:"config_4_1"`
		Config_5_1      int      `json:"config_5_1"`
		Config_6_4      int      `json:"config_6_4"`
		Config_7_1      int      `json:"config_7_1"`
		Config_8_1      int      `json:"config_8_1"`
		Config_9_1      int      `json:"config_9_1"`
		DoorlockTimeout int      `json:"doorlock_timeout"`
		Group1          []string `json:"group_1"`
		NodeID          int      `json:"node_id"`
		UsercodeCode1   string   `json:"usercode_code_1"`
		UsercodeCode2   string   `json:"usercode_code_2"`
		UsercodeCode3   string   `json:"usercode_code_3"`
		UsercodeCode4   string   `json:"usercode_code_4"`
		UsercodeCode5   string   `json:"usercode_code_5"`
		UsercodeCode6   string   `json:"usercode_code_6"`
		UsercodeCode7   string   `json:"usercode_code_7"`
		UsercodeCode8   string   `json:"usercode_code_8"`
		UsercodeCode9   string   `json:"usercode_code_9"`
		UsercodeCode10  string   `json:"usercode_code_10"`
		UsercodeCode11  string   `json:"usercode_code_11"`
		UsercodeCode12  string   `json:"usercode_code_12"`
		UsercodeCode13  string   `json:"usercode_code_13"`
		UsercodeCode14  string   `json:"usercode_code_14"`
		UsercodeCode15  string   `json:"usercode_code_15"`
		UsercodeCode16  string   `json:"usercode_code_16"`
		UsercodeCode17  string   `json:"usercode_code_17"`
		UsercodeCode18  string   `json:"usercode_code_18"`
		UsercodeCode19  string   `json:"usercode_code_19"`
		UsercodeCode20  string   `json:"usercode_code_20"`
		UsercodeCode21  string   `json:"usercode_code_21"`
		UsercodeCode22  string   `json:"usercode_code_22"`
		UsercodeCode23  string   `json:"usercode_code_23"`
		UsercodeCode24  string   `json:"usercode_code_24"`
		UsercodeCode25  string   `json:"usercode_code_25"`
		UsercodeCode26  string   `json:"usercode_code_26"`
		UsercodeCode27  string   `json:"usercode_code_27"`
		UsercodeCode28  string   `json:"usercode_code_28"`
		UsercodeCode29  string   `json:"usercode_code_29"`
		UsercodeCode30  string   `json:"usercode_code_30"`
	} `json:"configuration"`
	Label      string `json:"label"`
	StatusInfo struct {
		Status       string `json:"status"`
		StatusDetail string `json:"statusDetail"`
	} `json:"statusInfo"`
	ThingTypeUID string `json:"thingTypeUID"`
	UID          string `json:"UID"`
}

const (
	DeviceStatusOffline = "OFFLINE"
	DeviceStatusOnline  = "ONLINE"
	// I think "UNINITIALIZED" is another one, maybe.
)

/*
OpenHAB things:
{
    "bridgeUID": "zwave:serial_zstick:94acf23fa3",
    "channels": [
        {
            "channelTypeUID": "system:battery-level",
            "configuration": {},
            "defaultTags": [],
            "id": "battery-level",
            "itemType": "Number",
            "kind": "STATE",
            "label": "Battery Level",
            "linkedItems": [],
            "properties": {
                "binding:*:PercentType": "COMMAND_CLASS_BATTERY"
            },
            "uid": "zwave:device:94acf23fa3:node2:battery-level"
        },
        {
            "channelTypeUID": "zwave:alarm_access",
            "configuration": {},
            "defaultTags": [],
            "description": "Indicates if the access control alarm is triggered",
            "id": "alarm_access",
            "itemType": "Switch",
            "kind": "STATE",
            "label": "Alarm (access)",
            "linkedItems": [],
            "properties": {
                "binding:*:OnOffType": "COMMAND_CLASS_ALARM;type=ACCESS_CONTROL"
            },
            "uid": "zwave:device:94acf23fa3:node2:alarm_access"
        },
        {
            "channelTypeUID": "zwave:alarm_power",
            "configuration": {},
            "defaultTags": [],
            "description": "Indicates if a power alarm is triggered",
            "id": "alarm_power",
            "itemType": "Switch",
            "kind": "STATE",
            "label": "Alarm (power)",
            "linkedItems": [],
            "properties": {
                "binding:*:OnOffType": "COMMAND_CLASS_ALARM;type=POWER_MANAGEMENT"
            },
            "uid": "zwave:device:94acf23fa3:node2:alarm_power"
        },
        {
            "channelTypeUID": "zwave:alarm_raw",
            "configuration": {},
            "defaultTags": [],
            "description": "Provides alarm parameters as json string",
            "id": "alarm_raw",
            "itemType": "String",
            "kind": "STATE",
            "label": "Alarm (raw)",
            "linkedItems": [],
            "properties": {
                "binding:*:StringType": "COMMAND_CLASS_ALARM"
            },
            "uid": "zwave:device:94acf23fa3:node2:alarm_raw"
        },
        {
            "channelTypeUID": "zwave:alarm_system",
            "configuration": {},
            "defaultTags": [],
            "description": "Indicates if a system alarm is triggered",
            "id": "alarm_system",
            "itemType": "Switch",
            "kind": "STATE",
            "label": "Alarm (system)",
            "linkedItems": [],
            "properties": {
                "binding:*:OnOffType": "COMMAND_CLASS_ALARM;type=SYSTEM"
            },
            "uid": "zwave:device:94acf23fa3:node2:alarm_system"
        },
        {
            "channelTypeUID": "zwave:lock_door",
            "configuration": {},
            "defaultTags": [],
            "description": "Lock and unlock the door",
            "id": "lock_door",
            "itemType": "Switch",
            "kind": "STATE",
            "label": "Door Lock",
            "linkedItems": [],
            "properties": {
                "binding:*:OnOffType": "COMMAND_CLASS_DOOR_LOCK"
            },
            "uid": "zwave:device:94acf23fa3:node2:lock_door"
        }
    ],
    --"configuration": {
        "config_10_1": 3,
        "config_11_1": 0,
        "config_12_4": 4,
        "config_13_4": -1,
        "config_14_4": -1,
        "config_15_1": 0,
        "config_16_1": 4,
        "config_17_4": -1,
        "config_18_1": 6,
        "config_3_1": -1,
        "config_4_1": 0,
        "config_5_1": 0,
        "config_6_4": 50331648,
        "config_7_1": 0,
        "config_8_1": 3,
        "config_9_1": 3,
        "doorlock_timeout": 0,
        "group_1": [
            "controller"
        ],
        "node_id": 2,
        "usercode_code_1": "37 35 33 34 0A 0D",
        "usercode_code_10": "",
        "usercode_code_11": "",
        "usercode_code_12": "",
        "usercode_code_13": "",
        "usercode_code_14": "",
        "usercode_code_15": "",
        "usercode_code_16": "",
        "usercode_code_17": "",
        "usercode_code_18": "",
        "usercode_code_19": "",
        "usercode_code_2": "33 38 31 32 0A 0D",
        "usercode_code_20": "",
        "usercode_code_21": "",
        "usercode_code_22": "",
        "usercode_code_23": "",
        "usercode_code_24": "",
        "usercode_code_25": "",
        "usercode_code_26": "",
        "usercode_code_27": "",
        "usercode_code_28": "",
        "usercode_code_29": "",
        "usercode_code_3": "",
        "usercode_code_30": "",
        "usercode_code_4": "",
        "usercode_code_5": "",
        "usercode_code_6": "",
        "usercode_code_7": "",
        "usercode_code_8": "",
        "usercode_code_9": ""
    },
    "editable": true,
    --"label": "Murset's schlage lock",
    "properties": {
        "dbReference": "1223",
        "defaultAssociations": "1",
        "manufacturerId": "003B",
        "manufacturerRef": "0001:0468",
        "modelId": "BE468ZP",
        "vendor": "Allegion",
        "zwave_beaming": "true",
        "zwave_class_basic": "BASIC_TYPE_ROUTING_SLAVE",
        "zwave_class_generic": "GENERIC_TYPE_ENTRY_CONTROL",
        "zwave_class_specific": "SPECIFIC_TYPE_SECURE_KEYPAD_DOOR_LOCK",
        "zwave_deviceid": "1128",
        "zwave_devicetype": "1",
        "zwave_frequent": "true",
        "zwave_lastheal": "2021-05-31T08:35:32Z",
        "zwave_listening": "false",
        "zwave_manufacturer": "59",
        "zwave_neighbours": "1",
        "zwave_nodeid": "2",
        "zwave_plus_devicetype": "NODE_TYPE_ZWAVEPLUS_NODE",
        "zwave_plus_roletype": "ROLE_TYPE_SLAVE_SLEEPING_LISTENING",
        "zwave_routing": "true",
        "zwave_secure": "true",
        "zwave_version": "3.3"
    },
    --"statusInfo": {
        --"status": "ONLINE",
        --"statusDetail": "NONE"
    },
    --"thingTypeUID": "zwave:schlage_be468zp_00_000",
    --"UID": "zwave:device:94acf23fa3:node2"
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
*/
