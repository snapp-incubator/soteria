import click
import soteria
import pprint


pp = pprint.PrettyPrinter(indent=4)


@click.group()
@click.option("--base", "-b", required=True, type=str)
@click.pass_context
def cli(ctx, base):
    ctx.ensure_object(dict)

    ctx.obj["BASE"] = base


@cli.command()
@click.option(
    "--username",
    "-u",
    required=True,
    type=str,
    help="account username e.g. driver",
)
@click.option(
    "--password",
    "-p",
    required=True,
    type=str,
    help="account password e.g. password",
)
@click.pass_context
def show(ctx, username: str, password: str):
    account_manager = soteria.AccountManager(ctx.obj["BASE"])
    try:
        resp = account_manager.show(username, password)
        click.echo(pp.pformat(resp))
    except Exception as exep:
        raise click.ClickException(str(exep))


@cli.command()
@click.option(
    "--username",
    "-u",
    required=True,
    type=str,
    help="account username e.g. driver",
)
@click.option(
    "--password",
    "-p",
    required=True,
    type=str,
    help="account password e.g. password",
)
@click.option(
    "--user-type",
    "-t",
    required=True,
    type=click.Choice(["herald", "emq", "staff"]),
)
@click.pass_context
def new(ctx, username: str, password: str, user_type: str):
    account_manager = soteria.AccountManager(ctx.obj["BASE"])
    try:
        resp = account_manager.new(username, password, user_type)
        click.echo(pp.pformat(resp))
    except Exception as exep:
        raise click.ClickException(str(exep))


@cli.command()
@click.option(
    "--username",
    "-u",
    required=True,
    type=str,
    help="account username e.g. driver",
)
@click.option(
    "--password",
    "-p",
    required=True,
    type=str,
    help="account password e.g. password",
)
@click.option(
    "--access-type",
    "-a",
    required=True,
    type=click.Choice(["pub", "sub", "pubsub"]),
)
@click.option(
    "--topic",
    "-t",
    required=True,
    type=str,
    help="topic regular expressions are defined in soteria and has their name"
    "so please refer to soteria",
)
@click.pass_context
def rules_add(ctx, username: str, password: str, access_type: str, topic: str):
    account_manager = soteria.AccountManager(ctx.obj["BASE"])
    try:
        resp = account_manager.add_rule(username, password, topic, access_type)
        click.echo(pp.pformat(resp))
    except Exception as exep:
        raise click.ClickException(str(exep))


@cli.command()
@click.option(
    "--username",
    "-u",
    required=True,
    type=str,
    help="account username e.g. driver",
)
@click.option(
    "--password",
    "-p",
    required=True,
    type=str,
    help="account password e.g. password",
)
@click.option(
    "--expire",
    "-e",
    required=True,
    type=str,
    help="token expiration time, e.g. 1h",
)
@click.pass_context
def set_expire(
    ctx,
    username: str,
    password: str,
    expire: str,
):
    account_manager = soteria.AccountManager(ctx.obj["BASE"])
    try:
        resp = account_manager.set_expiration(username, password, expire)
        click.echo(pp.pformat(resp))
    except Exception as exep:
        raise click.ClickException(str(exep))


@cli.command()
@click.option(
    "--username",
    "-u",
    required=True,
    type=str,
    help="account username e.g. driver",
)
@click.option(
    "--password",
    "-p",
    required=True,
    type=str,
    help="account password e.g. password",
)
@click.option(
    "--secret",
    "-s",
    required=True,
    type=str,
    help="account secret which is different from password and"
    "is used to generate token",
)
@click.pass_context
def set_secret(
    ctx,
    username: str,
    password: str,
    secret: str,
):
    account_manager = soteria.AccountManager(ctx.obj["BASE"])
    try:
        resp = account_manager.set_secret(username, password, secret)
        click.echo(pp.pformat(resp))
    except Exception as exep:
        raise click.ClickException(str(exep))


@cli.command()
@click.option(
    "--username",
    "-u",
    required=True,
    type=str,
    help="account username e.g. driver",
)
@click.option(
    "--secret",
    "-s",
    required=True,
    type=str,
    help="account secret which is different from password and"
    "is used to generate token",
)
@click.option(
    "--grant-type",
    "-t",
    required=True,
    type=click.Choice(["pub", "sub", "pubsub"]),
)
@click.pass_context
def token(
    ctx,
    username: str,
    secret: str,
    grant_type: str,
):
    account_manager = soteria.AccountManager(ctx.obj["BASE"])
    try:
        token = account_manager.token(username, secret, grant_type)
        click.echo(token)
    except Exception as exep:
        raise click.ClickException(str(exep))


if __name__ == "__main__":
    cli(obj={})
