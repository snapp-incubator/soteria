# Soteria

Soteria is responsible for Authentication and Authorization of every request sent to EMQ.

# Deployment

## Staging

You can deploy `soteria` to the staging environments using
helm charts.

```bash
cd deployments/soteria
helm install soteria .
```

## Production

We deploy `soteria` on `Cloud (okd)` infrastructures.

### Cloud (okd)

[dispatching/ignite](https://gitlab.snapp.ir/dispatching/ignite) is responsible
for production deployments on Cloud (okd).

# Add Vendor
Soteria is a multivendor authenticator for EMQX. Follow instruction from [here](docs/vendor.md)

# Generate JWT Token
Replace `driver` and `0` for issuer and ID respectively.

```bash
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
