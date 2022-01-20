import click
import sys
import os
import tarfile

import src.network_interface.PackageOperations.PackageOperations_pb2
import src.network_interface.PackageOperations.PackageOperations_pb2_grpc
from src.network_interface.FileTransfer.fileserver_client import FileTransfer

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

# Create package from a development directory
@cli.command('createpackage', short_help='Creates a package from the current directory')
@click.option('--debug', is_flag=True, help='Debug option for')
@click.pass_context
def createpackage(ctx, debug):
    # Get packages and location
    project_name = os.path.basename(os.getcwd())
    version = 1
    message = f'Project Name: {project_name}, version: {version}'
    print(message)

    # Package name
    package_name = f'{project_name}-v{version}.tar.gz'

    # Add the contents
    with tarfile.open(package_name, 'w:gz', format=tarfile.GNU_FORMAT) as tar:
        for directory, directorynames, filenames in os.walk("."):
            for file in filenames:
                if (debug):
                    click.echo(f'Directory Visited: {directory}')

                if "/packages/" not in directory and "\\packages\\" not in directory and "./packages/" not in directory:
                    if (debug):
                        click.echo(f'Visited directory and file not in packages. File reached is {file}.')
                    
                    if file != package_name:
                        tar.add(os.path.join(directory, file))
                    else:
                        if (debug):
                            click.echo(f'File {file} is the temp client tar file, should not be bundled.')
                elif debug:
                    click.echo(f'File visited in packages directory! Directory is {directory}, file is {file}')
                    


    # Upload to server
    if (not debug):
        fileservice = FileTransfer()
        fileservice.upload('127.0.0.1', '50051', package_name)
    os.remove(package_name)