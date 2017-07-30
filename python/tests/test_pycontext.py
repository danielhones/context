import os
import sys
import unittest
if sys.version_info.major == 2:
    from StringIO import StringIO
else:
    from io import StringIO

from helpers import *
from context import BLUE, END_COLOR, insert_term_color, SEARCH_LINENO, SEARCH_REGEX


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


class TestInsertTermColor(unittest.TestCase):
    def test_color_whole_string(self):
        teststr = "some example test string\nwith some new lines\tand tabs."
        colored = insert_term_color(teststr, 0, len(teststr), BLUE)
        expected = BLUE + teststr + END_COLOR
        self.assertEqual(colored, expected)

    def test_color_substring(self):
        teststr = "only color THIS word"
        colored = insert_term_color(teststr, 11, 15, BLUE)
        expected = "only color " + BLUE + "THIS" + END_COLOR + " word"
        self.assertEqual(colored, expected)


class TestPycontextMainSearchLineNo(unittest.TestCase):
    """Functional tests to check accuracy of output from main"""

    def setUp(self):
        self.out = StringIO()
        self.errout = StringIO()

    def tearDown(self):
        self.print_errors()

    def print_errors(self):
        self.errout.seek(0)
        print("\nErrors:\n{}".format(self.errout.read()))

    def test_if_else_block_shows_if_line(self):
        test_source = FakeFile(IF_ELSE_EXAMPLE)
        result = context.main(4, test_source, search_type=SEARCH_LINENO,
                              verbose=True, output=self.out, errout=self.errout)
        expected = ["if a is not None:\n",
                    '    print("a was something. b ==", b)\n']
        self.assertEqual(result[test_source.name], expected)

    def test_if_else_block_shows_if_and_else_lines(self):
        test_source = FakeFile(IF_ELSE_EXAMPLE)
        result = context.main(6, test_source, search_type=SEARCH_LINENO,
                              verbose=True, output=self.out, errout=self.errout)
        expected = ["if a is not None:\n",
                    "else:\n",
                    "    b = 42\n"]
        self.assertEqual(result[test_source.name], expected)

    def test_try_except_else_block_shows_try_line(self):
        test_source = FakeFile(TRY_EXCEPT_ELSE_EXAMPLE)
        result = context.main(6, test_source, search_type=SEARCH_LINENO,
                              verbose=True, output=self.out, errout=self.errout)
        expected = ["try:\n",
                    '    print("7 ** 4 ==", c)\n']
        self.assertEqual(result[test_source.name], expected)

    def test_try_except_else_block_shows_try_and_except_lines(self):
        test_source = FakeFile(TRY_EXCEPT_ELSE_EXAMPLE)
        result = context.main(8, test_source, search_type=SEARCH_LINENO,
                              verbose=True, output=self.out, errout=self.errout)
        expected = ['try:\n',
                    'except Exception as e:\n',
                    '    c = None\n']
        self.assertEqual(result[test_source.name], expected)

    def test_try_except_else_block_shows_try_except_and_else_lines(self):
        test_source = FakeFile(TRY_EXCEPT_ELSE_EXAMPLE)
        result = context.main(12, test_source, search_type=SEARCH_LINENO,
                              verbose=True, output=self.out, errout=self.errout)
        expected = ['try:\n',
                    'except Exception as e:\n',
                    'else:\n',
                    '    print("Ran successfully")\n']
        self.assertEqual(result[test_source.name], expected)


if __name__ == "__main__":
    unittest.main()
