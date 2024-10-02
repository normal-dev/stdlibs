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
            # {
            #     "file": "global/Expr.py",
            #     "expected": [{
            #         "ident": "types.CodeType",
            #         "line": 3
            #     }]
            # },
            {
                "file": "global/?.py",
                "expected": [{
                    "ident": "datetime.timedelta",
                    "line": 3
                }]
            }
        ]

        for test in tests:
            src = open_test(test["file"])
            actual = extractor.extract(src)
            expected = test["expected"]

            self.assertListEqual(actual, expected)

    # def test_def(self):
    #     tests = [
    #         {
    #             "file": "def/Expr.py",
    #             "expected": [{
    #                 "ident": "sys.api_version",
    #                 "line": 4
    #             }]
    #         }
    #     ]

    #     for test in tests:
    #         src = open_test(test["file"])
    #         actual = extractor.extract(src)
    #         expected = test["expected"]

    #         self.assertListEqual(actual, expected)

if __name__ == '__main__':
    unittest.main()