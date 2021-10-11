import requests
import click

@click.command()
@click.option('--address', help='address of server')
@click.option('--name', help='name of client')

def request(address, name):
    endpoint = f"http://{address}/?name={name}/"
    click.echo(f"Endpoint:\t{endpoint}")
    request = requests.get(endpoint)
    click.echo(request.text)

if __name__ == 'main':
    print("Running Client.py")
    request()

request()