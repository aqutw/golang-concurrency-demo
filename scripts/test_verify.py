"""Regression checks for executable example coverage."""
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

import verify


class CoverageTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        self.patch = patch.object(verify, 'ROOT', self.root)
        self.patch.start()
        self.addCleanup(self.patch.stop)
        self.code = 'package main\n\nfunc main() {}\n'
        self.example = self.root / 'examples/basic/example/main.go'
        self.example.parent.mkdir(parents=True)
        self.example.write_text(self.code)
        self.entries = [{
            'path': 'examples/basic/example/main.go',
            'kind': 'success', 'contains': ['expected output'], 'counts': {},
        }]

    def test_complete_manifest_matches_public_examples(self):
        self.assertEqual(verify.coverage(self.entries), 1)

    def test_unlisted_executable_fails(self):
        (self.example.parent / 'extra.go').write_text('package main\n')
        with self.assertRaisesRegex(AssertionError, 'Go files and manifest do not match'):
            verify.coverage(self.entries)

    def test_missing_executable_fails(self):
        self.example.unlink()
        with self.assertRaisesRegex(AssertionError, 'Go files and manifest do not match'):
            verify.coverage(self.entries)

    def test_duplicate_paths_fail(self):
        with self.assertRaisesRegex(AssertionError, 'Duplicate manifest paths'):
            verify.coverage(self.entries * 2)

    def test_missing_behavior_assertions_fail(self):
        self.entries[0]['contains'] = []
        with self.assertRaisesRegex(AssertionError, 'Missing behavior assertions'):
            verify.coverage(self.entries)


if __name__ == '__main__':
    unittest.main()
