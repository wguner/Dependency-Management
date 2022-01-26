import click
import sys
import os
import tarfile
import grpc

from src.filesystem_interface.filesystem_interface import FilesystemInterface

import src.network_interface.PackageOperations.PackageOperations_pb2
import src.network_interface.PackageOperations.PackageOperations_pb2_grpc
from src.network_interface.FileTransfer.fileserver_client import FileTransfer

import src.network_interface.ProjectOperations.ProjectOperations_pb2 as ProjectOperations_pb2
import src.network_interface.ProjectOperations.ProjectOperations_pb2_grpc as ProjectOperations_pb2_grpc

import src.filesystem_interface.filesystem_interface
import src.network_interface.ServerUtils.ServerUtils as serverUtils

# Entry-point for the command line interface. Appears as 'packagebird'.
@click.group()
@click.pass_context
def cli(ctx):
    """Client utility for managing packages, is a group command for the subcommands below. Call --help on a subcommand to get further details."""
    # The entry point for the command line interface
    # Should ping server from this point to check that there is a server to connect with at the specified point
    if not serverUtils.ServerUtils.ping('127.0.0.1', '50051'):
        exit(0)


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
    
    # Close if called on directory not formatted as project
    if not FilesystemInterface.check_if_project_dir():
        click.echo('Directory is not configured as project. Please either configure as project or navigate to configured directory.')
        return
    
    project_configuration = FilesystemInterface.get_project_config()

    # Get name and version, 
    project_name = project_configuration['name']
    project_version = project_configuration['version']
    message = f'Creating package from project: {project_name}, version: {project_version}'
    print(message)

    # Package name
    package_name = f'{project_name}-v{project_version}.tar.gz'

    # Create compressed archive of file contents
    FilesystemInterface.make_archive(project_name, project_version)

    # Upload to server
    if (not debug):
        fileservice = FileTransfer()
        fileservice.upload('127.0.0.1', '50051', package_name)
    os.remove(package_name)

# Create project both locally and on the remote server
@cli.command('createproject', short_help='Creates a project in the local development environment and remotely on the server')
@click.option('-n', '--name', help='Name of project being created')
@click.option('-d', '--description', help='Description of new project being created')
@click.pass_context
def createproject(ctx, name, description):

    # Check if project already present in server
    request = ProjectOperations_pb2.ProjectRequest(name=name, description=description)
    
    # Needs moved to discrete logical section
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = ProjectOperations_pb2_grpc.ProjectOperationServicesStub(channel)
        response = stub.GetProject(request)
        if response.exist:
            click.echo(f"Project {name} already present on server. Please pick a different name.")
            return

    # Create the directory or notify user directory already present
    if FilesystemInterface.create_dir(name):
        click.echo(f"Creating {name} project directory...")
    else:
        click.echo(f"Directory {name} already present at location!")
        return
    
    # Formatting newly created project directory
    FilesystemInterface.make_project_dir(name, 0)

    pass