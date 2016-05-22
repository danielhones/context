import unittest
import ast
import sys
if sys.version_info.major == 2:
    from StringIO import StringIO
else:
    from io import StringIO
import os
import shutil

THIS_FILE_DIR = os.path.dirname(os.path.realpath(__file__))
context_dir = os.path.abspath(os.path.join(THIS_FILE_DIR, ".."))
sys.path.append(context_dir)

import context
from context import *


EXAMPLE_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "example_files", "example.py")
BAR_LINE_NOS = [4, 9, 21, 23]


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
        self.assertEqual(sorted(matches), BAR_LINE_NOS)


class TestFindTopLevel(unittest.TestCase):
    def test_class_and_func_definitions(self):
        tree = parse_source(EXAMPLE_FILE)
        self.assertEqual(find_top_level(tree), [4, 21])


class TestMain(unittest.TestCase):
    THIS_FILE_RELPATH = "./test.py"
    EXAMPLE_FILE_RELPATH = "./example_files/example.py"

    def setUp(self):
        # from this StackOverflow answer - http://stackoverflow.com/a/1218951/3199099
        self.orig_stdout = sys.stdout
        sys.stdout = self.stdout = StringIO()
        self.original_dir = os.getcwd()
        os.chdir(THIS_FILE_DIR)

    def tearDown(self):
        sys.stdout = self.orig_stdout
        os.chdir(self.original_dir)

    def test_recursive_current_dir_multiple_match(self):
        contexts = context.main("bar", ["."], recursive=True)
        self.assertEqual(len(contexts.keys()), 2)
        self.assertTrue(self.THIS_FILE_RELPATH in contexts)
        self.assertTrue(self.EXAMPLE_FILE_RELPATH in contexts)
        self.assertEqual(contexts[self.EXAMPLE_FILE_RELPATH], BAR_LINE_NOS)

    def test_recursive_current_dir_single_match(self):
        self.contexts = context.main("TestCase", ["."], recursive=True)
        self.assertAccurate(1, self.THIS_FILE_RELPATH, self.EXAMPLE_FILE_RELPATH)

    def test_recursive_multiple_directories(self):
        if not os.path.isdir("/tmp/context_test"):
            os.mkdir("/tmp/context_test")
        TMP_FILE = "/tmp/context_test/temp.py"
        shutil.copyfile(EXAMPLE_FILE, TMP_FILE)
        self.contexts = context.main("bar", [".", "/tmp"], recursive=True)
        self.assertAccurate(3,
                            (self.EXAMPLE_FILE_RELPATH, self.THIS_FILE_RELPATH, TMP_FILE),
                            "doesn't matter",
                            (self.EXAMPLE_FILE_RELPATH, BAR_LINE_NOS),
                            (TMP_FILE, BAR_LINE_NOS))
        os.remove(TMP_FILE)

    def test_ignore_files_with_directory(self):
        self.contexts = context.main("bar", ["."], ignore=["example_files/"], recursive=True)
        self.assertAccurate(1, self.THIS_FILE_RELPATH, self.EXAMPLE_FILE_RELPATH)

    def test_ignore_files_with_directory_and_filename(self):
        for path in ["example_files/example.py", "./example_files/example.py"]:
            self.contexts = context.main("bar", ["."], ignore=[path], recursive=True)
            self.assertAccurate(1, self.THIS_FILE_RELPATH, self.EXAMPLE_FILE_RELPATH)

    def test_ignore_files_with_filename(self):
        self.contexts = context.main("bar", ["."], ignore=["test.py"], recursive=True)
        self.assertAccurate(1, self.EXAMPLE_FILE_RELPATH, self.THIS_FILE_RELPATH, (self.EXAMPLE_FILE_RELPATH, BAR_LINE_NOS))

    def test_ignore_files_with_glob(self):
        self.contexts = context.main("bar", ["."], ignore=["*est*"], recursive=True)
        self.assertAccurate(1, self.EXAMPLE_FILE_RELPATH, self.THIS_FILE_RELPATH, (self.EXAMPLE_FILE_RELPATH, BAR_LINE_NOS))

    def test_search_default(self):
        self.contexts = context.main("bar", [EXAMPLE_FILE])
        self.assertAccurate(1, EXAMPLE_FILE, "doesn't matter", (EXAMPLE_FILE, BAR_LINE_NOS))

    def test_search_regex(self):
        self.contexts = context.main("bar", [EXAMPLE_FILE])
        self.assertAccurate(1, EXAMPLE_FILE, "doesn't matter", (EXAMPLE_FILE, BAR_LINE_NOS))

    def test_search_lineno(self):
        self.contexts = context.main("17", [EXAMPLE_FILE], SEARCH_LINENO)
        self.assertAccurate(1, EXAMPLE_FILE, "doesn't matter", (EXAMPLE_FILE, [4, 14, 16, 17]))

    def test_definitions(self):
        self.contexts = context.main("doesn't matter", [EXAMPLE_FILE], SEARCH_DEFINITIONS)
        self.assertAccurate(1, EXAMPLE_FILE, "doesn't matter", (EXAMPLE_FILE, [4, 21]))

    def assertAccurate(self, num_matches, should_have, should_not_have, *key_values):
        """
        Takes optional number of tuple arguments in the form (key, value) to match in self.contexts
        """
        should_have = should_have if type(should_have) is tuple else (should_have,)
        should_not_have = should_not_have if type(should_not_have) is tuple else (should_not_have,)
        self.assertEqual(len(self.contexts), num_matches)
        [self.assertTrue(i in self.contexts) for i in should_have]
        [self.assertFalse(i in self.contexts) for i in should_not_have]
        if len(key_values) > 0:
            [self.assertEqual(self.contexts[k], v) for k, v in key_values]


if __name__ == "__main__":
    unittest.main()
