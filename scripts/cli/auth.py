import requests

from vars import DRIVER_TOKEN

base_url = "https://soteria-snapp-ode-012.apps.private.teh-1.snappcloud.io"
auth_url = base_url + "/auth"


res = requests.post(
    auth_url,
    data={
        "token": DRIVER_TOKEN,
    },
)

print(res.content)
