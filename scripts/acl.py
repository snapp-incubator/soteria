import requests 

from vars import *

base_url = 'https://soteria-snapp-ode-012.apps.private.teh-1.snappcloud.io/'
acl_url = base_url + 'acl'


res = requests.post(acl_url, data={
    'access': PUBLISH,
    'token': DRIVER_TOKEN,
    'topic': 'snapp/driver/1234/location'
    })

print(res.content)