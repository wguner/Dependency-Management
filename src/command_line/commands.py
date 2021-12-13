import click
import os
import sys
import tarfile
from src.file_config import File_config
from src.manifest.manifest import Manifest
from src.requests.request_stub import Request_Stub
from src.fileserver.fileserver_client import FileTransfer
from src.command_line.utils import Utils

@click.group()
def cli():
    pass

@cli.command()
@click.option('-m', '--manifest', 'manifest')
@click.pass_context
def init(ctx, manifest):
    if manifest == None:
        if not os.path.exists('manifest.json'):
            click.echo("'manifest.json' file is not found! Please navigate to directory containing the manifest file.")
            return
        local_manifest = Manifest('manifest.json')
    else:
        if not os.path.exists(manifest):
            click.echo("Specified manifest file is not found! Check spelling or navigate to directory.")
            return
        local_manifest = Manifest(manifest)
    ctx.invoke(request, manifest=local_manifest, registry='http:///localhost:4040/')
    pass

@cli.command()
@click.option('-p', '--project', 'project')
@click.pass_context
def request(ctx, project):
    requester = Request_Stub()
    response = requester.make_request('http://localhost:4040/',project)
    with click.progressbar(length=len(response), label='Loading packages...') as bar:
        for package in response:
            ctx.invoke(addpackage, name=package, version=1)
            bar.update(1)
    pass

@cli.command()
@click.option('-p', '--package', 'name')
@click.option('-v', '--version', 'version')
@click.pass_context
def addpackage(ctx, name, version):
    # Makes the packages directory
    if not os.path.isdir('packages'):
        click.echo("Creating packages directory...")
        os.mkdir('packages')
    os.chdir('packages')

    fileservice = FileTransfer()
    request_string = f'{name}-v{version}.tar.gz'
    fileservice.download('127.0.0.1', '50051', request_string)

    #package_path = f'packages/{request_string}'
    #click.echo(f'Moving {request_string} to {package_path}.')
    #os.rename(request_string, package_path)

    with tarfile.open(request_string, 'r:gz') as archive:
        archive.extractall(path='.')
    os.remove(request_string)
    os.rename(f'{name}-v{version}', f'{name}')
    os.chdir('..')


@cli.command()
@click.option('-p', '--project', 'project')
@click.pass_context
def setup(ctx, project):
    fileservice = FileTransfer()
    utils = Utils()
    packages_strings = utils.project_packages(project)
    for package_addr in packages_strings:
        package_name = package_addr[0]
        package_version = package_addr[1]
        ctx.invoke(addpackage, name=package_name, version=package_version)