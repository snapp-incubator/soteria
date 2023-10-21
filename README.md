<h1 align="center"> Soteria </h1>

<p align="center">
    <img src="./.github/assets/logo.jpg" height="250px">
</p>

<p align="center">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/snapp-incubator/soteria/test.yaml?logo=github&style=for-the-badge">
    <img alt="Codecov" src="https://img.shields.io/codecov/c/github/snapp-incubator/soteria?logo=codecov&style=for-the-badge">
    <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/snapp-incubator/soteria?logo=github&style=for-the-badge">
</p>

## Introduction

Soteria is responsible for Authentication and Authorization of every request sent to [EMQ](https://github.com/emqx/emqx/).
The following configuration in [HOCON](https://github.com/lightbend/config/blob/master/HOCON.md)
format, configure EMQ to use HTTP Service for Authentication and Authorization.

```hocon
{
    mechanism = password_based
    backend = http
    enable = true

    method = post
    url = "http://127.0.0.1:8080/v2/auth"
    body {
        username = "${username}"
        password = "${password}"
        token = "${username}"
        clientid = "${clientid}"
    }
    headers {
        "Content-Type" = "application/json"
        "X-Request-Source" = "EMQX"
    }
}
```

```hocon
{
    type = http
    enable = true

    method = post
    url = "http://127.0.0.1:32333/v2/acl"
    body {
        username = "${username}"
        topic = "${topic}"
        action = "${action}"
    }
    headers {
        "Content-Type" = "application/json"
        "X-Request-Source" = "EMQX"
    }
}

```

We are using the [Authentication HTTP Service](https://www.emqx.io/docs/en/v5.2/access-control/authn/http.html)
and [Authorization HTTP Service](https://www.emqx.io/docs/en/v5.2/access-control/authn/http.html)
plugins of EMQ for forwarding these requests to Soteria and doing Authentication and Authorization.
EMQ has caching mechanism, but it sends requests almost for each Publish message to Soteria.
PS: On Subscribe we have only one message from client that need authorization and other messages are coming from server.

## Architecture

![arch](docs/arch.png)

## Support Vendors

Soteria supports having multiple vendors at the same time.
Means you can use single cluster for multiple companies at the same time and validate their tokens
and control accesses.

### Vendor Configuration

```yaml
company: "<<company_name>>"
driver_salt: ""
passenger_salt: ""
passenger_hash_length: 15
driver_hash_length: 15
allowed_access_types: ["pub", "sub"]
keys:
  iss-0: "key-value"
  iss-1: "key-value"
iss_entity_map:
  0: "entity-0"
  1: "entity-1"
  default: "default-entity"
iss_peer_map:
  0: "peer-0"
  1: "peer-1"
  default: "default-peer"
jwt:
  iss_name: "iss"
  sub_name: "sub"
  signing_method: "RS512"
topics:
  - topic1
  - topic2
  - ...
```

### HashID Manager

`driver_salt`,`passenger_salt`, `passenger_hash_length`, `driver_hash_length` are used for HashIDManager.
This component only works for passenger and driver issuers.

### Keys

This is a map of issuer to key for opening JWT token.

#### IssEntityMap & IssPeerMap

These two configuration map iss to entity and peer respectively.

**Note**: default case is `required`

```yaml
iss_entity_map:
  0: "driver"
  1: "passenger"
  default: "none"
iss_peer_map:
  0: "passenger"
  1: "driver"
  default: "none"
```

In the example above, we have two maps for entity & peer maps. As it's clear for **entity** structure **0** and **1** is mapped to **driver** and **passenger**, respectively. Vice Versa, for peer structure it can be seen that **1** and **0** is mapped to **driver** and **passenger**. We have also the **default** key for both two cases.

In the topic example, we have an accesses section in which **0** is mapped to **2** and **1** is mapped to **-1** which can be interpreted as a map from **IssEntity's Keys** to **Access Types**. In the other words this structure means:

- **Driver** has a **Pub** access on topic
- **Passenger** has a **None** access on topic (No Access)

### JWT

This is the JWT configuration. `iss_name` and `sub_name` are the name of issuer
and subject in the JWT token's payload respectively.

`signing_method` is the method that is used to sign the JWT token.
Here are list of different signing methods

- ES384
- RS512 \*
- PS512
- RS384 \*
- HS256 \*
- HS384 \*
- RS256 \*
- PS384
- ES256
- ES512
- EdDSA
- HS512 \*
- PS256
  **Note**: only the methods with `*` are supported for now.

### Topic Configuration

```yaml
type: "<<Name>>"
template: "<<regex template>>"
hash_type: 0|1
accesses:
  iss-0: "<<access>>"
  iss-1: "<<access>>"
```

### Template

Topic template is a string consist of [Variables](##Available_Variables) and [Functions](##Available_Functions)
and regular expressions.

Variables and Function are replaced first, and then the whole template will compile as a regular expression.
The end result will be compared against the requested topic.

#### Example

This is template topic given in `vendor:topics[#]:template`.

```yaml
- type: driver_location
  template: ^{{.company}}/driver/{{.sub}}/location$
  accesses:
    0: "2"
    1: "-1"
```

```regex
^{{.company}}/driver/{{HashID .hashType .sub (IssToSnappID .iss)}}/location/[a-zA-Z0-9-_]+$
```

After parsing the template we get something like this

```regex
// company=snapp
// hashType=0
// sub=D96ZbvJakLp4PYd
// iss=0
^snapp/driver/D96ZbvJakLp4PYd/location/[a-zA-Z0-9-_]+$
```

Now if the requested topic match the created topic it is considered as a valid topic for that particular user.

```text
requested_topic: snapp/driver/D96ZbvJakLp4PYd/location/23fw49vxd
created_topic_regex: ^snapp/driver/D96ZbvJakLp4PYd/location/[a-zA-Z0-9-_]+$
```

#### Available Variables

These are the variables available to use in the topic templates.

- `iss`
  issuer obtained from JWT token
- `sub`
  subject obtained from JWT token
- `hashType`
  Hash type field defined in topic template configuration

  | HashType | Value |
  | -------- | ----- |
  | HashID   | 0     |
  | MD5      | 1     |

- `company`
  company field defined in vendor configuration

#### Available Functions

These are the function available to use in the topic templates.

- `IssToEntity(iss string) string`
  convert `iss` obtained from JWT token to defined entity in `issEntityMap`
- `IssToPeer(iss string) string`
  convert `iss` obtained from JWT token to define peer in `issPeerMap`
- `IssToSnappID(iss string) string`
  convert `iss` obtained from JWT token to `snappid.audience`
- `HashID(hashType int, sub string, snappID snappid.audience)`
  generated `hashID` for the given `subject` base on the `hashType` and `snappid.audience`

**Note**: `snappid.audience` only is available for issuer 0 and 1 which are for driver and passenger respectively.

#### Accesses

List of all types of access on a topic.

| Access              | Value |
| ------------------- | ----- |
| Subscribe           | 1     |
| Publish             | 2     |
| Subscribe & Publish | 3     |
| None                | -1    |

#### Suggested Issuers

Use any value for issuer but if you have an entity called `Driver` or `Passenger`,
we recommend use the following issuers for them.

| Issuer    | Value |
| --------- | ----- |
| Driver    | 0     |
| Passenger | 1     |
