# goxpath
An XPath implementation in Go.

###Axii supported:
    ancestor
    ancestor-or-self
    attribute
    child
    descendant
    descendant-or-self
    parent
    self

###Axii not supported (yet):
    following
    following-sibling
    preceding
    preceding-sibling
    namespace

###NodeTypes supported:
    node()

###NodeTypes not supported (yet):
    comment()
    text()
    processing-instruction()

###Shorthand's supported
    .
    ..
    @

###Installation
    go get github.com/ChrisTrenkamp/goxpath

###Example


#####test.xml
    <?xml version="1.0" encoding="UTF-8"?>
    <p1>
      <p2>
        <p3/>
      </p2>
      <p2>
        <p3/>
      </p2>
    </p1>

#####Absolute path
    $ goxpath '/p1' test.xml
    <p1>
      <p2>
        <p3></p3>
      </p2>
      <p2>
        <p3></p3>
      </p2>
    </p1>

    $ goxpath '/p1/p2/p3' test.xml
    <p3></p3>
    <p3></p3>

#####Abbreviated Relative path
    $ goxpath '//p2' test.xml
    <p2>
        <p3></p3>
      </p2>
    <p2>
        <p3></p3>
      </p2>

    $ goxpath '//p3' test.xml
    <p3></p3>
    <p3></p3>

###API

####import
    import "github.com/ChrisTrenkamp/goxpath/xpath"

####Usage
    res, _ := xpath.FromStr(xp, x)
    for i := range res {
        str, _ := xpath.Print(res[i])
        fmt.Println(str)
    }
