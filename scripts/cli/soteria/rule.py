import dataclasses


@dataclasses.dataclass
class Rule:
    access_type: str
    endpoint: str
    topic: str
    uuid: str
