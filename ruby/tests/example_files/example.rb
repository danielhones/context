class FooClass
  attr_accessor :a, :b

  def initialize(a, b)
    @a = a
    @b = b
  end

  def bar
    puts "in bar"
    if @a > @b
      nil
    end
  end

  def baz
    puts "in baz"
    if @a == @b
      puts "@a == @b"
      puts "baz: #{@b}"
    end
  end
end


def main
  foo = FooClass.new(2, 3)
  foo.bar
  foo.a = foo.b
  foo.baz
end


if __FILE__ == $0
  main
else
  baz
end
  
