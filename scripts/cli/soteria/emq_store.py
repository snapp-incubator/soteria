from requests_toolbelt import sessions


class EMQStore:
    def __init__(self, base_url: str):
        self.session = sessions.BaseUrlSession(base_url=base_url)

    def new(self, username: str, password: str, duration: int):
        res = self.session.post(
            "emq/",
            json={
                "username": username,
                "password": password,
                "duration": duration,
            },
        )
        return res.json()
