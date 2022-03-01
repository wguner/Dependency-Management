from ast import Str
from distutils.version import Version
from inspect import getmembers
from posixpath import split
from re import T
from urllib import request
import click
import sys
import os
import tarfile
import grpc

from src.filesystem_interface.filesystem_interface import FilesystemInterface

import src.network_interface.PackageOperations.PackageOperations_pb2 as packageOperations
import src.network_interface.PackageOperations.PackageOperations_pb2_grpc as packageOperationsGRPC
from src.network_interface.FileTransfer.fileserver_client import FileTransfer

import src.network_interface.ProjectOperations.ProjectOperations_pb2 as ProjectOperations_pb2
import src.network_interface.ProjectOperations.ProjectOperations_pb2_grpc as ProjectOperations_pb2_grpc

import src.filesystem_interface.filesystem_interface
import src.network_interface.ServerUtils.ServerUtils as serverUtils

import src.network_interface.ListContents.ListContents_pb2 as ListContents_pb2
import src.network_interface.ListContents.ListContents_pb2_grpc as ListContents_pb2_grpc

import src.network_interface.BuildTest.BuildTest_pb2 as BuildTest_pb2
import src.network_interface.BuildTest.BuildTest_pb2_grpc as BuildTest_pb2_grpc

# Entry-point for the command line interface. Appears as 'packagebird'.
@click.group()
@click.pass_context
def cli(ctx):
    """Client utility for managing packages, is a group command for the subcommands below. Call --help on a subcommand to get further details."""
    # The entry point for the command line interface
    # Should ping server from this point to check that there is a server to connect with at the specified point
    if not serverUtils.ServerUtils.ping('127.0.0.1', '50051'):
        exit(0)

# Sync project contents with remote server
@cli.command('sync', short_help='Syncs source of project with server')
@click.pass_context
def sync(ctx):
    # Close if called on directory not formatted as project
    if not FilesystemInterface.check_if_project_dir():
        click.echo('Directory is not configured as project. Please either configure as project or navigate to configured directory.')
        return

    pr_config = FilesystemInterface.get_project_config()
    pr_name = pr_config["name"]

    click.echo('Syncing project source files')

    # Package name
    project_name = f'{pr_name}'

    click.echo(f'Uploaded project: {project_name} source')

    # Create compressed archive of file contents
    FilesystemInterface.make_archive(project_name)

    project_archive_name = f'{project_name}.tar.gz'

    # Upload to server
    fileservice = FileTransfer()
    fileservice.upload('127.0.0.1', '50051', project_archive_name, "project")
    os.remove(project_archive_name)

# Add package to the development directory
@cli.command('addpackage', short_help='Adds a package to development directory packages subdirectory')
@click.option('-p', '--package', 'name')
@click.option('-v', '--version', 'version')
@click.pass_context
def addpackage(ctx, name, version):
    # bad idea approaching...
    # get project name
    pr_config = FilesystemInterface.get_project_config()
    pr_name = pr_config["name"]

    # Makes the packages directory
    if not os.path.isdir('packages'):
        click.echo("Creating packages directory...")
        os.mkdir('packages')
    os.chdir('packages')

    # Download the archived package
    fileservice = FileTransfer()

    # add project name here
    request_string = f'{pr_name}-v{name}-v{version}'

    # Make a request for a package-list
    # grpc packagelist service here
    stub = packageOperationsGRPC.PackageOperationServicesStub(grpc.insecure_channel('127.0.0.1:50051'))
    packageList = stub.GetPackageList(packageOperations.PackageRequest(packageitem=request_string))

    # Debug output
    # print(f'Response from server: {packageList.packageitem}')
    package_request = f'{name}-v{version}'

    # Add the originally requested package
    packageList.packageitem.append(package_request)
    
    # Iterate through list of returned packages
    for packageListItem in packageList.packageitem:
        pstring = packageListItem.split('-v')
        itemName = pstring[0]
        itemVersion = pstring[1]
    
        itemRequestString = f'{itemName}-v{itemVersion}.tar.gz'

        # Create Package Directory
        os.mkdir(f'{itemName}')
        os.chdir(f'{itemName}')
    
        # Download the package source into the package directory
        fileservice.download('127.0.0.1', '50051', f'{itemRequestString}')
    
        # Extract the contents
        with tarfile.open(itemRequestString, 'r:gz') as archive:
            try:
                archive.extractall(path='.')
            except tarfile.ReadError:
                click.echo('Empty file found.')
    
        # Remove the archive file
        os.remove(itemRequestString)

        # Move back to package directory
        os.chdir('..')
    
    # change back into project root directory
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
    package_name = f'{project_name}-{project_version}'

    # Create compressed archive of file contents
    FilesystemInterface.make_archive(package_name)

    package_name = f'{package_name}.tar.gz'

    # Upload to server
    if (not debug):
        fileservice = FileTransfer()
        fileservice.upload('127.0.0.1', '50051', package_name, "package")
    os.remove(package_name)

# Builds package on the server
@cli.command('buildpackage', short_help='Builds a package on the server and relays the result')
@click.option('-n', '--name', help='Name of package being built')
@click.option('-v', '--version', help='Version of package being built')
@click.pass_context
def buildpackage(ctx, name, version):
    # Requested package name and version
    version = int(version)
    request = BuildTest_pb2.PackageInfo(name=name, version=version)

    with grpc.insecure_channel('localhost:50051') as channel:
        stub = BuildTest_pb2_grpc.BuildTestServicesStub(channel)
        response = stub.Build(request)
        if response != None:
            click.echo(response)
    pass

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

# List entites in server
@cli.command('list', short_help='List Projects, Packages, Members associated with configured Packagebird server')
@click.option('-pro', '--projects', help='All projects registered on server', is_flag=True)
@click.option('-pac', '--packages', help='All packages registered on server', is_flag=True)
@click.option('-mem', '--members', help='All members registered on server', is_flag=True)
@click.pass_context
def gList(ctx, projects, packages, members):
    # Using Click's flag options
    getProjects, getPackages, getMembers = False, False, False
    if projects:
        getProjects = True
    if packages:
        getPackages = True
    if members:
        getMembers = True

    # Exit early if nothing
    if not getProjects and not getPackages and not getMembers:
        return  

    listRequest = ListContents_pb2.ContentRequest(ListProjects=getProjects, ListPackages=getPackages, ListMembers=getMembers)
    listResponse = {}

    #DEBUG
    # click.echo(f"{listRequest.ListProjects}\n{listRequest.ListPackages}\n{listRequest.ListMembers}")

    with grpc.insecure_channel('localhost:50051') as channel:
        stub = ListContents_pb2_grpc.ListContentServicesStub(channel)
        listResponse = stub.GetContent(listRequest)
    
    hasProjects, hasPackages, hasMembers = len(listResponse.Projects)>0, len(listResponse.Packages)>0, len(listResponse.Members)>0

    if hasProjects:
        click.echo(f"Projects:")
        for project in listResponse.Projects:
            click.echo(f"\t{project}")

    if hasPackages:
        click.echo(f"Packages:")
        for package in listResponse.Packages:
            click.echo(f"\t{package}")

    if hasMembers:
        click.echo(f"Members:")
        for member in listResponse.Members:
            click.echo(f"\t{member}")

    pass