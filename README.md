MiniLisp
---

An interpreter for a minimal lisp dialect written in golang.
Check out the _examples_ directory to see things in action.

#### Functions

```lisp
# simple addition
(+ 4 5) # = 9
(+ (+ 2 2) 5) # = 9

# equality check
(= 5 6) # equates to () which is false
(= 5 5) # equates to 'true'

# print
(print) # just outputs a newline
(print 5) # prints 5
(print (+ 5 5)) # prints 10

# list operations
(list 1 2 3)            # -> (1 2 3)
(first (1 2 3))         # -> 1
(rest (1 2 3))          # -> (2 3)
(append (1 2 3) 4)      # -> (1 2 3 4)
(append (1 2 3) (4 5))  # -> (1 2 3 4 5)

# if statements
    condition          evaluated on false
     |                  |
(if (= 7 7) (print 2) (print 1)) # prints 2
              |
             evaluated on true

(if (= 7 8) (+) (print 1)) # an expression can be invalid if it is not executed
             |
            invalid

# define a symbol
(define a 5)
(print a) # prints 5

(define sum lambda (a b) (+ a b))
(sum 5 6) # -> 11
```

#### Parser Automaton

![](https://i.imgur.com/KAreuSe.jpg)
