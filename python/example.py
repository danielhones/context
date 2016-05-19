"""
A simple Python file for testing context.py
"""


class FooClass:
    def __init__(self, a, b):
        self.a = a
        self.b = b

    def bar(self):
        print("in bar")
        if self.a > self.b:
            print("a > b")
            print("bar: {}".format(self.a))

    def baz(self):
        print("in baz")
        if self.a == self.b:
            print("a == b")
            print("baz: {}".format(self.b))


def main():
    foo = FooClass(2, 3)
    foo.bar()
    foo.a == foo.b
    foo.baz()


if __name__ == "__main__":
    main()

if __name__ == "something different":
    baz()
    print("nope")
