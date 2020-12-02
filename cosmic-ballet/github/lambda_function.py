import json
import uuid
import os
import boto3
from http import client
from urllib.parse import urlencode
from urllib.request import Request, urlopen
from base64 import b64decode, b64encode

github_root = 'https://github.com/login/oauth'
allow_signup = 'false'

app_domain = os.environ.get('SIMPSONS_CONFIG_APP_DOMAIN')
redirect_url = os.environ.get('SIMPSONS_CONFIG_REDIRECT_URL')
scope = os.environ.get('SIMPSONS_CONFIG_SCOPES')
client_id = os.environ.get('GITHUB_CLIENT_ID')
kms_key_alias = os.environ.get('SIMPSON_CONFIG_KMS_ALIAS')


def lambda_handler(event, context):

    if ("callback" in event['path']):
        return handle_callback(event)

    return handle_start(event)


def handle_start(event):
    state = uuid.uuid4()
    print(f'Generate State {state} for OAuth flow')
    return {
        'statusCode': 302,
        'headers': {
            'location': f'{github_root}/authorize?client_id={client_id}&redirect_url={redirect_url}&allow_signup={allow_signup}&scope={scope}&state={state}',
            'set-cookie': f'session_state={state}; Domain={app_domain}; Path=/; Max-Age=20; Secure; HttpOnly'
        }
    }


def handle_callback(event):
    session = boto3.session.Session()
    client = session.client('kms')

    encrypted_secret = b64decode(os.environ.get('GITHUB_CLIENT_SECRET'))

    client_secret = client.decrypt(CiphertextBlob=bytes(
        encrypted_secret), EncryptionContext={'LambdaFunctionName': os.environ['AWS_LAMBDA_FUNCTION_NAME']})['Plaintext'].decode('utf-8')

    if 'Cookie' not in event['headers']:
        print('Cookie not in header, kicking off OAuth flow')
        return handle_start(event)

    cookie_state = event['headers']['Cookie'].split("=")[1]
    param_state = event['queryStringParameters']['state']
    param_code = event['queryStringParameters']['code']

    if cookie_state != param_state:
        print('Cookie state did not match query param state, returning to main app')
        print(f'Expected {cookie_state} but got {param_state}')
        return {
            'statusCode': 302,
            'headers': {
                'location': f'https://{app_domain}',
            }
        }

    print(
        f'Creating Auth Param Call with ClientID {client_id} and code {param_code}')
    params = {
        'client_id': client_id,
        'client_secret': client_secret,
        'code': param_code
    }

    headers = {'Accept': 'application/json'}

    encodedParams = urlencode(params).encode("utf-8")

    req = Request(url=f'{github_root}/access_token',
                  data=encodedParams, method="POST", headers=headers)

    with urlopen(req) as res:
        token = json.load(res)
        session_token = b64encode(client.encrypt(
            KeyId=kms_key_alias, Plaintext=bytes(token['access_token'].encode("utf-8")))["CiphertextBlob"])

        return {
            'statusCode': 302,
            'headers': {
                'location': f'https://{app_domain}',
                'set-cookie': f'session_access={str(session_token, "utf-8")}; Domain={app_domain}; Path=/; Secure; HttpOnly'
            }
        }
