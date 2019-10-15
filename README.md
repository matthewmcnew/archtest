# Archtest

[![Build Status](https://travis-ci.org/matthewmcnew/archtest.svg?branch=master)](https://travis-ci.org/matthewmcnew/archtest)

Unit test your golang architecture

Complete examples in [archtest_test.go](archtest_test.go)

#### Checking for dependencies

```golang
archtest.Package(t, "github.com/myorg/myproject/....").
    ShouldNotDependOn("github.com/some/package")
```

#### Checking for direct dependencies

```golang
archtest.Package(t, "github.com/myorg/myproject/....").
    ShouldNotDependDirectlyOn("github.com/some/package")
```

#### Including Tests

```golang
archtest.Package(t, "github.com/myorg/myproject/....").
    IncludeTests().
    ShouldNotDependDirectlyOn("github.com/some/package")
```