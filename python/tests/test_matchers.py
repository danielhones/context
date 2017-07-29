import ast
import unittest

from helpers import *
from context import LineNoMatcher, RegexMatcher, LIGHT_GREEN, END_COLOR


class TestLineNoMatcherFormatLookFor(unittest.TestCase):
    pass


class TestLineNoMatcherNodeMatches(unittest.TestCase):
    def setUp(self):
        self.test_node = ast.Name(id="get_files", ctx=ast.Store(), lineno=23, col_offset=4)

    def test_matching_line(self):
        matcher = LineNoMatcher(23, [])
        self.assertTrue(matcher.node_matches(self.test_node))

    def test_non_matching_line(self):
        matcher = LineNoMatcher(7, [])
        self.assertFalse(matcher.node_matches(self.test_node))


class TestLineNoMatcherColorLine(unittest.TestCase):
    def setUp(self):
        self.test_line = "if 'key' in settings and settings['key'] == 'value':"
        self.matcher = LineNoMatcher(10, [])

    def test_matching_line_with_color(self):
        colored = self.matcher.color_line(self.matcher.look_for[0], self.test_line, LIGHT_GREEN)
        expected = LIGHT_GREEN + self.test_line + END_COLOR
        self.assertEqual(colored, expected)

    def test_non_matching_line_with_color(self):
        colored = self.matcher.color_line(99, self.test_line, LIGHT_GREEN)
        expected = self.test_line
        self.assertEqual(colored, expected)

    def test_matching_line_without_color(self):
        colored = self.matcher.color_line(self.matcher.look_for[0], self.test_line, None)
        expected = self.test_line
        self.assertEqual(colored, expected)

    def test_non_matching_line_without_color(self):
        colored = self.matcher.color_line(99, self.test_line, None)
        expected = self.test_line
        self.assertEqual(colored, expected)


if __name__ == '__main__':
    unittest.main()
