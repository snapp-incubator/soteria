import requests 

from vars import *

base_url = 'https://soteria-snapp-ode-012.apps.private.teh-1.snappcloud.io/'
auth_url = base_url + 'acl'


res = requests.post(auth_url, data={
    'access': PUBLISH,
    'token': DRIVER_TOKEN,
    'topic': 'snapp/driver/1234/location'
    })

print(res.content)