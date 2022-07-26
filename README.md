# Soteria
Soteria is responsible for Authentication and Authorization of every request sent to EMQ and Herald.

# Deployment

## Staging

You can deploy `soteria` to the staging environments using
helm charts.

```
cd deployments
helm install soteria --generate-name
```

## Production

We deploy `soteria` on two different infrastructures.

- VM
- Cloud (okd)

### VM

For VM deployments following 2 steps are required.

- Preparing VMs which is done by `ansible playbooks`. You can find `herald`'s
  `ansible playbooks` in the following link
  [bravo/new-ansible-playbook](https://gitlab.snapp.ir/bravo/new-ansible-playbook)

- Deploying `soteria` with CI/CD pipelines.

### Cloud (okd)

[dispatching/ignite](https://gitlab.snapp.ir/dispatching/ignite) is responsible
for production deployments on Cloud (okd).

# Add Vendor
Soteria is a multivendor authenticator for EMQX. to add a vendor for authentication, go to chart directory and in the`values.yaml`, add the following template named after the desired vendor:
```yaml
snapp:
  company: "snapp"
  driver_salt: "secret"
  passenger_salt: "secret"
  driver_hash_length: 15
  passenger_hash_length: 15
  allowed_access_types: [ "pub", "sub" ]
  topics:
    - type: cab_event
      template: ^{{.audience}}-event-{{.hashId}}$
      hash_type: 1
      accesses:
        driver: '1'
        passenger: '1'
    - ...
  driver_key: |-
    ...
  passenger_key: |-
    ...
```

# Folder Structure

- `.api`: API documentation like swagger files
- `.gitlab`: Gitlab CI templates
- `.okd`: OpenShift deployment configs (no longer is use. please use Helm charts)
- `deployments`: Helm Charts
- `internal`: Main application directory for codes
- `pkg`: Go packages that their logic is independent of this project and can become handy in other projects as well.
- `test`: test data like jwt keys
