#!/usr/bin/python3


import ast
import sys


def echo(*args):
    """Avoiding print so it works with Python 2 and 3"""
    output = " ".join([str(i) for i in args]) + "\n"
    sys.stdout.write(output)


class SourceCode(object):
    def __init__(self, filename, offset=1):
        """
        offset is the difference between 0 and the number of the first line of the file.  It's probably 1
        """
        self.offset = offset
        self.filename = filename
        self.lines = []
        with open(filename, "r") as f:
            self.lines = f.readlines()

    def __repr__(self):
        return "SourceCode('{}', {})".format(self.filename, self.offset)

    def line(self, lineno):
        return self.lines[lineno + self.line_offset]

    def numlines(self):
        return len(self.lines)

    def format_line(self, lineno):
        return "{}:  {}".format(str(lineno).ljust(self.numlines), self.line(lineno))


def usage():
    echo(
        "Usage:\n"
        "context.py <filename> <object>"
    )


def parse_args():
    if len(sys.argv) < 3:
        usage()
        sys.exit(1)

    return sys.argv[1], sys.argv[2]


def main(source_file, look_for):
    source = SourceCode(source_file)
    echo(source)


if __name__ == "__main__":
    main(*parse_args())
