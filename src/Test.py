from typing import Dict, List
import unittest
import string

import requests
from manifest import Manifest
from file_config import File_config
from request_bundler import make_request

class TestManifest(unittest.TestCase):
    sample_manifest_path: str = 'src/sample_manifest.json'
    sample_manifest_file: str = '{"dependencies":{"name":"Sample","version":"0.0.1"},"registry":"http://localhost:9091/registry"}'
    sample_manifest_parse: Dict = {
        "dependencies": [
        {
            "name": "Sample",
            "version": "0.0.1"
        }
        ],
        "registry": "http://localhost:9091/registry"
    }
    sample_manifest_dependencies: List = ["Sample-v0.0.1"]
    sample_manifest_request: str = '{"requests:["Sample-v0.0.1"]"}'

    """Testing conversion between JSON file and Dictionary"""
    def test_convert(self):
        self.assertEqual(self.sample_manifest_parse, Manifest.file_parse(self, self.sample_manifest_path), "Files should match.")
    
    """Testing class conversion between JSON file and class dictionary"""
    def test_class_convert(self):
        sample_manifest = Manifest(self.sample_manifest_path)
        self.assertEqual(self.sample_manifest_parse, sample_manifest.manifest, "Class attribute should match sample")
    
    """Testing dependency list results"""
    def test_class_dependencies(self):
        self.assertEqual(self.sample_manifest_dependencies, Manifest(self.sample_manifest_path).dependencies(), "Dependency list should match")
    
    """Testing request and return results"""
    def test_class_request(self):
        pass
        # sample_manifest = Manifest(self.sample_manifest_path)
        # self.assertRaises(Exception,make_request(sample_manifest.registry, sample_manifest.dependencies))

class TestFileConfig(unittest.TestCase):
    sample_name = 'Sample'

    """Testing file configuration"""
    def test_file_config(self):
        pass

if __name__ == '__main__':
    unittest.main()