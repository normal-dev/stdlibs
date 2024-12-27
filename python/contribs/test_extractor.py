import sys
import unittest
import os
import extractor

def open_test(test):
    f = open("./tests/" + test)
    s = f.read()
    f.close()
    return s

class TestExtractor(unittest.TestCase):
    def test_global(self):
        tests = [
            {
                "file": "Assign.py",
                "expected": [
                    {
                        "ident": "sys.abiflags",
                        "line": 3
                    },
                    {
                        "ident": "sys.base_prefix",
                        "line": 5
                    },
                    {
                        "ident": "sys.base_prefix",
                        "line": 7
                    }
                ]
            },
            {
                "file": "Attribute.py",
                "expected": [
                    {
                        "ident": "types.CodeType",
                        "line": 7
                    },
                    {
                        "ident": "sys.api_version",
                        "line": 5
                    },
                ]
            },
            {
                "file": "Call.py",
                "expected": [
                    {
                        "ident": "datetime.timedelta",
                        "line": 3
                    }
                ]
            },
            {
                "file": "Import.py",
                 "expected": [
                    {
                        "ident": "sys.stdlib_module_names",
                        "line": 5
                    },
                    {
                        "ident": "ast",
                        "line": 6
                    }
                ]
            },
            {
                "file": "ImportFrom.py",
                 "expected": [
                    {
                        "ident": "collections.abc",
                        "line": 6
                    },
                    {
                        "ident": "collections.abc",
                        "line": 7
                    },
                    {
                        "ident": "typing.Any",
                        "line": 8
                    },
                    {
                        "ident": "ast.parse",
                        "line": 9
                    },
                    {
                        "ident": "datetime.datetime",
                        "line": 10
                    }
                ]
            }
        ]

        for test in tests:
            src = open_test(test["file"])
            actual = extractor.extract(src)
            expected = test["expected"]

            self.assertListEqual(actual, expected)

if __name__ == '__main__':
    unittest.main()