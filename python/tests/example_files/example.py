from otherfoo import baz


class FooClass:
    def __init__(self, var_a, var_b):
        self.some_variable = var_a
        self.other_variable = var_b

    def bar(self):
        print("in bar")
        if self.some_variable > self.other_variable:
            self.other_variable = self.some_variable
        elif self.some_variable < self.other_variable:
            self.some_variable = self.other_variable
        else:
            print("I give up")

    def baz(self):
        print("in baz")
        if self.some_variable == self.other_variable:
            print("some_variable == other_variable")
            print("baz: {}".format(self.other_variable))
        else:
            self.bar()


def main():
    foo = FooClass(2, 3)
    foo.bar()
    foo.some_variable = foo.other_variable
    foo.baz()


if __name__ == "__main__":
    main()
else:
    baz()
