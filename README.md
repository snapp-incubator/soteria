# Soteria

# What is Soteria?

Soteria is responsible for Authentication and Authorization of every request sent to EMQ and Herald.

# How to compile?

- [Install Golang](https://golang.org/doc/install)

- Set `snapp goproxy`

`go env -w GOPROXY="https://repo.snapp.tech/repository/goproxy/"`

- Run the following command to compile the application

`make compile`

# How to run it locally?

By executing the following command Herald will be up with EMQX and RabbitMQ brokers.

`make up`

# How to test?

## Unit testing

`make test`

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

# Folder Structure

- `.api`: API documentation like swagger files
- `.gitlab`: Gitlab CI templates
- `.okd`: OpenShift deployment configs (no longer is use. please use Helm charts)
- `deployments`: Helm Charts
- `internal`: Main application directory for codes
- `pkg`: Go packages that their logic is independent of this project and can become handy in other projects as well.
- `test`: test data like jwt keys
