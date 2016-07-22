#!/usr/bin/env ruby
require 'ripper'
require 'stringio'
require 'pp'
require 'optparse'

oldstderr, $stderr = $stderr, StringIO.new  # suppress warnings
require 'parser/current'
$stderr = oldstderr


SEARCH_DEFAULT, SEARCH_LINENO, SEARCH_REGEX, SEARCH_DEFINITIONS = (0..3).to_a
DEFAULT_NUM_COLOR = "blue"
DEFAULT_LINE_COLOR = "red"
IGNORE_DIRECTORIES = [".git"]


def make_formatter(format)
  # format should be a terminal color/format like "[1m[31m" for bold red
  endformat = "[0m"
  lambda { |string| "\e#{format}" + string + "\e#{endformat}" }
end


class SourceCode
  def initialize(filename, :look_for => nil, offset => 1, num_format => nil, line_format => nil)
    @filename = filename
    @offset = offset
    @lines = File.readlines(filename)
    @look_for = look_for
    @number_formatter = num_format.nil ? make_formatter(num_format) : lambda { |x| x }
    @line_formatter = (look_for && line_format) ? make_line_formatter(line_format) : lambda { |x| x }
  end

  def make_line_formatter(format)
    formatted_lookfor = make_formatter(format).call(look_for)
    lambda do |string|
      string.gsub(Regexp.new(look_for), formatted_lookfor)
    end
  end
  
  def line(lineno)
    @lines[lineno - @offset]
  end

  def numlines
    @lines.length
  end

  def format_line(lineno)
    formatted_lineno = @number_formatter.call( lineno.to_s.rjust(numlines.to_s.length) )
    formatted_line = @line_formatter.call( line(lineno) )
    "#{formatted_lineno}:  #{formatted_line}"
  end
end


def make_matcher(search_type, look_for)
  if search_type == SEARCH_REGEX
    look_for = Regexp.new(look_for)
  elsif search_type == SEARCH_LINENO
    look_for = look_for.to_i
  end

  _default = lambda do |node, lineno|
    if node.respond_to?(:children) && node.children.map(&:to_s).include?(look_for)
      return lineno
    end
  end

  _regex = lambda do |node, lineno|
    node._fields.each do |attr|
      begin
        return lineno if look_for.match(getattr(node, attr))
      rescue
      end
    end
  end

  _lineno = lambda do |node, lineno|
    return lineno == look_for ? lineno : nil
  end

  {SEARCH_DEFAULT => _default,
   SEARCH_LINENO => _lineno,
   SEARCH_REGEX => _regex}[search_type]
end


def walk(node, matcher, history=nil)
  history = history.nil? ? [] : history.dup.uniq
  matches = []

  lineno = node.location.line rescue nil
  match = matcher.call(node, lineno)
  if match
    matches.push(match)
    matches = (matches + history).uniq
  end
 
  if !node.is_a?(AST::Node) || node.children.nil? || node.children.length == 0
    return matches.uniq
  end

  history.push(lineno)
  node.children.each do |child|
    history.push(child.location.line) rescue nil
    matches += walk(child, matcher, history)
    history.pop
  end
  history.pop

  return matches.compact.sort.uniq
end

  
def find_top_level(tree, depth=1)
end


def main(look_for, files, search_type=SEARCH_DEFAULT, recursive=false, ignore=IGNORE_DIRECTORIES, verbose=false)
  # TODO: get this working for recursive option

  # check for recursive, walk files/directories
  all_contexts = {}
  skipped_files = {}
  
  files.each do |f|
    source_file = f  # grab absolute path here
    begin
      source = SourceCode.new(source_file)
      tree = Parser::CurrentRuby.parse_file(source_file)

      if search_type == SEARCH_DEFINITIONS
        context = find_top_level(tree)
      else
        context = walk(tree, make_matcher(search_type, look_for))
      end
    rescue SystemExit, Interrupt
      raise
    end
    next if context.length == 0

    puts source_file if recursive || files.length > 1

    context.each { |lineno| puts source.format_line(lineno) }
  end

  if verbose && skipped_files.length > 0
    puts "Skipped these files due to errors:"
    skipped_files.each { |k, v| puts "#{k}:  #{v}" }
  end
  return all_contexts
end


if __FILE__ == $0
  options = {}
  OptionParser.new do |opts|
    opts.set_banner("Find lines in a Ruby source file and the context they're in\n\n" +
                    "Usage: #{$0} [options] look_for paths [paths...]\n\n" +
                    "Options:\n")
    opts.on("-v", "--verbose", "display information about errors and skipped files") { |v| options[:verbose] = true }
    opts.on("-r", "--recursive", "recursively search directory") { |v| options[:recursive] = true }
    opts.on("-n", "--search-line", "search by line number") { |v| options[:search_type] = SEARCH_LINENO }
    opts.on("-e", "--search-regex", "search by regular expression") { |v| options[:search_type] = SEARCH_REGEX }
    opts.on("-d", "--search-defs", "look for class, module, and function definitions") { |v|
      options[:search_type] = SEARCH_DEFINITIONS
    }
    opts.on("-h", "--help", "print this help") { puts opts; exit }
    opts.on("-i IGNORE", "--ignore IGNORE", "comma-separated list of files and directories to ignore") { |v|
      options[:ignore] = IGNORE_DIRECTORIES + v.split(",")
    }
    if ARGV[0].nil? || ARGV[1].nil?
      puts opts; exit
    end
  end.parse!

  look_for = ARGV[0]
  paths = ARGV.slice(1,ARGV.length)

  main(look_for,
       paths,
       options[:search_type] || SEARCH_DEFAULT,
       options[:recursive],
       options[:ignore],
       options[:verbose])
end
