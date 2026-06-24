# go-radir

[![Go Report Card](https://goreportcard.com/badge/github.com/oid-directory/go-radir)](https://goreportcard.com/report/github.com/oid-directory/go-radir) [![codecov](https://codecov.io/github/oid-directory/go-radir/graph/badge.svg?token=WE2QWITRXA)](https://codecov.io/github/oid-directory/go-radir) [![Go Reference](https://pkg.go.dev/badge/github.com/oid-directory/go-radir.svg)](https://pkg.go.dev/github.com/oid-directory/go-radir) ![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square) ![Volatility Warning](https://img.shields.io/badge/volatile-darkred?label=%F0%9F%92%A5&labelColor=white&color=orange&cacheSeconds=86400) [![RADIR](https://img.shields.io/badge/RA_DIR_(I--D)-red?style=flat-square&label=%F0%9F%93%84&cacheSeconds=86400&link=https%3A%2F%2Fdatatracker.ietf.org%2Fdoc%2Fhtml%2Fdraft-coretta-oiddir-roadmap)](https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-roadmap) [![RADIT](https://img.shields.io/badge/RA_DIT_(I--D)-red?style=flat-square&label=%F0%9F%93%84&cacheSeconds=86400&link=https%3A%2F%2Fdatatracker.ietf.org%2Fdoc%2Fhtml%2Fdraft-coretta-oiddir-radit)](https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radit) [![RADUA](https://img.shields.io/badge/RA_DUA_(I--D)-red?style=flat-square&label=%F0%9F%93%84&cacheSeconds=86400&link=https%3A%2F%2Fdatatracker.ietf.org%2Fdoc%2Fhtml%2Fdraft-coretta-oiddir-radua)](https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radua) [![RADSA](https://img.shields.io/badge/RA_DSA_(I--D)-red?style=flat-square&label=%F0%9F%93%84&cacheSeconds=86400&link=https%3A%2F%2Fdatatracker.ietf.org%2Fdoc%2Fhtml%2Fdraft-coretta-oiddir-radsa)](https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-radsa) [![RASCHEMA](https://img.shields.io/badge/RA_SCHEMA_(I--D)-red?style=flat-square&label=%F0%9F%93%84&cacheSeconds=86400&link=https%3A%2F%2Fdatatracker.ietf.org%2Fdoc%2Fhtml%2Fdraft-coretta-oiddir-schema)](https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema)


Package radir implements subsets of the ***EXPERIMENTAL*** OID Directory I-D (Internet-Draft) series.

## About

This package serves as the de facto Go-based OID-Directory library. It is the foundation of the entire OID-Directory Internet-Draft series, offering myriad functions, types and methods for use in creating an ASN.1 Object Identifier repository within an X.500 or LDAP architecture.

The go-radit (Directory Information Tree) and go-radua (Directory User Agent) packages rely upon go-radir.

## OID Utility

This package includes a subdirectory package known simply as "oid". It offers DotNotation, ASN1Notation and IRINotation types for expressing the various incarnations of Object Identifiers. See the package documentation for extensive information on making use of this package.

