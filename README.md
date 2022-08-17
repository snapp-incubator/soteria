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
Soteria is a multivendor authenticator for EMQX. Follow instruction from [here](docs/vendor.md)

# Generate JWT Token
replace `driver` and `0` for issuer and id respectively.

```sh
curl -s -u 'admin:admin' -L https://doago-snapp-ode-020.apps.private.teh-1.snappcloud.io/api/snapp/driver/0  | jq '.Token' -r
```

# Folder Structure

- `.api`: API documentation like swagger files
- `.gitlab`: Gitlab CI templates
- `.okd`: OpenShift deployment configs (no longer is use. please use Helm charts)
- `deployments`: Helm Charts
- `internal`: Main application directory for codes
- `pkg`: Go packages that their logic is independent of this project and can become handy in other projects as well.
- `test`: test data like jwt keys
