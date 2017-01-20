package results
import (
    "fmt"
    "github.com/gin-gonic/gin"
    "runtime"
    "path"
    "prosnav.com/wxserver/modules/log"
    "encoding/json"
)

const (
    UNKNOWN_ERROR = "UnknownError"
)

type BaseError struct {
}

func (e *BaseError) Error() string {
    data, _ := json.Marshal(e)
    return string(data)
}

type BusinessError struct {
    BaseError
    ErrCode string
    ErrMsg  string
}
func (e *BusinessError) Error() string {
    data, _ := json.Marshal(e)
    return string(data)
}

type BadRequestError struct {
    BaseError
}
func (e *BadRequestError) Error() string {
    data, _ := json.Marshal(e)
    return string(data)
}

type UnAuthorizedError struct {
    BaseError
}
func (e *UnAuthorizedError) Error() string {
    data, _ := json.Marshal(e)
    return string(data)
}
type NotAllowedError struct {
    BaseError
}

func (e *NotAllowedError) Error() string {
    data, _ := json.Marshal(e)
    return string(data)
}
type NotFoundError struct {
    BaseError
}

func (e *NotFoundError) Error() string {
    data, _ := json.Marshal(e)
    return string(data)
}
func NewBaseError(err error) *BusinessError {
    baseErr := new(BusinessError)
    baseErr.ErrCode = "5001"
    baseErr.ErrMsg = err.Error()
    return baseErr
}
func NewBusinessError(errCode string) *BusinessError {
    busErr := new(BusinessError)
    busErr.ErrCode = errCode
    return busErr
}

func stack(skip int) string {
    stk := make([]uintptr, 32)
    str := ""
    l := runtime.Callers(skip, stk[:])
    for i := 0; i < l; i++ {
        f := runtime.FuncForPC(stk[i])
        name := f.Name()
        file, line := f.FileLine(stk[i])
        str += fmt.Sprintf("\n    %-30s [%s:%d]", name, path.Base(file), line)
    }
    return str
}

func ErrorHandler() gin.HandlerFunc{
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Error(3, "Error: %v\n%v", err, stack(3))
                switch e := err.(type) {
                    case *BadRequestError:
                        c.JSON(400, e)
                        c.Abort()
                    case *UnAuthorizedError:
                        c.JSON(401, e)
                        c.Abort()
                    case *NotAllowedError:
                        c.JSON(403, e)
                        c.Abort()
                    case *NotFoundError:
                        c.JSON(404, e)
                        c.Abort()
                    case *BusinessError:
                        c.JSON(500, e)
                        c.Abort()
                    default:
                        c.JSON(500, NewBaseError(e.(error)))
                        c.Abort()
                }
            }
        }()

        c.Next()
    }
}


