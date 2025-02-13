// Copyright © 2018 The ELPS authors

package lisp_test

import (
	"testing"

	"github.com/luthersystems/elps/elpstest"
)

func TestEval(t *testing.T) {
	tests := elpstest.TestSuite{
		{"raw strings", elpstest.TestSequence{
			{`"""a raw string"""`, `"a raw string"`, ""},
			{`"""a raw
string"""`, `"a raw\nstring"`, ""},
			{`"""""a raw
string"""`, `"\"\"a raw\nstring"`, ""},
		}},
		{"string escape sequences", elpstest.TestSequence{
			{`"a string"`, `"a string"`, ""},
			{`"a\nstring"`, `"a\nstring"`, ""},
			{`"\"\"a\nstring"`, `"\"\"a\nstring"`, ""},
		}},
		{"quotes", elpstest.TestSequence{
			{"3", "3", ""},
			// a single quote on a self-evaluating expression does not show up.
			{"'3", "3", ""},
			// a double quotes on a self-evaluating expression both show up.
			{"''3", "''3", ""},
		}},
		{"symbols", elpstest.TestSequence{
			{"()", "()", ""},
			{"'true", "'true", ""},
			{"true", "true", ""},
			{"'false", "'false", ""},
			{"false", "false", ""},
			// A bit brittle, but it's ok for now. Replace with a more robust
			// test later if problematic.
			{"a", "test:1: unbound symbol: a", ""},
		}},
		{"set", elpstest.TestSequence{
			{`(set 'x 1)`, `1`, ``},
			{`x`, `1`, ``},
			{`(set 'x 2)`, `2`, ``},
			{`x`, `2`, ``},
			{`(set 'true 2)`, `test:1: lisp:set: cannot rebind constant: true`, ``},
		}},
		{"lists basics", elpstest.TestSequence{
			{"'()", "'()", ""},
			{"'(1 2 3)", "'(1 2 3)", ""},
			{"(nth '() 0)", "()", ""},
			{"(nth '() 1)", "()", ""},
			{"(nth '() 2)", "()", ""},
			{"(nth '(1) 0)", "1", ""},
			{"(nth '(1) 1)", "()", ""},
			{"(nth '(1) 2)", "()", ""},
			{`(rest '(1 2))`, `'(2)`, ``},
			{`(rest '(1))`, `()`, ``},
			{`(car ())`, `()`, ``},
			{`(cdr '(1 2))`, `'(2)`, ``},
			{`(cdr '(1))`, `()`, ``},
			{`(car '(a b c))`, `a`, ``}, // not quoted
		}},
		{"types", elpstest.TestSequence{
			{`(symbol? true)`, `true`, ""},
			{`(symbol? false)`, `true`, ""},
			{`(symbol? 1)`, `false`, ""},
			{`(symbol? ())`, `false`, ""},
			{`(bytes? (to-bytes "abc"))`, `true`, ""},
			{`(bytes? "abc")`, `false`, ""},
			{`(string? "abc")`, `true`, ""},
			{`(string? (to-bytes "abc"))`, `false`, ""},
			{`(string? 'abc)`, `false`, ""},
			{`(int? 0)`, `true`, ""},
			{`(int? 0.0)`, `false`, ""},
			{`(int? -12)`, `true`, ""},
			{`(int? "1")`, `false`, ""},
			{`(int? (to-int 0))`, `true`, ""},
			{`(int? (to-int 0.0))`, `true`, ""},
			{`(int? (to-int -12))`, `true`, ""},
			{`(int? (to-int "1"))`, `true`, ""},
			{`(to-int "123")`, `123`, ""},
			{`(float? 0)`, `false`, ""},
			{`(float? 0.0)`, `true`, ""},
			{`(float? -12)`, `false`, ""},
			{`(float? "1")`, `false`, ""},
			{`(float? (to-float 0))`, `true`, ""},
			{`(float? (to-float 0.0))`, `true`, ""},
			{`(float? (to-float -12))`, `true`, ""},
			{`(float? (to-float "1"))`, `true`, ""},
			{`(to-float "-12.75")`, `-12.75`, ""},
			{`(bool? true)`, `true`, ""},
			{`(bool? false)`, `true`, ""},
			{`(bool? lisp:true)`, `true`, ""},
			{`(bool? lisp:false)`, `true`, ""},
			{`(bool? ())`, `false`, ""},
			{`(bool? 't)`, `false`, ""},
			{`(list? ())`, `true`, ""},
			{`(list? '())`, `true`, ""},
			{`(list? '(1 2 3))`, `true`, ""},
			{`(list? [1 2 3])`, `true`, ""},
			{`(list? (vector 1 2 3))`, `false`, ""},
			{`(list? (sorted-map))`, `false`, ""},
			{`(vector? ())`, `false`, ""},
			{`(vector? '())`, `false`, ""},
			{`(vector? '(1 2 3))`, `false`, ""},
			{`(vector? [1 2 3])`, `false`, ""},
			{`(vector? (vector 1 2 3))`, `true`, ""},
			{`(vector? (sorted-map))`, `false`, ""},
			{`(array? ())`, `false`, ""},
			{`(array? '())`, `false`, ""},
			{`(array? '(1 2 3))`, `false`, ""},
			{`(array? [1 2 3])`, `false`, ""},
			{`(array? (vector 1 2 3))`, `true`, ""},
			{`(array? (sorted-map))`, `false`, ""},
			{`(sorted-map? ())`, `false`, ""},
			{`(sorted-map? '())`, `false`, ""},
			{`(sorted-map? '(1 2 3))`, `false`, ""},
			{`(sorted-map? [1 2 3])`, `false`, ""},
			{`(sorted-map? (vector 1 2 3))`, `false`, ""},
			{`(sorted-map? (sorted-map))`, `true`, ""},
		}},
		{"function basics", elpstest.TestSequence{
			{"(lambda ())", "(lambda ())", ""},
			{"((lambda ()))", "()", ""},
			{"(lambda (x) x)", "(lambda (x) x)", ""},
			{"((lambda (x) x) 1)", "1", ""},
			{"(lambda (x) (+ x 1))", "(lambda (x) (+ x 1))", ""},
			{"((lambda () (+ 1 1)))", "2", ""},
			{"((lambda (n) (+ n 1)) 1)", "2", ""},
			{"((lambda (x y) (+ x y)) 1 2)", "3", ""},
			{"((lambda (x &rest y) (cons x y)) 1 2 3)", "'(1 2 3)", ""},
			{"((lambda (&rest x) (reverse 'list x)) 1 2 3)", "'(3 2 1)", ""},
		}},
		{"concat", elpstest.TestSequence{
			{`(concat 'vector () ())`, `(vector)`, ""},
			{`(concat 'string "a" (to-bytes "b") '(99))`, `"abc"`, ""},
			{`(to-string (concat 'bytes "a" (to-bytes "b") (vector 99)))`, `"abc"`, ""},
		}},
		{"length", elpstest.TestSequence{
			{`(length "abc")`, `3`, ``},
			{`(length '(a b))`, `2`, ``},
			{`(length [])`, `0`, ``},
			{`(length (vector 1 2 3 4))`, `4`, ``},
		}},
		{"emtpy?", elpstest.TestSequence{
			{`(empty? "abc")`, `false`, ``},
			{`(empty? '(a b))`, `false`, ``},
			{`(empty? [])`, `true`, ``},
			{`(empty? (vector 1 2 3 4))`, `false`, ``},
			{`(empty? (vector))`, `true`, ``},
			{`(empty? "")`, `true`, ``},
		}},
		{"funcall", elpstest.TestSequence{
			{`(funcall (lambda () 1))`, `1`, ""},
			{`(funcall (lambda (x) (+ 1 x)) 1)`, `2`, ""},
			{`(funcall (lambda (&rest xs) (cdr xs)) 1 2 3)`, `'(2 3)`, ""},
			{`(funcall '+ 1 2 3)`, `6`, ""},
			{`(funcall 'lisp:+ 1 2 3)`, `6`, ""},
			{`(defun f () "outer")`, `()`, ""},
			{`(defun inner ()
				(labels ((f () "innermost"))
					#^(apply 'f %&rest)))`, `()`, ""},
			{`(defun middle ()
				(let ([one (inner)])
					(labels ((f () "middle"))
						#^(apply one %&rest))))`, `()`, ""},
			{`(funcall (inner))`, `"outer"`, ""},
			{`(funcall (middle))`, `"outer"`, ""},
			{`(defun inner ()
				(labels ((f () "innermost"))
					#^(apply f %&rest)))`, `()`, ""},
			{`(funcall (inner))`, `"innermost"`, ""},
			{`(funcall (middle))`, `"innermost"`, ""},
			{`(defun inner ()
				(labels ((f () "innermost"))
					(lambda (&rest args) (apply #'f args))))`, `()`, ""},
			{`(funcall (inner))`, `"innermost"`, ""},
			{`(funcall (middle))`, `"innermost"`, ""},
		}},
		{"apply", elpstest.TestSequence{
			{`(apply (lambda () 1) '())`, `1`, ""},
			{`(apply (lambda (x) (+ 1 x)) '(1))`, `2`, ""},
			{`(apply (lambda (x) (+ 1 x)) 1 '())`, `2`, ""},
			{`(apply (lambda (&rest xs) (cdr xs)) 1 2 '(3))`, `'(2 3)`, ""},
			{`(apply '+ 1 2 '(3))`, `6`, ""},
			{`(apply 'lisp:+ 1 2 3 '())`, `6`, ""},
		}},
		{"variadic arguments", elpstest.TestSequence{
			{"((lambda (&rest x) (reverse 'list x)) 1 2 3)", "'(3 2 1)", ""},
			{"((lambda (x &rest y) (cons x y)) 1 2 3)", "'(1 2 3)", ""},
			{"((lambda (x &rest y) (cons x y)) 1 2 3)", "'(1 2 3)", ""},
			{"((lambda (x y &rest z) (cons x (cons y z))) 1 2 3)", "'(1 2 3)", ""},
			{"((lambda (x y &rest z) (cons x (cons y z))) 1 2)", "'(1 2)", ""},
		}},
		{"optional arguments", elpstest.TestSequence{
			{"((lambda (&optional x) (cons 1 x)))", "'(1)", ""},
			{"((lambda (&optional x) (cons 1 x)) '(2))", "'(1 2)", ""},
			{"((lambda (&optional x y) (+ (or x 1) (or y 2))))", "3", ""},
			{"((lambda (&optional x y) (+ (or x 1) (or y 2))) 2)", "4", ""},
			{"((lambda (&optional x y) (+ (or x 1) (or y 2))) 2 3)", "5", ""},
			{"((lambda (r &optional x) (cons r (cons 1 x))) 0)", "'(0 1)", ""},
			{"((lambda (r &optional x) (cons r (cons 1 x))) 0 '(2))", "'(0 1 2)", ""},
			{"((lambda (r &optional x y) (+ r (or x 1) (or y 2))) 1)", "4", ""},
			{"((lambda (r &optional x y) (+ r (or x 1) (or y 2))) 1 2)", "5", ""},
			{"((lambda (r &optional x y) (+ r (or x 1) (or y 2))) 1 2 3)", "6", ""},
		}},
		{"keyword arguments", elpstest.TestSequence{
			{"((lambda (&key x) (reverse 'list x)))", "'()", ""},
			{"((lambda (&key x) (reverse 'list x)) :x '(1 2 3))", "'(3 2 1)", ""},
			{"((lambda (&key x y) (cons (or x 1) y)) :y '(2))", "'(1 2)", ""},
			{"((lambda (r &key x) (cons (or x 1) (reverse 'list r))) '(2 3))", "'(1 3 2)", ""},
			{"((lambda (r &key x) (cons (or x 1) (reverse 'list r))) '(2 3) :x 4)", "'(4 3 2)", ""},
			{"((lambda (r &key x y) (cons r (cons (or x 1) y))) 0 :y '(2))", "'(0 1 2)", ""},
		}},
		//{"partial evaluation", elpstest.TestSequence{
		//	{"((lambda (x y) (+ x y)) 1)", "(lambda (y (x 1)) (+ x y))", ""},
		//	{"(((lambda (x y) (+ x y)) 1) 2)", "3", ""},
		//}},
		{"lists", elpstest.TestSequence{
			{"(cons 1 (cons 2 (cons 3 ())))", "'(1 2 3)", ""},
			{"(list 1 2 3)", "'(1 2 3)", ""},
			{"(concat 'list (list 1 2) (list 3))", "'(1 2 3)", ""},
			{"(cons 1 (cons 2 (cons 3 ())))", "'(1 2 3)", ""},
			{"(list 1 2 3)", "'(1 2 3)", ""},
			{"(reverse 'list (list 1 2 3))", "'(3 2 1)", ""},
			{"(reverse 'list (list 1 2))", "'(2 1)", ""},
			{"(reverse 'vector (list 1 2 3))", "(vector 3 2 1)", ""},
			{"(reverse 'vector (list 1 2))", "(vector 2 1)", ""},
			{"(reverse 'list (vector 1 2 3))", "'(3 2 1)", ""},
			{"(reverse 'list (vector 1 2))", "'(2 1)", ""},
			{"(reverse 'vector (vector 1 2 3))", "(vector 3 2 1)", ""},
			{"(reverse 'vector (vector 1 2))", "(vector 2 1)", ""},
			{"(concat 'list (list 1 2) (list 3))", "'(1 2 3)", ""},
			{"(slice 'list '(0 1 2 3 4) 1 3)", "'(1 2)", ""},
		}},
		{"insert-index", elpstest.TestSequence{
			{`(insert-index 'list '() 0 1)`, `'(1)`, ""},
			{`(insert-index 'list '(2) 0 1)`, `'(1 2)`, ""},
			{`(insert-index 'list '(1) 1 2)`, `'(1 2)`, ""},
			{`(insert-index 'list '(1 3) 1 2)`, `'(1 2 3)`, ""},
			{`(insert-index 'vector (vector) 0 1)`, `(vector 1)`, ""},
			{`(insert-index 'vector (vector 2) 0 1)`, `(vector 1 2)`, ""},
			{`(insert-index 'vector (vector 1) 1 2)`, `(vector 1 2)`, ""},
			{`(insert-index 'vector (vector 1 3) 1 2)`, `(vector 1 2 3)`, ""},
		}},
		{"append 'list", elpstest.TestSequence{
			{"(set 'v (list))", "'()", ""},
			{"(set 'v1 (append 'list v 1))", "'(1)", ""},
			{"(set 'v12 (append 'list v1 2))", "'(1 2)", ""},
			{"(set 'v123 (append 'list v12 3))", "'(1 2 3)", ""},
			{"(set 'v1234 (append 'list v123 4))", "'(1 2 3 4)", ""},
			{"v", "'()", ""},
			{"v1", "'(1)", ""},
			{"v12", "'(1 2)", ""},
			{"v123", "'(1 2 3)", ""},
			{"v1234", "'(1 2 3 4)", ""},
			{"(set 'v1235 (append 'list v123 5))", "'(1 2 3 5)", ""},
			// There is never any slice memory shared between the argument and
			// list return values of append, so the previous append cannot
			// modify the previously computed value of v1234 lists do not
			// behave like go slices in this way and are often a poor choice
			// for an expanding buffer.
			{"v1234", "'(1 2 3 4)", ""},
		}},
		{"make-sequence", elpstest.TestSequence{
			{"(make-sequence 0 5)", "'(0 1 2 3 4)", ""},
			{"(make-sequence 0 5 2)", "'(0 2 4)", ""},
		}},
		{"filtering", elpstest.TestSequence{
			{"(select 'list #^(< % 3) '())", "'()", ""},
			{"(select 'list #^(< % 3) '(0 1 2 3 4 5))", "'(0 1 2)", ""},
			{"(select 'list #^(< % 3) '(3 4 5 6))", "'()", ""},
			{"(reject 'list #^(< % 3) '())", "'()", ""},
			{"(reject 'list #^(< % 3) '(0 1 2 3 4 5))", "'(3 4 5)", ""},
			{"(reject 'list #^(< % 3) '(0 1 1 -1 2 2))", "'()", ""},
			{"(select 'list #^(< % 3) (vector))", "'()", ""},
			{"(select 'list #^(< % 3) (vector 0 1 2 3 4 5))", "'(0 1 2)", ""},
			{"(select 'list #^(< % 3) (vector 3 4 5 6))", "'()", ""},
			{"(reject 'list #^(< % 3) (vector ))", "'()", ""},
			{"(reject 'list #^(< % 3) (vector 0 1 2 3 4 5))", "'(3 4 5)", ""},
			{"(reject 'list #^(< % 3) (vector 0 1 1 -1 2 2))", "'()", ""},
			{"(select 'vector (expr (< % 3)) '())", "(vector)", ""},
			{"(select 'vector (expr (< % 3)) '(0 1 2 3 4 5))", "(vector 0 1 2)", ""},
			{"(select 'vector (expr (< % 3)) '(3 4 5 6))", "(vector)", ""},
			{"(reject 'vector (expr (< % 3)) '())", "(vector)", ""},
			{"(reject 'vector (expr (< % 3)) '(0 1 2 3 4 5))", "(vector 3 4 5)", ""},
			{"(reject 'vector (expr (< % 3)) '(0 1 1 -1 2 2))", "(vector)", ""},
		}},
		{"defun", elpstest.TestSequence{
			// defun macro
			{"(defun fn0 () (+ 1 1))", "()", ""},
			{"(defun fn1 (n) (+ n 1))", "()", ""},
			{"(defun fn2 (x y) (+ x y))", "()", ""},
			{"(fn0)", "2", ""},
			{"(fn1 1)", "2", ""},
			{"(fn2 1 2)", "3", ""},
		}},
		{"errors", elpstest.TestSequence{
			{`(list 1 2 (error 'test-error "test message") 4)`, "test:1: test-error: test message", ""},
		}},
	}
	elpstest.RunTestSuite(t, tests)
}
