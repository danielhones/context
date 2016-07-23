#!/usr/bin/env ruby
require 'ripper'
require 'stringio'
require 'pp'
require 'optparse'
require 'find'

oldstderr, $stderr = $stderr, StringIO.new  # suppress warnings
require 'parser/current'
$stderr = oldstderr


SEARCH_DEFAULT, SEARCH_LINENO, SEARCH_REGEX, SEARCH_DEFINITIONS = (0..3).to_a
IGNORE_DIRECTORIES = [".git"]
BLUE = "\e[34m"
LIGHT_GREEN =  "\e[92m"
END_COLOR = "\e[0m"


class SourceCode
  def initialize(filename, look_for: nil, num_color: nil, line_color: nil, offset: 1)
    @filename = filename
    @offset = offset
    @lines = File.readlines(filename)
    @look_for = look_for
    @num_color = num_color || ""
    @line_color = line_color || ""
    @end_color = (num_color || line_color) ? END_COLOR : ""
  end
  
  def add_line_color(line)
    #formatted_lookfor = make_formatter(format).call(@look_for)
    #string.gsub(Regexp.new(@look_for), formatted_lookfor)
    colored_look_for = @line_color + @look_for + @end_color
    line.gsub(Regexp.new(@look_for), colored_look_for)
  end

  def add_num_color(lineno)
    @num_color + lineno + @end_color
  end
  
  def line(lineno)
    @lines[lineno - @offset]
  end

  def numlines
    @lines.length
  end

  def format_line(lineno)
    colored_lineno = add_num_color( lineno.to_s.rjust(numlines.to_s.length) )
    colored_line = add_line_color( line(lineno) )
    "#{colored_lineno}:  #{colored_line}"
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


def main(look_for, files, search_type: SEARCH_DEFAULT, recursive: false,
         ignore: IGNORE_DIRECTORIES, verbose: false, color: false)
  num_color = color ? BLUE : nil
  line_color = color ? LIGHT_GREEN : nil

  # check for recursive, walk files/directories
  all_contexts = {}
  skipped_files = {}

  if recursive
    files = files.map{ |dir|
      Find.find(dir).map { |f|
        File.expand_path(f, Dir.getwd)
      }
    }.flatten
  end

  files.each do |f|
    source_file = f  # grab absolute path here
    context = []
    begin
      source = SourceCode.new(source_file, look_for: look_for, num_color: num_color, line_color: line_color)
      tree = Parser::CurrentRuby.parse_file(source_file)

      if search_type == SEARCH_DEFINITIONS
        context = find_top_level(tree)
      else
        context = walk(tree, make_matcher(search_type, look_for))
      end
    rescue SystemExit, Interrupt
      raise
    rescue => e
      skipped_files[f] = e
    end
    next if context.length == 0

    puts source_file if recursive || files.length > 1

    context.each { |lineno| puts source.format_line(lineno) }
    puts ""
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
    opts.on("-c", "--color", "highlight line numbers and matches") { |v| options[:color] = true }
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
  paths = ARGV.slice(1, ARGV.length)

  main(look_for, paths,
       search_type: options[:search_type] || SEARCH_DEFAULT,
       recursive: options[:recursive],
       ignore: options[:ignore],
       verbose: options[:verbose],
       color: options[:color])
end
