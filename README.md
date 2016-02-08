# goxpath
An XPath implementation in Go.

###Installation
    go get github.com/ChrisTrenkamp/goxpath

###Example

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
    "github.com/ChrisTrenkamp/goxpath/goxpath"
    "github.com/ChrisTrenkamp/goxpath/tree/xmltree"

####Usage
    xp := goxpath.MustParse(`/path`)
    t := xmltree.MustParseXML(bytes.NewBufferString(xml.Header + `<path>Hello world</path>`))
    res := goxpath.MustExec(xp, t, nil)
    fmt.Println(res[0]) //Hello world
    err := xmltree.Marshal(res[0], os.Stdout) //<path>Hello world</path>
