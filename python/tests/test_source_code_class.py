import os
import sys
import unittest
from io import StringIO

from helpers import *


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


class TestSourceCodeContextSearchLineNo(unittest.TestCase):
    def test_if_else_block_shows_if_line(self):
        test_file = FakeFile(IF_ELSE_EXAMPLE)
        source = context.SourceCode(test_file, context.LineNoMatcher, 4, filename=test_file.name)
        result = source.find_context()
        expected = ["if a is not None:\n",
                    '    print("a was something. b ==", b)\n']
        self.assertEqual(result, expected)

    def test_if_else_block_shows_if_and_else_lines(self):
        test_file = FakeFile(IF_ELSE_EXAMPLE)
        source = context.SourceCode(test_file, context.LineNoMatcher, 6, filename=test_file.name)
        result = source.find_context()
        expected = ["if a is not None:\n",
                    "else:\n",
                    "    b = 42\n"]
        self.assertEqual(result, expected)

    def test_try_except_else_block_shows_try_line(self):
        test_file = FakeFile(TRY_EXCEPT_ELSE_EXAMPLE)
        source = context.SourceCode(test_file, context.LineNoMatcher, 6, filename=test_file.name)
        result = source.find_context()
        expected = ["try:\n",
                    '    print("7 ** 4 ==", c)\n']


    def test_try_except_else_block_shows_try_and_except_lines(self):
        test_file = FakeFile(TRY_EXCEPT_ELSE_EXAMPLE)
        source = context.SourceCode(test_file, context.LineNoMatcher, 8, filename=test_file.name)
        result = source.find_context()
        expected = ['try:\n',
                    'except Exception as e:\n',
                    '    c = None\n']
        self.assertEqual(result, expected)

    def test_try_except_else_block_shows_try_except_and_else_lines(self):
        test_file = FakeFile(TRY_EXCEPT_ELSE_EXAMPLE)
        source = context.SourceCode(test_file, context.LineNoMatcher, 12, filename=test_file.name)
        result = source.find_context()
        expected = ['try:\n',
                    'except Exception as e:\n',
                    'else:\n',
                    '    print("Ran successfully")\n']
        self.assertEqual(result, expected)
