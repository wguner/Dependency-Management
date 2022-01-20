import click
import sys

import src.network_interface.PackageOperations_pb2
import src.network_interface.PackageOperations_pb2_grpc

# Entry-point for the command line interface. Appears as 'packagebird'.
@click.group()
@click.pass_context
def cli(ctx):
    """Client utility for managing packages, is a group command for the subcommands below. Call --help on a subcommand to get further details."""
    # The entry point for the command line interface

# Add package to the development directory
@cli.command('addpackage', short_help='Adds a package to development directory packages subdirectory')
@click.option('-p', '--package', 'name')
@click.option('-v', '--version', 'version')
@click.pass_context
def addpackage(ctx, name, version):
    # Makes the packages directory
    if not os.path.isdir('packages'):
        click.echo("Creating packages directory...")
        os.mkdir('packages')
    os.chdir('packages')
    
    # Create Package Directory
    os.mkdir(f'{name}')
    os.chdir(f'{name}')

    # Download the archived package
    fileservice = FileTransfer()
    request_string = f'{name}-v{version}.tar.gz'
    fileservice.download('127.0.0.1', '50051', f'{request_string}')

    # Extract the contents
    with tarfile.open(request_string, 'r:gz') as archive:
        try:
            archive.extractall(path='.')
        except tarfile.ReadError:
            click.echo('Empty file found.')
    
    # Remove the archive file
    os.remove(request_string)

    # Move back to the development root directory
    os.chdir('..')
    os.chdir('..')