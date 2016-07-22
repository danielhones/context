#!/usr/bin/env python3
import ast
import argparse
import sys
import os
import glob
import re
import itertools


IGNORE_DIRECTORIES = ["__pycache__", ".git"]
SEARCH_DEFAULT, SEARCH_LINENO, SEARCH_REGEX, SEARCH_DEFINITIONS = range(4)


def echo(*args):
    """A poor man's print"""
    output = " ".join([str(i) for i in args]) + "\n"
    sys.stdout.write(output)


class SourceCode(object):
    def __init__(self, filename, offset=1):
        """
        offset is the difference between the number of the first line of the file and 0.  It's probably 1.
        """
        self.offset = offset
        self.filename = filename
        self.lines = []
        with open(filename, "r") as f:
            self.lines = f.readlines()

    def __repr__(self):
        return "SourceCode('{}', {})".format(self.filename, self.offset)

    def line(self, lineno):
        return self.lines[lineno - self.offset]

    @property
    def numlines(self):
        return len(self.lines)

    def format_line(self, lineno):
        return "{}:  {}".format(str(lineno).rjust(len(str(self.numlines))), self.line(lineno))


def make_matcher(search_type, look_for):
    def _regex(node):
        try: lineno = node.lineno
        except: return None
        for attr in node._fields:
            try:
                if look_for.match(getattr(node, attr)):
                    return lineno
            except:
                pass

    def _default(node):
        try: lineno = node.lineno
        except: return None
        for attr in node._fields:
            if getattr(node, attr) == look_for:
                return lineno       

    def _lineno(node):
        try: lineno = node.lineno
        except: return None
        return lineno if lineno == look_for else None

    # Convert look_for to the type that the matcher functions need.
    # Doing it here rather than inside the matcher function because it only needs to happen once
    look_for = {SEARCH_DEFAULT: lambda x: x,
                SEARCH_LINENO: int,
                SEARCH_REGEX: re.compile}[search_type](look_for)
    return {SEARCH_DEFAULT: _default,
            SEARCH_LINENO: _lineno,
            SEARCH_REGEX: _regex}[search_type]


def walk(node, matcher, history=None):
    history = [] if history is None else history[:]
    matches = []

    def append_to_history(item):
        try: history.append(item.lineno)
        except AttributeError: pass

    def pop_from_history():
        try: history.pop()
        except IndexError: pass

    children = list(ast.iter_child_nodes(node))
    match = matcher(node)

    if match:
        matches.append(match)
        matches = list(set(matches + history))

    if len(children) == 0:
        return set(matches)

    append_to_history(node)
    for i in children:
        append_to_history(i)
        matches.extend(walk(i, matcher, history))
        pop_from_history()
    pop_from_history()

    return set(matches)


def find_top_level(tree, depth=1):
    matches = [i.lineno for i in ast.iter_child_nodes(tree) if isinstance(i, (ast.FunctionDef, ast.ClassDef))]
    return matches


def find_context(tree, look_for, matcher):
    """
    Return a list of lines indicating branches that lead to the object we're looking for
    """
    # TODO: Add something to parse look_for, for example FooClass.bar should only match bar
    #       function calls that are attributes of a FooClass instance
    matches = walk(tree, matcher)
    return sorted(list(matches))


def parse_source(filename):
    with open(filename, "r") as f:
        tree = ast.parse(f.read(), filename=filename)
    return tree


def main(look_for, files, search_type=SEARCH_DEFAULT, recursive=False, ignore=IGNORE_DIRECTORIES, verbose=False):
    """
    look_for is a string, files is a list of paths
    """
    if recursive:
        files = itertools.chain(*[os.walk(f) for f in files])
    else:
        # format returned by os.walk:
        files = [(os.curdir, [], [f]) for f in files]

    all_contexts = {}
    skipped_files = {}

    for (directory, _, filenames) in files:
        abspath = os.path.abspath(directory)

        current_ignore = []
        for i in ignore:
            current_ignore.extend(glob.glob(os.path.abspath(i)))

        if abspath in current_ignore:
            continue
        for fn in filenames:
            source_file = os.path.join(directory, fn)
            if os.path.abspath(source_file) in current_ignore or source_file in current_ignore:
                continue
            try:
                source = SourceCode(source_file)
                tree = parse_source(source_file)
                if search_type == SEARCH_DEFINITIONS:
                    context = find_top_level(tree)
                else:
                    context = find_context(tree, look_for, make_matcher(search_type, look_for))
            except KeyboardInterrupt as e:
                sys.exit(1)
            except Exception as e:
                skipped_files[source_file] = "{}: {}".format(e.__class__.__name__, e)
                continue
            if len(context) == 0:
                continue

            if recursive or len(files) > 1:
                echo("\n{}".format(source_file))

            all_contexts[source_file] = context
            echo("\n" + "".join([source.format_line(i) for i in context]))

    if verbose and skipped_files:
        echo("Skipped these files due to errors:\n{}".format(
            "\n".join(["{}: {}".format(key, skipped_files[key]) for key in skipped_files])))

    # to make testing easier, don't have to capture and parse stdout:
    return all_contexts


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Find lines in a Python source file and the context they're in."
    )
    parser.add_argument("look_for", help="object to look for in the file")
    parser.add_argument("paths", help="files (or directories if -r option) to look in", nargs="+")
    parser.add_argument("-r", "--recursive", action="store_true", help="recursively search directory")
    parser.add_argument("-n", "--search-line", dest="search_type", action="store_const", const=SEARCH_LINENO,
                        help="search by line number")
    parser.add_argument("-e", "--search-regex", dest="search_type", action="store_const", const=SEARCH_REGEX,
                        help="search by regexp")
    parser.add_argument("-d", "--search-defs", dest="search_type", action="store_const", const=SEARCH_DEFINITIONS,
                        help=("just look for class and function definitions.  The look_for argument "
                              "should be an integer indicating the maxmimum depth of search"))
    parser.add_argument("-v", "--verbose", action="store_true",
                        help="display information about errors and skipped files")
    parser.add_argument("-i", "--ignore",
                        help=("comma-separated list of files and directories "
                              "to ignore, default is {}".format(IGNORE_DIRECTORIES)))
    args = parser.parse_args()

    if args.ignore:
        ignore = args.ignore.split(",") + IGNORE_DIRECTORIES
    else:
        ignore = IGNORE_DIRECTORIES

    main(args.look_for,
         args.paths,
         search_type=args.search_type or SEARCH_DEFAULT,
         recursive=args.recursive,
         ignore=ignore,
         verbose=args.verbose)
