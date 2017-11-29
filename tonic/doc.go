/*
Package tonic lets you write simpler gin handlers.
The way it works is that it generates wrapping gin-compatible handlers,
that do all the repetitive work and wrap the call to your simple tonic handler.

Package tonic handles path/query/body parameter binding in a single consolidated input object
which allows you to remove all the boilerplate code that retrieves and tests the presence
of various parameters.

Here is an example input object.

    type MyInput struct {
        Foo int    `path:"foo,required"`
        Bar float  `query:"bar,default=foobar"`
        Baz string `json:"baz" binding:"required"`
    }

Output objects can be of any type, and will be marshaled to JSON.

The handler can return an error, which will be returned to the caller.

Here is a basic application that greets a user on http://localhost:8080/hello/me

    import (
         "errors"
         "fmt"

         "github.com/gin-gonic/gin"
         "github.com/loopfz/gadgeto/tonic"
    )

    type GreetUserInput struct {
        Name string `path:"name,required" description:"User name"`
    }

    type GreetUserOutput struct {
        Message string `json:"message"`
    }

    func GreetUser(c *gin.Context, in *GreetUserInput) (*GreetUserOutput, error) {
        if in.Name == "satan" {
            return nil, errors.New("go to hell")
        }
        return &GreetUserOutput{Message: fmt.Sprintf("Hello %s!", in.Name)}, nil
    }

    func main() {
        r := gin.Default()
        r.GET("/hello/:name", tonic.Handler(GreetUser, 200))
        r.Run(":8080")
    }


If needed, you can also override different parts of the logic via certain available hooks in tonic:
    - binding
    - error handling
    - render

We provide defaults for these (bind from JSON, render into JSON, error = http status 400).
You will probably want to customize the error hook, to produce finer grained error status codes.

The role of this error hook is to inspect the returned error object and deduce the http specifics from it.
We provide a ready-to-use error hook that depends on the juju/errors package (richer errors):
    https://github.com/loopfz/gadgeto/tree/master/tonic/utils/jujerr

Example of the same application as before, using juju errors:

    import (
         "fmt"

         "github.com/gin-gonic/gin"
         "github.com/juju/errors"
         "github.com/loopfz/gadgeto/tonic"
         "github.com/loopfz/gadgeto/tonic/utils/jujerr"
    )

    type GreetUserInput struct {
        Name string `path:"name,required" description:"User name"`
    }

    type GreetUserOutput struct {
        Message string `json:"message"`
    }

    func GreetUser(c *gin.Context, in *GreetUserInput) (*GreetUserOutput, error) {
        if in.Name == "satan" {
            return nil, errors.NewForbidden(nil, "go to hell")
        }
        return &GreetUserOutput{Message: fmt.Sprintf("Hello %s!", in.Name)}, nil
    }

    func main() {
        tonic.SetErrorHook(jujerr.ErrHook)
        r := gin.Default()
        r.GET("/hello/:name", tonic.Handler(GreetUser, 200))
        r.Run(":8080")
    }


You can also easily serve auto-generated (using tonic data) swagger documentation:

    import (
         "fmt"

         "github.com/gin-gonic/gin"
         "github.com/juju/errors"
         "github.com/loopfz/gadgeto/tonic"
         "github.com/loopfz/gadgeto/tonic/utils/jujerr"
         "github.com/loopfz/gadgeto/tonic/utils/swag"
    )

    func main() {
        tonic.SetErrorHook(jujerr.ErrHook)
        r := gin.Default()
        r.GET("/hello/:name", tonic.Handler(GreetUser, 200))
        r.GET("/swagger.json", swag.Swagger(r, "MyAPI", swag.Version("v1.0"), swag.BasePath("/foo/bar")))
        r.Run(":8080")
    }
*/
package tonic
