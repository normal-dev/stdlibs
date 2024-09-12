import sys
import unittest
import os
import extractor

def open_test(test):
    f = open("./tests/data/" + test)
    s = f.read()
    f.close()
    return s

class Default(unittest.TestCase):
    def test_global(self):
        f = open_test("default/global/?.py")
        locus = extractor.extract(f)

        print(locus)

if __name__ == '__main__':
    unittest.main()