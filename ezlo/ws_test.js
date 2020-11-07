const WebSocket = require('ws');
const url = "wss://nma-server6-ui-cloud.ezlo.com/"
const connection = new WebSocket(url)

connection.onopen = () => {
	console.log("onopen")
	sendMessage();
}

connection.onerror = error => {
  console.log(`WebSocket error: ${error}`)
}

connection.onmessage = e => {
	console.dir(JSON.parse(e.data), {depth: null, colors: true})
	//console.log(JSON.stringify(JSON.parse(e.data), null, "  "))
	console.log()
	sendMessage()
}


let messages = [
	{
		"method":"loginUserMios",
		"id":"loginUser",
		"params":{
			// Get from Auth (postman) call.
			"MMSAuth": "eyJFeHBpcmVzIjoxNjA0ODA2NDM5LCJHZW5lcmF0ZWQiOjE2MDQ3MjAwMzksIlBlcm1pc3Npb25zRW5hYmxlZCI6WzEsMiwzLDQsNiw5LDE0LDE1LDE3LDE4LDE5LDIwLDIxLDIyLDIzLDI1LDI2LDI3LDI4LDI5LDMwLDQxLDQyLDQzLDUyLDYyLDcyLDgyLDkyLDEwMiwxMTIsMTIyLDEzMiwxNDIsMTUyLDE2MiwxNzIsMTgyLDE5MiwyMDIsMjEyLDIyMiwyMzIsMjQyLDI1MiwyNjIsMjcyLDI4MiwyOTIsMzAyLDMxMiwzMjIsMzMyLDM0Miw0ODIsNDkyLDU2Miw1ODIsNjgyLDY4NSw3MDUsNzMxLDE1OTYsMTYwNiwxNjE2LDE2MjYsMTYzNSwxNjM2LDE2NDgseyJQSyI6MTY1MywiQXJndW1lbnRzIjoiWzksMTAsMTEsMTIsMTNdIn1dLCJQZXJtaXNzaW9uc0Rpc2FibGVkIjpbXSwiVmVyc2lvbiI6MiwiUEtfQWNjb3VudCI6MTE1NTUwMiwiUEtfQWNjb3VudFR5cGUiOjUsIlBLX0FjY291bnRDaGlsZCI6MCwiUEtfQWNjb3VudF9QYXJlbnQiOjIzNzUsIlBLX0FjY291bnRfUGFyZW50MiI6MSwiUEtfQWNjb3VudF9QYXJlbnQzIjowLCJQS19BY2NvdW50X1BhcmVudDQiOjAsIlBLX0FwcCI6MCwiUEtfT2VtIjoxLCJQS19PZW1fVXNlciI6IiIsIlBLX1Blcm1pc3Npb25Sb2xlIjoxMCwiUEtfVXNlciI6Mjc3MzI4MiwiUEtfU2VydmVyX0F1dGgiOjEsIlBLX1NlcnZlcl9BY2NvdW50Ijo1LCJQS19TZXJ2ZXJfRXZlbnQiOjEzLCJTZXJ2ZXJfQXV0aCI6InZlcmEtdXMtb2VtLWF1dGhhMTEubWlvcy5jb20iLCJTZXEiOjU1MTAzNCwiVXNlcm5hbWUiOiJ6aW9uc2NhbXBhbmRjb3R0YWdlcyJ9",
			"MMSAuthSig": "FtASRZZZg4w5Hn5HfPenCvnS7m5ibxgvRulrLzp57o2YO8XSLSOi0WIPCFMNK9nCErlEM8mWZuXOUURvz7LXqGZIyvvr+5a6oqP9J7AX85/dNgNUYi5YAPL489ijkAy97+0zFKTGOcvNHS5itxTyWxtVRsqb8wtF0+Uc5Ydn1qt3ELZq+HK8X2mCRYcTDI3s7DDgQzExCoaffZmrTfdWJ6suSIscEhySLncaaFe2DEHUYD40r5ogeMUrLW3KorxQArJSj/hFx/9CojQhOqVTdEphZbg3mMPW7CCTLqgQajRjgwsUUC+J5Sz1S9iwdoDFuKNrErfC10Yyl3iC7TuNdg==",
		}
	}, {
		"method":"register",
		"id":"register",
		"jsonrpc":"2.0",
		"params":{
			"serial":"45052564",
		},
	}, {
		"method":"hub.devices.list",
		"id":"1604717670276",
		"params":{},
	}, {
		"method":"hub.items.list",
		"id":"1604717670277",
		"params":{},
	},
	/*{
		"method":"hub.item.value.set",
		"id":"1604717670278",
		"params":{
			"_id": "5fa6292700000030860bd211",
			"valueType": "dictionary",
			"elementType": "dictionary.userCode",
			"value": {
				"4": {
					"code": "16345",
					"name": "my code",
				},
			},
		},
	},*/
	/*{
		"method":"hub.item.value.set",
		"id":"1604717670279",
		"params":{
			"_id": "5fa6292700000030860bd211",
			"value": {
				"1": {
					"code": "16345",
					"name": "my code",
				},
			},
		},
	},*/
	/*{
		// works - locks the door
		"method":"hub.item.value.set",
		"id":"1604717670278",
		"params":{
			"_id": "5fa6292700000030860bd20e",
			"value": "secured",
		},
	},*/
]

function sendMessage() {
	if (messages.length == 0) {
		console.log("no more messages")
		return
	}
	let next = JSON.stringify(messages.shift())
	if (next == "null") {
		// silently skip
		return
	}
	console.log(`\nsending message:\n${next}\n`)
	connection.send(next)
}

console.log("running?")

/*
{
  id: 'ui_broadcast',
  msg_subclass: 'hub.item.dictionary.updated',
  result: {
    _id: '5fa6292700000030860bd211',
    deviceCategory: 'door_lock',
    deviceId: '5fa6292600000030860bd206',
    deviceName: 'SmartCode 888',
    deviceSubcategory: '',
    element: { '1': { value: { code: '1634', name: 'Code 1' } } },
    elementType: 'userCode',
    elementsMaxNumber: 30,
    name: 'user_codes',
    notifications: null,
    operation: 'removed',
    roomName: '',
    serviceNotification: false,
    syncNotification: false,
    userCodeRestriction: '^\\d{4,10}$',
    userNotification: false,
    valueFormatted: '',
    valueType: 'dictionary'
  }
}
*/
