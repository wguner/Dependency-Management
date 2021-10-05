import argparse

# Server address
server_adr = 'localhost'

# Contact the server, with custom address
parser = argparse.ArgumentParser(description='Contacts the token server.')
parser.add_argument('address', metavar=server_adr, type=str, nargs='?', help='The address of the server to conact')

# More code below