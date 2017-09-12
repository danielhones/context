# Context

See lines of code in context by folding/hiding irrelevant lines.  Sort of like grep that also shows you the branch that leads to a particular line or function call/object reference.  Imagine working inside a 5th or 6th nested if/else block with each block taking up more than a page of screen.  It would be nice to have a bird's-eye view of the if/elif conditions that lead to the line you're working on, without the clutter of the in-between lines of code and tracing through it manually.  That's what this tool does.

Right now there are tools for Python and Ruby, each using native parsers to build the abstract syntax tree.  This might be a questionable decision since it means maintaining the same program in two different languages... So that might eventually change if I find a good library/libraries for parsing each different target language and can write something to provide a shared interface for them across languages.  Until then, each file is only around 200 lines so it's not so bad yet.

If you use Emacs, take a look at the directions for [code-context.el](#emacs) included in this repo.


# Examples

Using `pycontext` itself as our example file, here's what we get if we `grep` for "SourceCode" (note this example will probably be outdated as things change in the file):

```
~/context$ grep -n SourceCode python/pycontext
93:class SourceCode(object):
177:            source = SourceCode(f, matcher, look_for, filename=f.name,
```

And here's using pycontext to do the same thing:

```
~/context$ pycontext -n SourceCode python/pycontext

 93: class SourceCode(object):
153: def main(look_for, files, search_type=None, ignore=None,
174:     for f in files:
176:         try:
177:             source = SourceCode(f, matcher, look_for, filename=f.name,
```

`pycontext` shows that the `SourceCode` code class is used in a try block in a for loop within the `main` function.  It's also very useful for showing the logic branch that leads to a particular line of code running.


# Usage

`pycontext` and `rbcontext` use the same command-line interface.  They read input from a file or list of files specified on the command-line or from stdin if there are no files listed.  There is no recursive option like there is with grep but you can get the same functionality by using `find` or something similar in a sub-shell: `pycontext -n skipped_files $(find . -name '*.py')`.  You can see the complete list of options from the command line with the `-h` flag:

```
$ pycontext -h
usage: pycontext [-h] [-c] [-n] [-l] [-e] [-v] look_for [paths [paths ...]]

Find lines in a Python source file and the context they're in.

positional arguments:
  look_for            Object to look for in the file. For certain searchtypes
                      this can be a comma-separated list or a Pythonlist such
                      as [6,8,10]
  paths               Files (or directories if -r option) to look in. If there
                      is no argument passed here, read from STDIN

optional arguments:
  -h, --help          show this help message and exit
  -c, --color         colorize output
  -n, --number-lines  include line numbers in results
  -l, --search-line   Search by line number. Allows a list as the look_for
                      argument
  -e, --search-regex  search by regexp
  -v, --verbose       Display information about errors and skipped files
```

```
$ rbcontext -h
Usage: rbcontext [options] look_for [paths]

Find lines in a Ruby source file and the context they're in

Options:
    -v, --verbose                    Display information about errors and skipped files
    -l, --search-line                Search by line number.  Can take a single integer or a Ruby list of integers as the look_for argument
    -e, --search-regex               Search by regular expression
    -c, --color                      Highlight line numbers and matches
    -n, --number-lines               Show line numbers in results
    -h, --help                       Print this help
```


# Installing

`pycontext` works out of the box.  Just put it somewhere in your path and make it executable and you're good to go.  For
`rbcontext`, do the same and install its dependency:

`gem install parser`


## Known Bugs

- Due to the hack of replacing `else` by `elsif true`, a Ruby file that uses try/except/else will not parse (because the `else` is subbed with `elsif True` which is a SyntaxError).  Hopefully I'll fix this soon; I recently fixed the same bug for pycontext.  The reason for this hack in the first place is that else lines have no line number in Ruby's AST parser.


## Updates

- The else bug that still exists in rbcontext has been fixed in pycontext


## rbcontext quirks

`rbcontext` does a few things a little strangely, but most of the bugs have been fixed and these remain as harmless oddities:

- The first non-comment line of each file/branch will be included in the results.  This is due to the AST including a `begin` node there that doesn't show up in the sourcefile but nonetheless starts a new branch and has a line number.
- Lines consisting only of `end` do not show up in the results.  This is not a huge problem because `end` only terminates a branch.  But it has the effect that if you search for a line number consisting only of `end`, there will be no results returned.


# Emacs

In the emacs directory of this repo, there's a file called `code-context.el`.  This defines some helpful functions that integrate `pycontext` and `rbcontext` with Emacs.  To use it, copy it into your `.emacs.d` directory and load it in `init.el`:

```
(load "~/.emacs.d/code-context.el")   ;; or whatever the correct path is for you
```

In order for this to work, `pycontext` and `rbcontext` need to be in a location included in your PATH environment variable within your Emacs environment.  This could be different than your usual PATH environment variable in a terminal, so if you see "command not found" when trying to use the context commands, double-check the file location against your Emacs env with `M-x shell-command env<RET>`.  Alternately, you can change the `context-script-map` variable to contain absolute paths to `pycontext` and `rbcontext`.

Here are the default keybindings for Emacs.  They use `C-c q` as the prefix, but that is configurable by changing the `context-keybinding-prefix` variable in `code-context.el`.

```
C-c q l         show context for line number at point
C-c q p         show context for symbol at point
C-c q c         prompt in the mini-buffer and show context for it
```
