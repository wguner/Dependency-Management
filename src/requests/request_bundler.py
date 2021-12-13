import json
from typing import List
import json
import requests

class make_request:

    """Bundle request together"""
    def __init__(self, address, dependencies: List) -> bool:
        request_data = {'dependencies': dependencies}
        request = requests.post(address, request_data)
        return request.ok