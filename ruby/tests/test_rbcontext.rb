require 'minitest/autorun'
require 'parser/current'
load File.expand_path('../rbcontext', File.dirname(__FILE__))


class TestInsertTermColor < MiniTest::Test
  def test_color_whole_string
    teststr = "some example test string\nwith some new lines\tand tabs."
    colored = insert_term_color(teststr, 0, teststr.length, BLUE)
    expected = BLUE + teststr + END_COLOR
    assert_equal(colored, expected)
  end

  def test_color_substring
    teststr = "only color THIS word"
    colored = insert_term_color(teststr, 11, 15, BLUE)
    expected = "only color " + BLUE + "THIS" + END_COLOR + " word"
    assert_equal(colored, expected)
  end
end
