def main(shift=2):
    print("In main():")
    i = 0
    while True:
        try:
            print("In try block")
            a = 10
            b = 3
            c = a ** b
            d = c << shift
            print("a,b,c,d:", a, b, c, d)
        except Exception as e:
            print("error in try block: ", e)
            print("continuing after error")
        else:
            print("successful try block")
        finally:
            i += 1


if __name__ == '__main__':
    print("Running main function")
    main()
