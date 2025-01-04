# Glossary

## API

An API is an importable symbol from a (standard) library, e. g. `node:fs.mkdir`,
`net/http.NewRequest` or `urllib.parse`.

## Catalogue

A catalogue is a collection of meta information. For [APIs](#API), it contains
the amount of APIs, the amount of namespaces, namespaces symbols and technlogy
version. For contributions, it contains the amount of contributions and
repositories.

## Contribution

A contribution is an open-source file that contains location information
([locus](#locus)), uniquely identifiable repository information, the source
code, file name, and file path.

## Locus

A locus is part of [contribution](#contribution) and consists of an [API](#API) and
a line number of its occurrence.