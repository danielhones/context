import unittest
import ast
import re
import sys, os

this_file_dir = os.path.dirname(os.path.realpath(__file__))
context_dir = os.path.abspath(os.path.join(this_file_dir, ".."))
sys.path.append(context_dir)

from context import *


EXAMPLE_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "example_files", "example.py")


class TestMakeMatcher(unittest.TestCase):
    def setUp(self):
        self.node = ast.Name(id="get_files", ctx=ast.Store(), lineno=23, col_offset=4)

    def test_search_regex(self):
        matcher = make_matcher(SEARCH_REGEX, "get_fi[a-z]{2,3}")
        self.assertTrue(matcher(self.node))

    def test_search_default(self):
        matcher = make_matcher(SEARCH_DEFAULT, 'get_files')
        self.assertTrue(matcher(self.node))

    def test_search_lineno(self):
        matcher = make_matcher(SEARCH_LINENO, str(self.node.lineno))
        self.assertTrue(matcher(self.node))


class TestWalk(unittest.TestCase):
    def setUp(self):
        self.source_code = SourceCode(EXAMPLE_FILE)
        self.tree = parse_source(EXAMPLE_FILE)

    def test_num_matches(self):
        def _matcher(node):
            try: return node.lineno
            except: return None

        # empty lines from here - http://stackoverflow.com/a/35026033/3199099
        empty_lines = sum(not line.strip() for line in self.source_code.lines)

        self.assertEqual(len(walk(self.tree, _matcher)), self.source_code.numlines - empty_lines)
        self.assertEqual(len(walk(self.tree, lambda _: None)), 0)

    def test_accurate_matches(self):
        matcher = make_matcher(SEARCH_DEFAULT, "bar")
        matches = walk(self.tree, matcher)
        # The list to match against is line numbers from example.py
        self.assertEqual(sorted(matches), [4, 9, 21, 23])


class TestFindTopLevel(unittest.TestCase):
    def test_class_and_func_definitions(self):
        tree = parse_source(EXAMPLE_FILE)
        self.assertEqual(find_top_level(tree), [4, 21])


if __name__ == "__main__":
    unittest.main()
