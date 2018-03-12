/*
Package raws currently provides simplicity - one package vs multitude in AWS -
as well as multi-region management - all calls are done for each selected region(s).
Region's parameter also supports globbing, thus allowing to fetch data from all
eu with: 'eu-*' or all eu-west with 'eu-west-*'

Currently only a couple of the most used information is gathered, but adding extra
calls should not be complicated, as they all have the same logic.

For the sake of avoiding repetitive documentation, each function that this package
contains with the context.Context type as a first parameter, won't be documented,
at least if there is not special usage of it inside of the function, because the
context.Context is provided for implementing the most adopted concurrency pattern,
used by the Go community (see https://blog.golang.org/context).
*/
package raws
