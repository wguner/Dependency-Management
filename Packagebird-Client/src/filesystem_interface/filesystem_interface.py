from distutils.log import error
import os

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
