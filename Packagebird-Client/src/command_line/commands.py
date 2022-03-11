from ast import Str
from distutils.version import Version
from inspect import getmembers
from posixpath import split
from re import M, T
from tokenize import Name
from urllib import request
from xmlrpc.client import Boolean
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

# Member operations
import src.network_interface.member_operations.MemberOperations_pb2 as MemberOperations_pb2
import src.network_interface.member_operations.MemberOperations_pb2_grpc as MemberOperations_pb2_grpc

# Global variable for connection address
address =  "localhost" # "149.28.65.7" # "localhost"

# Entry-point for the command line interface. Appears as 'packagebird'.
@click.group()
@click.pass_context
def cli(ctx):
    """Client utility for managing packages, is a group command for the subcommands below. Call --help on a subcommand to get further details."""
    # The entry point for the command line interface
    # Should ping server from this point to check that there is a server to connect with at the specified point
    address = os.getenv("PACKAGEBIRD_ADDRESS")
    if address == None:
        address = "localhost"
        click.echo("Server connection address not set, defaulting to localhost. Please use the 'config' subcommand to set the appropriate connection address.")
    if not serverUtils.ServerUtils.ping(address, '50051'):
        exit(0)

# Configurate values on the client, first being the server connection string
@cli.command('config', short_help='Configurate local variables, such as connection address and other settings')
@click.option('-a', '--address', 'address', default=None, help='Sets the environment variable to the connection server address')
@click.option('-u', '--username', 'username', default=None, help='Sets the client configured username')
@click.option('p', '--password', 'userpassword', default=None, help='Sets the client configured password')
@click.pass_context
def config(ctx, address, username, userpassword):
    if address != None:
        click.echo(f"Setting Packagebird server connection address to {address}")
        os.environ["PACKAGEBIRD_ADDRESS"] = address
    if username != None: 
        click.echo(f"Setting client username to {username}")
        os.environ["PACKAGEBIRD_USERNAME"] = username
    if userpassword != None:
        click.echo(f"Setting client password to {userpassword}")
        os.environ["PACKAGEBIRD_USERPASSWORD"] = userpassword

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
    fileservice.upload(address, '50051', project_archive_name, "project")
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
    
    # Get project dependencies
    pr_dependencies = pr_config["dependencies"]
    for dep in pr_dependencies:
        if name==dep["name"] and version==dep["version"]:
            click.echo(f'Package "{name}" with version {version} already a dependency of the project.\nCancelling operation.')
            return
    
    # Check if server if online
    if not serverUtils.ServerUtils.ping(f'{address}', '50051'):
        click.echo(f'Cannot connect to specified server. Please check network connection and configured server address and port.')
        return

    # Make a request to the server
    # If package wasn't found, server will return signal to use and empty list
    # If package was found, server will return list of depedencies to fetch and install
    # For every package in list, retreive archived package source from the server and configure in the project
    
    # Function should be called for every dependency


    # Makes the packages directory if not already present
    # Should be part of different function for project configuration
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
    stub = packageOperationsGRPC.PackageOperationServicesStub(grpc.insecure_channel(f'{address}:50051'))
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
        fileservice.download(address, '50051', f'{itemRequestString}')
    
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
        fileservice.upload(address, '50051', package_name, "package")
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

    with grpc.insecure_channel(f'{address}:50051') as channel:
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
    with grpc.insecure_channel(f'{address}:50051') as channel:
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
    # Set all three flags true if no flags passed
    if not projects and not packages and not members:
        getProjects = True
        getPackages = True
        getMembers = True
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

    with grpc.insecure_channel(f'{address}:50051') as channel:
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

# Add member to server
@cli.command('addmember', short_help='Creates a new member on the configured server.\nCAN ONLY BE PERFORMED BY A USER WITH ADMIN LEVEL PRIVILEGES.')
@click.option('-n', '--name', 'name', default=None, required=True, help='The name of the new member being created')
@click.option('-p', '--password', 'password', default=None, required=True, help='The password of the user being created')
@click.option('-l', '--level', 'level', default=None, required=True, help='The privilege level of the user being created')
@click.pass_context
def addmember(ctx, name, password, level):

    # Ensure all parameters are specified, all flags used
    if name is None:
        click.echo(f'Name of user must be specified! Please provide a name of the user being created.')
        return
    if password is None:
        click.echo(f'Password for the user must be specified! Please provide a password for the user being created.')
        return
    if level is None:
        click.echo(f'Privilege-level for the user must be specified! Please provide a privilege-level for the user being created.')
        return

    # Make request to server, get message and relay to user
    # Items needed for request
    # currently configured client username, client password
    # new user username, new user password, new user level
    # on server, check if passed client username, client password match what is specified in database
    # if so, check if user has appropriate privilege levels
    # if so, proceed
    if not ClientAuth():
        return
    
    if name is 'admin':
        click.echo('Cannot add new admin user to server')
        return
            
    with grpc.insecure_channel(f'{address}:50051') as channel:
        memberCRUDStub = MemberOperations_pb2_grpc.MemberCRUDServicesStub(channel)
        createMemberRequest = MemberOperations_pb2.MemberAuthenticationRequest(Name=name, Password=password, Level=level)
        memberResponse = memberCRUDStub.CreateMember(createMemberRequest)
        click.echo(memberResponse.Response)

# Remove member from server
@cli.command('removemember', short_help='Removes specified member from server.\nCAN ONLY BE PERFORMED BY A USER WITH ADMIN LEVEL PRIVILEGES.')
@click.option('-n', '--name', 'name', default=None, required=True, help='The name of the member being removed')
@click.pass_context
def removemember(ctx, name):

    # Ensure parameter is specified
    if name is None:
        click.echo(f'Name of member must be specified! Please specify the member name.')
        return
    
    if not ClientAuth():
        return
    
    if name is 'admin':
        click.echo('Cannot remove admin user from server')
        return
    
    with grpc.insecure_channel(f'{address}:50051') as channel:
        memberCRUDStub = MemberOperations_pb2_grpc.MemberCRUDServicesStub(channel)
        createMemberRequest = MemberOperations_pb2.MemberAuthenticationRequest(Name=name)
        memberResponse = memberCRUDStub.CreateMember(createMemberRequest)
        click.echo(memberResponse.Response)

def ClientAuth() -> Boolean:
    # Make request to server, get message and relay to user
    with grpc.insecure_channel(f'{address}:50051') as channel:
        memberAUTHStub = MemberOperations_pb2_grpc.MemberAuthenticationStub(channel)
        
        memberAuthRequest = MemberOperations_pb2.MemberAuthenticationRequest(Name=os.getenv("PACKAGEBIRD_USERNAME"), Password=os.getenv("PACKAGEBIRD_USERPASSWORD", IsLevel=True))
        memberAuthResponse = memberAUTHStub.AuthenticateMember(memberAuthRequest)
        if memberAuthResponse.Is_Authenticated == False:
            click.echo(memberAuthResponse.Response)
            return False
        return True
