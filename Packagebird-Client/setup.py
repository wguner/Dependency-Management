from importlib.metadata import entry_points
from setuptools import setup, find_packages

setup(
    name='packagebird',
    version='0.0.5',
    packages=find_packages(),
    include_package_data=True,
    install_requires=[
        'Click',
    ],
    entry_points={
        'console_scripts': [
            'packagebird=src.command_line.commands:cli'
        ]
    }
)