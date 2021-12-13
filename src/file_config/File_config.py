import os
import sys
import click

"""Handles configuration and local environment modification"""
class File_config:
    def __init__(self) -> None:
        pass

    """Creates a directory with a name"""
    def create_dir(file_name: str) -> bool:
        os.mkdir(file_name)
        return True
    
    """Creates a package formatted directory"""
    def create_package_dir() -> bool:
        if os.path.isdir('packages'):
            return False
        else:
            os.mkdir('packages')
            return True
    
    """Creates a package formatted and bundles properly"""
    def create_package(name: str) -> bool:
        if os.path.isdir(f'packages/{name}'):
            click.echo(f"\nPackage '{name}'' already present!...")
            return

        if File_config.create_package_dir():
            click.echo("\nCreating packages directory...")
        
        if not os.path.isdir(f'packages/{name}'):
            os.mkdir(f'packages/{name}')
            os.mkdir(f'packages/{name}/src')

        filename = f'{name}.pkg'
        manifest = 'manifest.json'
        src_file = 'src.java'        

        with open(os.path.join(f'packages/{name}', filename), 'w') as package:
            package.write(f'Fill in later...')
        
        with open(os.path.join(f'packages/{name}', manifest), 'w') as manifest:
            manifest.write('{"name":"'+name+'"}')
        
        with open(os.path.join(f'packages/{name}/src', src_file), 'w') as file:
            file.write("class src_file { public static void main(String args) { /* Stub */ } }")

        
"""Extra miscellanious"""
if __name__=='__main__':
    File_config.create_package('Sample')
    File_config.create_package('Meta')