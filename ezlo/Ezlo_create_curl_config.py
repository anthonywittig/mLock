#!/usr/bin/python
# Script for local logon to Ezlo hub to get logon details for curl HTTPS requests.
# Normal output can be used as curl config file.
# Written by Rene Boer , 31 July 2020
# Version 1.0
#
# Use option --help for full usage details.

#
#
# Taken from https://community.getvera.com/t/python-script-for-ezlo-fw-http-api-curl-commands/214852
#
#

debug = False
certverify = False

import requests
import json
import sys
import hashlib
import urllib3
import uuid
import http.client
import logging
import argparse
import base64

# WS/HTTPS port to use for local access
ezlo_port = '17000'
wss_user = ''
wss_token = ''
wss_authorization = ''

def EzloLoginLocal(user_id, password, controller_ip, serial):
    # If we have the logon data for local, no need to get it again.
    global wss_user
    global wss_token
    global wss_authorization
    wss_user = ''
    wss_token = ''
    wss_authorization = ''

    # Not stored, logon to portal to get them.
    if wss_user == '' or wss_token == '' or wss_authorization == '':
        Ezlo_MMS_salt = "oZ7QE6LcLJp6fiWzdqZc"
        authentication_url = 'https://vera-us-oem-autha11.mios.com/autha/auth/username/{user}?SHA1Password={pwd}&PK_Oem=1&TokenVersion=2'
        get_token_url = 'https://cloud.ezlo.com/mca-router/token/exchange/legacy-to-cloud/'
        sync_token_url = 'https://api-cloud.ezlo.com/v1/request'

        # Get Tokens
        request_headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/json;charset=UTF-8',
            'Access-Control-Allow-Origin': '*',
            'User-Agent': 'RB HTTP 1.0'
        }
        hash = user_id.lower()+password+Ezlo_MMS_salt
        SHA1pwd = hashlib.sha1(hash.encode())
        if debug:
            print("====================================================================")
            print("Logon to portal for MMS Keys")
        #print(authentication_url.format(user=user_id,pwd=SHA1pwd.hexdigest()))
        response = requests.get(authentication_url.format(user=user_id,pwd=SHA1pwd.hexdigest()), headers=request_headers, verify=certverify)
        if response.status_code != 200:
            return 'cannot logon to portal',0
        if debug:
            print("\n")
            print(response.text)
            print("====================================================================\n")
        js_resp = response.json()

        if debug:
            print("====================================================================")
            print("Get Token")
        MMSAuth = js_resp.get('Identity')
        MMSAuthSig = js_resp.get('IdentitySignature')
        request_headers["MMSAuth"] = MMSAuth
        request_headers["MMSAuthSig"] = MMSAuthSig
        response = requests.get(get_token_url, headers=request_headers, verify=certverify)
        if response.status_code != 200:
            return 'cannot get token',0
        if debug:
            print("\n")
            print(response.text)
            print("====================================================================\n")
        js_resp = response.json()
        token=js_resp.get('token')

        # Get controller keys (user & token)
        if debug:
            print("====================================================================")
            print("Get controller keys")
        new_uuid = uuid.uuid4()
        post_headers = {
            'Authorization': 'Bearer '+token,
            'Accept': 'application/json',
            'Content-Type': 'application/json;charset=UTF-8',
            'Access-Control-Allow-Origin': '*',
            'User-Agent': 'RB Vera Bridge 1.0'
        }
        post_data = {
            'call': 'access_keys_sync',
            'version': '1',
            'params': {
                'version': 53,
                'entity': 'controller',
                'uuid': str(new_uuid)
            }
        }
        response = requests.post(sync_token_url, json=post_data, headers=post_headers, verify=certverify)
        if response.status_code != 200:
            return '',0
        if debug:
            print("\n")
            print(response.text)
            print("====================================================================\n")
        js_resp = response.json()
        # Get user and token from response.
        data = js_resp.get('data')
		# first look up controller uuid
        contr_uuid = ''
        for key in data.get('keys'):
            key_data = data.get('keys').get(key)
            if key_data.get('meta'):
                if key_data.get('meta').get('entity'):
                    if key_data.get('meta').get('entity').get('id'):
                        if key_data.get('meta').get('entity').get('id') == serial:
                            contr_uuid = key_data.get('meta').get('entity').get('uuid')
                        else:
                            print("Non-matching serial found: " + key_data.get('meta').get('entity').get('id'))
        if contr_uuid == '':
            print("Controller serial not found\n")
            return '',0
        for key in data.get('keys'):
            key_data = data.get('keys').get(key)
            if key_data.get('data') and wss_user == '' and wss_token == '':
                if key_data.get('data').get('string'):
                    if key_data.get('meta').get('target').get('uuid') == contr_uuid:
                        wss_token = key_data.get('data').get('string')
                        wss_user = key_data.get('meta').get('entity').get('uuid')
                        wss_data = wss_user+":"+wss_token
                        print("--- wss_data: " + wss_data)
                        wss_authorization = str(base64.b64encode(wss_data.encode("utf-8")), "utf-8")
    else:
        if debug:
            print("Using credentials from file.")

    return 'OK'

if __name__ == '__main__':
    # parse arguments
    parser = argparse.ArgumentParser(description='Create a curl config file to send commands to Ezlo Hub using HTTPS')
    parser.add_argument('ip', help='Your Ezlo Hub IP Address.')
    parser.add_argument('serial', help='Your Ezlo Hub Serial.')
    parser.add_argument('user', help='Your Ezlo Hub user id.')
    parser.add_argument('password', help='Your Ezlo Hub password.')
    parser.add_argument('-d', '--debug', action="store_true", help='Show debug commands.')
    args = parser.parse_args()
    user_id = args.user
    pwd = args.password
    controller_ip = args.ip
    serial = args.serial
    if args.debug:
        debug = True

    # Enable debugging of http requests (gives more details on Python 2 than 3 it seems)
    if debug:
        http.client.HTTPConnection.debuglevel = 1
        logging.basicConfig()
        logging.getLogger().setLevel(logging.DEBUG)
        requests_log = logging.getLogger("urllib3")
        requests_log.setLevel(logging.DEBUG)
        requests_log.propagate = True
    else:
        urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

    msg = 'No Hub serial or IP address specified'
    if controller_ip:
        msg = EzloLoginLocal(user_id, pwd, controller_ip, serial)
    if msg != 'OK':
        print('Failed to login ', msg)
        sys.exit()
    # Echo curl config file
    print("-H \"Authorization: Basic {}\"".format(wss_authorization))
    #print("-H \"user: {}\"".format(wss_user))
    #print("-H \"token: {}\"".format(wss_token))
    print("--insecure")
    print("--http1.1")

