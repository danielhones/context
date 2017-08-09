require 'minitest/autorun'
require 'parser/current'
load File.expand_path('../rbcontext', File.dirname(__FILE__))


class NodeVisitorTestSubclass < NodeVisitor
  attr_accessor :specific_visits, :generic_visits

  def initialize
    @specific_visits = []
    @generic_visits = []
  end

  def visit_send(node)
    @specific_visits.push("test visited send")
  end

  def generic_visit(node)
    @generic_visits.push(node.type)
  end
end


class TestNodeVisitor < MiniTest::Test
  def setup
    @visitor = NodeVisitorTestSubclass.new
  end

  def test_specific_visit_method_right_type
    node = AST::Node.new('send')
    @visitor.visit(node)
    assert_equal(@visitor.specific_visits.last, 'test visited send')
  end

  def test_specific_visit_method_wrong_type
    node = AST::Node.new('int')
    @visitor.visit(node)
    assert_equal(@visitor.specific_visits.length, 0)
  end

  def test_generic_visit
    node = AST::Node.new('int')
    @visitor.visit(node)
    assert_equal(@visitor.generic_visits.last, node.type)
  end
end
