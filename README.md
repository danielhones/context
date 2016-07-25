# Context

See lines of code in context.  Sort of like grep that also shows you the branch that leads to a particular line or function call/object reference.  Imagine working inside a 5th or 6th nested if/else block with each block taking up more than a page of screen.  Not a great situation, but it happens sometimes.  It would be nice to have a bird's-eye view of what if/elif clauses lead to the line you're working on, without the clutter of the in-between lines of code and tracing through it manually.  That's what this does.

Right now there are tools for Python and Ruby, each using native parsers to build the abstract syntax tree.  This might be a questionable decision since it means maintaining the same program in two different languages... So that might eventually change if I find a good library/libraries for parsing each different target language and can write something to provide a shared interface for them across languages.  Until then, each file is only around 200 lines so it's not so bad yet.

If you use Emacs, take a look at [emacs-context](https://github.com/danielhones/emacs-context) which provides some functions for calling these command-line tools from within a buffer you're working on and displaying the output with syntax highlighting.


# Examples

Using `pycontext` itself as our example file, here's what we get if we `grep` for "walk" (note this example will probably be outdated as things change in the file):

```
~/context$ grep -n "walk" python/pycontext 
97:def walk(node, matcher, history=None):
122:        matches.extend(walk(i, matcher, history))
140:    matches = walk(tree, matcher)
158:        files = itertools.chain(*[os.walk(f) for f in files])
160:        # format returned by os.walk:
```

And here's using pycontext to do the same thing:

```
~/context$ pycontext walk python/pycontext 

 97:  def walk(node, matcher, history=None):
120:      for i in children:
122:          matches.extend(walk(i, matcher, history))
134:  def find_context(tree, look_for, matcher):
140:      matches = walk(tree, matcher)
150:  def main(look_for, files, search_type=SEARCH_DEFAULT, recursive=False,
157:      if recursive:
158:          files = itertools.chain(*[os.walk(f) for f in files])
```

`pycontext` shows the functions that `walk` is called from, the `if`, `elif` branches, `for` loops, etc.  In the grep example, you can't know for sure that line 122 calls `walk` recursively from within the `walk` function itself.  But you can see that with `pycontext`, as line 122 would otherwise show some other function definition before it.  `pycontext` also doesn't include the comment at line 160.


# Usage

`pycontext` and `rbcontext` use the same command-line interface.  You can see the options from the command line with the `-h` flag:

```
$ rbcontext -h
Find lines in a Ruby source file and the context they're in

Usage: rbcontext [options] look_for [paths]
Options:
    -v, --verbose                    display information about errors and skipped files
    -r, --recursive                  recursively search directory
    -n, --search-line                search by line number
    -e, --search-regex               search by regular expression
    -c, --color                      highlight line numbers and matches
    -h, --help                       print this help
    -i, --ignore IGNORE              comma-separated list of files and directories to ignore
```


# Installing

`pycontext` works out of the box.  Just put it somewhere in your path and make it executable and you're good to go.  For
`rbcontext`, do the same and install its dependency:

`gem install parser`


## rbcontext quirks

`rbcontext` does a few things a little strangely.  It's very likely this is due to copying the Python implementation for walking the AST and matching.  I'll work to get these fixed and behaving more sensibly, but until then here's how it works:

- The first non-comment line of each file/branch will be included in the results.  This is due to the AST including a `begin` node there that doesn't show up in the sourcefile but nonetheless starts a new branch and has a line number.
- Lines consisting only of `else` seem to not show up at all.  This is not ideal and I'll work on fixing it first because it can make the results look misleading.
- The same thing happens with lines consisting only of `end`.  This is not as big a problem as with `else` since `end` only terminates a branch.  But it has the additional effect that if you search for a line number consisting only of `end`, there will be no results returned.
