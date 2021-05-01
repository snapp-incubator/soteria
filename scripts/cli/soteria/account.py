import dataclasses
import typing
from .rule import Rule


@dataclasses.dataclass
class Account:
    password: str
    rules: typing.Optional[typing.List[Rule]]
    secret: str
    type: str
    username: str
    token_expiration_duration: int
