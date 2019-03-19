# handle
This package provides a way to make the execution paths of http handlers easier to trace while eliminating the possibility of certain mistakes like responding twice to an http request. It does this by providing the `Response` type so that http responses can be treated primarily as return values instead of side-effects, as well as the `With` function for ease of use.

`handle` is intended to be small and modular. That is, you can decide whether to use this package or the more flexible streaming semantics of standard go on a per-route basis.

Package `respondwith` is also provided as an optional way to simplify some common types of responses, like strings or JSON.