require 'minitest/autorun'
load File.expand_path('../rbcontext', File.dirname(__FILE__))


class TestLineNoMatcherFormatLookFor < MiniTest::Test
  def test_look_for_single_line_number
    matcher = LineNoMatcher.new(10, [])
    assert_equal(matcher.look_for, [10])
  end

  def test_look_for_multiple_line_numbers
    line_nos = [10, 19, 31, 42]
    matcher = LineNoMatcher.new(line_nos, [])
    assert_equal(matcher.look_for, line_nos)
  end
end


class TestLineNoMatcherNodeMatches < MiniTest::Test
  def setup
    @expr = 
    @location = Parser::Source::Map.new(@expr)
    @test_node = Parser::AST::Node.new('send', @location)
  end

  def test_matching_line
    matcher = LineNoMatcher.new(23, [])
    assert(matcher.node_matches(@test_node))
  end

  def test_non_matching_line
    matcher = LineNoMatcher.new(7, [])
    assert_nil(matcher.node_matches(@test_node))
  end
end
