import os
import sys
import unittest
from io import StringIO

from helpers import *
from context import LineNoMatcher, Matcher, RegexMatcher, SourceCode


IF_ELSE_EXAMPLE = """
if a is not None:
    b = a * 2
    print("a was something. b ==", b)
else:
    b = 42
    print("a was nothing. b ==", b)
"""
TRY_EXCEPT_ELSE_EXAMPLE = """
try:
    a = 7
    b = a ** 3
    c = a * b
    print("7 ** 4 ==", c)
except Exception as e:
    c = None
    print("Unexpected error:", e)
    print("c == ", str(c))
else:
    print("Ran successfully")
"""


class FakeFile(StringIO):
    name = "test"


class SourceCodeTestCase(unittest.TestCase):
    matcher_type = Matcher

    def assert_context_accurate(self, test_file, look_for):
        self.source = SourceCode(test_file, self.matcher_type, look_for, filename=test_file.name)
        self.result = self.source.find_context()
        print("\nRESULT:\n{}".format(self.result))
        self.assertEqual(self.result, self.expected)


class TestSourceCodeContextSearchLineNo(SourceCodeTestCase):
    matcher_type = LineNoMatcher

    def test_if_else_block_shows_if_line(self):
        self.expected = ["if a is not None:\n",
                         '    print("a was something. b ==", b)\n']
        self.assert_context_accurate(FakeFile(IF_ELSE_EXAMPLE), 4)

    def test_if_else_block_shows_if_and_else_lines(self):
        self.expected = ["if a is not None:\n",
                         "else:\n",
                         "    b = 42\n"]
        self.assert_context_accurate(FakeFile(IF_ELSE_EXAMPLE), 6)

    def test_try_except_else_block_shows_try_line(self):
        self.expected = ["try:\n",
                         '    print("7 ** 4 ==", c)\n']
        self.assert_context_accurate(FakeFile(TRY_EXCEPT_ELSE_EXAMPLE), 6)

    def test_try_except_else_block_shows_try_and_except_lines(self):
        self.expected = ['try:\n',
                         'except Exception as e:\n',
                         '    c = None\n']
        self.assert_context_accurate(FakeFile(TRY_EXCEPT_ELSE_EXAMPLE), 8)

    def test_try_except_else_block_shows_try_except_and_else_lines(self):
        self.expected = ['try:\n',
                         'except Exception as e:\n',
                         'else:\n',
                         '    print("Ran successfully")\n']
        self.assert_context_accurate(FakeFile(TRY_EXCEPT_ELSE_EXAMPLE), 12)


if __name__ == '__main__':
    unittest.main()
