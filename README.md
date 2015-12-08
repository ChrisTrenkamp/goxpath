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
    import "github.com/ChrisTrenkamp/goxpath/xpath"

####Usage
    res, _ := xpath.FromStr(xp, x)
    for i := range res {
        str, _ := xpath.Print(res[i])
        fmt.Println(str)
    }
