import requests
import base64
import pprint


base_url = "https://soteria-snapp-ode-012.apps.private.teh-1.snappcloud.io/"
pp = pprint.PrettyPrinter(indent=4)


def get_account(username, password):

    url = base_url + "accounts/" + username

    res = requests.get(url, auth=(username, password))
    return pp.pprint(res.json())


print(get_account("driver", "password"))
print(get_account("daghigh", "secret"))


def update_account(username, password, secret):
    url = base_url + "accounts/" + username

    body = {
        secret: "secret",
    }

    res = requests.put(url, auth=(username, password), json=body)
    return pp.pprint(res.json())


print(update_account("daghigh", "secret", "secret"))
