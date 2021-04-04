import requests

from vars import *

base_url = "https://soteria-snapp-ode-012.apps.private.teh-1.snappcloud.io/"
token_url = base_url + "token"


def get_token(grant_type, client_id, client_secret):
    res = requests.post(
        token_url,
        data={
            "grant_type": grant_type,
            "client_id": client_id,
            "client_secret": client_secret,
        },
    )
    return res.content.decode("utf-8")


print(get_token(SUBSCRIE, "box", "secret"))
print(get_token(SUBSCRIE, "daghigh", "secret"))
