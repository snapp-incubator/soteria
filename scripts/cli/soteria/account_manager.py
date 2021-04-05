from requests_toolbelt import sessions
import dacite

from .rule import Rule
from .account import Account


class AccountManager:
    user_types = {
        "herald": "HeraldUser",
        "emq": "EMQUser",
        "staff": "Staff",
    }

    grant_types = {
        "sub": "1",
        "pub": "2",
        "pubsub": "3",
    }

    def __init__(self, base_url: str):
        self.session = sessions.BaseUrlSession(base_url=base_url)

    def show(self, username: str, password: str) -> Account:
        res = self.session.get(
            f"accounts/{username}", auth=(username, password)
        ).json()

        if res["data"] is not None:
            return dacite.from_dict(data_class=Account, data=res["data"])
        raise Exception(res["message"])

    def new(self, username: str, password: str, user_type: str):
        res = self.session.post(
            "accounts/",
            json={
                "username": username,
                "password": password,
                "user_type": self.user_types[user_type],
            },
        )
        return res.json()

    def add_rule(
        self, username: str, password: str, topic: str, access_type: str
    ) -> Rule:
        res = self.session.post(
            f"accounts/{username}/rules",
            json={
                "topic": topic,
                "access_type": access_type,
            },
            auth=(username, password),
        ).json()

        if res["data"] is not None:
            return dacite.from_dict(data_class=Rule, data=res["data"])
        raise Exception(res["message"])

    def set_secret(self, username: str, password: str, secret: str):
        res = self.session.put(
            f"accounts/{username}",
            json={
                "secret": secret,
            },
            auth=(username, password),
        )
        return res.json()

    def set_expiration(self, username: str, password: str, expiration: int):
        res = self.session.put(
            f"accounts/{username}",
            json={
                "token_expiration": expiration,
            },
            auth=(username, password),
        )
        return res.json()

    def token(self, username: str, secret: str, grant_type: str) -> bytes:
        """
        Generates a token from soteria.
        """
        res = self.session.post(
            "token",
            json={
                "client_id": username,
                "client_secret": secret,
                "grant_type": self.grant_types[grant_type],
            },
        )
        return res.content
