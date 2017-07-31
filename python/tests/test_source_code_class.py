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
IF_ELIF_ELSE_EXAMPLE = """
if a is not None:
    b = a * 2
    print("a was something. b ==", b)
elif b is not None:
    b *= 2
    print("b was something. b ==", b)
else:
    b = 42
    print("a was nothing. b ==", b)
"""
NESTED_IF_ELSE_EXAMPLE = """
if b:
    a = b + 2
    print("a:", a)
    if a % 2 == 0:
        print("EVEN!")
    else:
        print("ODD!")
else:
    print("No b")
    c = 2 + 7
    if c % a == 0:
        print('yes')
    else:
        # comment
        print('no')
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
        self.assertEqual(self.result, self.expected)


class TestSourceCodeContextSearchLineNo(SourceCodeTestCase):
    matcher_type = LineNoMatcher

    def test_if_else_block_shows_if_line(self):
        self.expected = ["if a is not None:\n",
                         '    print("a was something. b ==", b)\n']
        self.assert_context_accurate(FakeFile(IF_ELSE_EXAMPLE), 4)

    @unittest.expectedFailure
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

    @unittest.expectedFailure
    def test_try_except_else_block_shows_try_except_and_else_lines(self):
        self.expected = ['try:\n',
                         'except Exception as e:\n',
                         'else:\n',
                         '    print("Ran successfully")\n']
        self.assert_context_accurate(FakeFile(TRY_EXCEPT_ELSE_EXAMPLE), 12)

    def test_if_elif_else_block_shows_elif_line(self):
        self.expected = ['if a is not None:\n',
                         'elif b is not None:\n',
                         '    print("b was something. b ==", b)\n']
        self.assert_context_accurate(FakeFile(IF_ELIF_ELSE_EXAMPLE), 7)

    @unittest.expectedFailure
    def test_if_elif_else_block_shows_else_line(self):
        self.expected = ['if a is not None:\n',
                         'elif b is not None:\n',
                         'else:\n',
                         '    b = 42\n']
        self.assert_context_accurate(FakeFile(IF_ELIF_ELSE_EXAMPLE), 9)

    def test_nested_if_else_with_only_if_branches(self):
        self.expected = ['if b:\n',
                         '    if a % 2 == 0:\n',
                         '        print("EVEN!")\n']
        self.assert_context_accurate(FakeFile(NESTED_IF_ELSE_EXAMPLE), 6)

    @unittest.expectedFailure
    def test_nested_if_else_with_else_branch(self):
        self.expected = ['if b:\n',
                         '    if a % 2 == 0:\n',
                         '    else:\n',
                         '        print("ODD!")\n']
        self.assert_context_accurate(FakeFile(NESTED_IF_ELSE_EXAMPLE), 8)

    @unittest.expectedFailure
    def test_nested_if_else_with_both_else_branches(self):
        self.expected = ['if b:\n',
                         'else:\n',
                         '    else:\n',
                         "        print('no')\n"]
        self.assert_context_accurate(FakeFile(NESTED_IF_ELSE_EXAMPLE), 16)


if __name__ == '__main__':
    unittest.main()
