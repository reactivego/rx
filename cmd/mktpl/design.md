# Make Template from Example

I want to be able to use the generated reactivex.go file as the template for
generating another reactivex.go with additional types supported.
Don't want to have to convert reactivex.go to a template myself every time.

In order to do that I need to either do a text analysis on the code or
alternatively I can use the go/parser to produce an ast.

Using an ast I can analyze it to produce a template that I can then expand.
Or I could prune parts of the ast and then duplicate and attach them to a new
ast that I then write out back to source. So second aproach would not be using
a template expander but a custom written ast manipulator.
