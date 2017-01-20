package sync
import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/services/syncservice"
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/modules/log"
)

type Adapter interface {
    adapt() interface{}
}

type userForm struct {
    Userid string   `json:"userid"`
    Product []string    `json:"product"`
}

type proForm struct {
    ProductCode string  `json:"productcode"`
    Brief string        `json:"brief"`
}

func (a *userForm)adapt() *domain.User {
    user := new(domain.User)
    user.UserId = a.Userid
    user.Tags = a.Product
    return user
}

func (a *proForm)adapt() *domain.Tag {
    tag := new(domain.Tag)
    tag.TagCode = a.ProductCode
    tag.Brief = a.Brief
    return tag
}

func SyncUsers(c *gin.Context) {
    var forms []*userForm
    c.Bind(&forms)
    users := make([]*domain.User, len(forms))
    for i, u := range forms {
        log.Debug("form: %v", u)
        users[i] = u.adapt()
    }
    resultMap := syncservice.SyncUsers(users)
    c.JSON(200, resultMap)
}



func SyncProducts(c *gin.Context) {
    var forms []*proForm
    c.Bind(&forms)
    pros := make([]*domain.Tag, len(forms))
    for i, p := range forms  {
        pros[i] = p.adapt()
    }
    resultMap := syncservice.SyncTags(pros)
    c.JSON(200, resultMap)
}
