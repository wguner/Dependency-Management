import argparse
import requests

# Server address
server_adr = 'http://localhost:8080'

# Contact the server, with custom address
parser = argparse.ArgumentParser(description='Contacts the token server.')
parser.add_argument('address', metavar='server_adr', type=str, nargs='?', help='The address of the server to conact')

ars = parser.parse_args()

# Connect to server endpoint, await response
request = requests.get(server_adr)

# Get status
status = request.status_code

# Parse to JSON
response = request.json()

# Print status and payload
print(f'Status Code: {status}\nPayload:\n{response}')