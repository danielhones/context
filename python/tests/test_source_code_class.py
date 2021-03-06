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
        if False:
            print("that would be weird")
        elif None:
            print("so would this")
        else:
            x = 27
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
finally:
    print("running cleanup")
    sock.close()
"""
MIXED_TRY_IF_ELSE_EXAMPLE = """

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

    def test_if_without_else(self):
        test_file = FakeFile("""
if True:
    a = 2
    print(a)
""")
        self.expected = ['if True:\n',
                         '    print(a)\n']
        self.assert_context_accurate(test_file, 4)

    def test_if_else_block_shows_if_line(self):
        self.expected = ["if a is not None:\n",
                         '    print("a was something. b ==", b)\n']
        self.assert_context_accurate(FakeFile(IF_ELSE_EXAMPLE), 4)

    def test_if_else_block_shows_if_and_else_lines(self):
        self.expected = ["if a is not None:\n",
                         "else:\n",
                         "    b = 42\n"]
        self.assert_context_accurate(FakeFile(IF_ELSE_EXAMPLE), 6)

    def test_if_elif_else_block_shows_elif_line(self):
        self.expected = ['if a is not None:\n',
                         'elif b is not None:\n',
                         '    print("b was something. b ==", b)\n']
        self.assert_context_accurate(FakeFile(IF_ELIF_ELSE_EXAMPLE), 7)

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

    def test_nested_if_else_with_else_branch(self):
        self.expected = ['if b:\n',
                         '    if a % 2 == 0:\n',
                         '    else:\n',
                         '        print("ODD!")\n']
        self.assert_context_accurate(FakeFile(NESTED_IF_ELSE_EXAMPLE), 8)

    def test_nested_if_else_with_all_else_branches(self):
        self.expected = ['if b:\n',
                         'else:\n',
                         '    if c % a == 0:\n',
                         '    else:\n',
                         '        if False:\n',
                         '        elif None:\n',
                         '        else:\n',
                         "            x = 27\n"]
        self.assert_context_accurate(FakeFile(NESTED_IF_ELSE_EXAMPLE), 20)

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

    def test_try_except_else_finally(self):
        self.expected = ['try:\n',
                         'except Exception as e:\n',
                         'else:\n',
                         'finally:\n',
                         '    sock.close()\n']
        self.assert_context_accurate(FakeFile(TRY_EXCEPT_ELSE_EXAMPLE), 15)

    def test_try_except_finally_without_else(self):
        test_file = FakeFile("""
try:
    a = 2
    b = 3
except KeyboardInterrupt:
    sys.exit(0)
except TypeError as e:
    print("error:", e)
    print("continuing")
finally:
    print("finally")
""")
        self.expected = ['try:\n',
                         'except KeyboardInterrupt:\n',
                         'except TypeError as e:\n',
                         'finally:\n',
                         '    print("finally")\n']
        self.assert_context_accurate(test_file, 11)

    def test_try_finally_without_except_or_else(self):
        test_file = FakeFile("""
try:
    a = 2
    b = 3
finally:
    a = None
    print("finally")
""")
        self.expected = ['try:\n',
                         'finally:\n',
                         '    print("finally")\n']
        self.assert_context_accurate(test_file, 7)


if __name__ == '__main__':
    unittest.main()
