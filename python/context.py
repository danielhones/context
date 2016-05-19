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
        return self.lines[lineno - self.offset]

    @property
    def numlines(self):
        return len(self.lines)

    def format_line(self, lineno):
        return "{}:  {}".format(str(lineno).rjust(len(str(self.numlines))), self.line(lineno))


def walk(node, matcher=lambda x: None, history=None):
    matches = []
    history = [] if history is None else history[:]

    def add_to_history(item):
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
        return matches
    
    add_to_history(node)
    
    for i in children:
        add_to_history(i)
        matches.extend(walk(i, matcher, history))
        pop_from_history()

    pop_from_history()
    
    return matches


def find_context(source, tree, look_for):
    """
    Return a list of lines indicating branches that lead to the object we're looking for
    """
    matching_lines = []

    def matcher(node):
        for attr in node._fields:
            value = getattr(node, attr)
            if value == look_for:
                try:
                    return node.lineno
                except AttributeError:
                    return None
        return None
        
    matches = walk(tree, matcher)
    matches = list(set(matches))
    return sorted(matches)


def parse_source(filename):
    with open(filename, "r") as f:
        tree = ast.parse(f.read(), filename=filename)
    return tree

    
def usage():
    echo("Usage:\ncontext.py <filename> <object>")


def parse_args():
    if len(sys.argv) < 3:
        usage()
        sys.exit(1)

    return sys.argv[1], sys.argv[2]


def main(source_file, look_for):
    source = SourceCode(source_file)
    tree = parse_source(source_file)
    context = find_context(source, tree, look_for)
    echo("\n" + "".join([source.format_line(i) for i in context]))


if __name__ == "__main__":
    main(*parse_args())
