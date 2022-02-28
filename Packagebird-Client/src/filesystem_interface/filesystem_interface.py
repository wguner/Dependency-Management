from distutils.log import error
import os
import json
import tarfile
from typing import Dict

from click import File

class FilesystemInterface(object):

    # Checks if a directory exist at location
    @staticmethod
    def check_if_dir_exist(name: str) -> bool:
        return os.path.isdir(name)

    # Creates a directory if not already present
    @staticmethod
    def create_dir(name: str) -> bool:
        if not FilesystemInterface.check_if_dir_exist(name):
            os.mkdir(name)
            return True
        return False
    
    # Checks if a 'projectconfig.json'
    @staticmethod
    def check_if_project_dir() -> bool:
        return os.path.exists('projectconfig.json')
    
    # Gets the contents of a '.json' formatted file as a dictionary
    @staticmethod
    def get_project_config() -> Dict:
        pcFile = open('projectconfig.json')
        
        project_config = json.load(pcFile)

        pcFile.close()

        return project_config
    
    # Really crazy way of making a package from a particular directory and also filtering out excluding files
    @staticmethod
    def make_archive(projectName, projectVersion):
        package_name = f'{projectName}-v{projectVersion}.tar.gz'

        with tarfile.open(package_name, 'w:gz', format=tarfile.GNU_FORMAT) as tar:
            for dir, dirs, files in os.walk("."):
                for file in files:
                    # Special check for zip-slip vulnerability
                    if '/packages' not in dir and '\\packages\\' not in dir and './packages/' not in dir and '/../' not in dir:
                        if file != package_name and file != 'projectconfig.json':
                            tar.add(os.path.join(dir,file))
    
    # Project source alternative
    @staticmethod
    def make_archive(projectName):
        fileName = f'{projectName}.tar.gz'
        with tarfile.open(fileName, 'w:gz', format=tarfile.GNU_FORMAT) as tar:
            for dir, dirs, files in os.walk("."):
                for file in files:
                    # Special check for zip-slip vulnerability
                    if '/packages' not in dir and '\\packages\\' not in dir and './packages/' not in dir and '/../' not in dir:
                        if file != fileName and file != 'projectconfig.json':
                            tar.add(os.path.join(dir,file)),

    # Adds appropriate subdirectories to freshly created project directory
    @staticmethod
    def make_project_dir(projectName, projectVersion):
        
        # Changes into project directory
        os.chdir(projectName)

        # Creates the packages subdirectory
        os.mkdir('packages')

        # Creates the project config JSON file
        config = {
            "name": projectName,
            "version": projectVersion
        }

        with open('projectconfig.json', 'w') as projectConfig:
            json.dump(config, projectConfig)
        
        # Return to parent directory
        os.chdir('..')