# goxpath [![GoDoc](https://godoc.org/gopkg.in/src-d/go-git.v2?status.svg)](https://godoc.org/github.com/ChrisTrenkamp/goxpath) [![codecov.io](https://codecov.io/github/ChrisTrenkamp/goxpath/coverage.svg?branch=master)](https://codecov.io/github/ChrisTrenkamp/goxpath?branch=master)
An XPath implementation in Go.

###Installation
To retrieve the API:

    go get -u github.com/ChrisTrenkamp/goxpath

To install the command-line client:

    go get -u github.com/ChrisTrenkamp/goxpath/cmd/goxpath

###Command-line examples

#####montypython.xml
    <?xml version="1.0" encoding="UTF-8"?>
    <grail>
        <quest>
            <for>shrubbery</for>
        </quest>
        <knight xmlns="http://monty.python">
            <who say="ni!"/>
        </knight>
    </grail>

#####Absolute path
    $ goxpath '/grail/quest' montypython.xml
    <quest>
            <for>shrubbery</for>
        </quest>

#####Absolute path value
    $ goxpath -v '/grail/quest' montypython.xml

            shrubbery

#####Namespace mapping
    $ goxpath '/grail/knight' montypython.xml
    $ #Nothing is returned because 'knight' is not in a known namespace

    $ goxpath -ns monty=http://monty.python '/grail/monty:knight' montypython.xml
    <knight xmlns="http://monty.python">
            <who xmlns="http://monty.python" say="ni"></who>
        </knight>

    $ goxpath -v -ns monty=http://monty.python '/grail/monty:knight/monty:who/@say' montypython.xml
    ni!

###API

####import
    "github.com/ChrisTrenkamp/goxpath"
    "github.com/ChrisTrenkamp/goxpath/tree/xmltree"

####Usage
    xp := goxpath.MustParse(`/path`)
    t := xmltree.MustParseXML(bytes.NewBufferString(xml.Header + `<path>Hello world</path>`))
    res := goxpath.MustExec(xp, t, nil)
    fmt.Println(res[0]) //Hello world
    err := xmltree.Marshal(res[0], os.Stdout) //<path>Hello world</path>
