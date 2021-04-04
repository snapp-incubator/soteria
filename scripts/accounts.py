import requests
import json
import base64


base_url = 'https://soteria-snapp-ode-012.apps.private.teh-1.snappcloud.io/'


def get_account(username, password):
    url = base_url + 'accounts/' + username
    auth_bytes = '{}:{}'.format(username, password).encode('ascii')
    auth = base64.b64encode(auth_bytes).decode('ascii')
    
    res = requests.get(url, headers={
        "Authorization": "Basic " + auth
    })
    return json.dumps(res.json(), indent=4, sort_keys=True)



print(get_account('driver', 'password'))