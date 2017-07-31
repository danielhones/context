import os
import sys
import unittest

from helpers import *
from context import BLUE, END_COLOR, insert_term_color


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


if __name__ == "__main__":
    unittest.main()
