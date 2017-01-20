package news
import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/midwares"
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/services/newsservice"
    "prosnav.com/wxserver/modules/log"
    "encoding/json"
    "strings"
    "prosnav.com/wxserver/modules/setting"
)


var KEYSECRET = setting.Cfg.Section("auth").Key("KEYSECRET").String()

func QueryNewsByUserid(c *gin.Context) {
    sess := midwares.GetSession(c)
    user := sess.Get("user").(domain.User)
    news := newsservice.QueryNewsByUserid(user.UserId)

    c.JSON(200, news)
}


type NewsForm struct {
    Content   string   `form:"content"`
    Summary   string   `form:"summary"`
    Title     string   `form:"title"`
    TagCode string  `form:"tagcode"`
    PageCount int     `form:"pageCount" json:"pageCount"`
    CurrPage  int     `form:"currPage" json:"currPage"`
    UserName  string    `json:"username"`
    Total     int64       `json:"total"`
    News  []*newsservice.News `json:"news"`
}


func QueryNewsByUseridWithPage(c *gin.Context) {
    sess := midwares.GetSession(c)
    userid := sess.Get("user").(string)
    username := sess.Get("username").(string)
    var form NewsForm
    c.Bind(&form)
    if form.PageCount == 0 {
        form.PageCount = 20
    }
    log.Debug("form: %v", form)
    total := newsservice.CountNewsByUserid(userid)
    news := newsservice.QueryShortNewsByUseridWithPage(userid, form.CurrPage, form.PageCount)
    log.Debug("news length: %d", len(news))
    c.JSON(200, NewsForm{
        CurrPage: form.CurrPage,
        PageCount: form.PageCount,
        Total: total,
        News: news,
        UserName: username,
    })
}


func QueryNewsWithPage(c *gin.Context) {
    var form NewsForm
    c.Bind(&form)
    if form.PageCount == 0 {
        form.PageCount = 20
    }
    log.Debug("form: %v", form)
    total := newsservice.CountNews()
    news := newsservice.QueryShortNewsWithPage(form.CurrPage, form.PageCount)
    log.Debug("news length: %d", len(news))
    c.JSON(200, NewsForm{
        CurrPage: form.CurrPage,
        PageCount: form.PageCount,
        Total: total,
        News: news,
    })
}

func InsertNews(c *gin.Context) {
    var form NewsForm
    c.Bind(&form)
    news := new(domain.News)
    news.Content = form.Content
    tagCode := strings.TrimSpace(form.TagCode)
    news.TagCode = tagCode
    news.Summary = form.Summary
    news.Title = form.Title
    newsservice.InsertNews(news)
    c.JSON(200, nil)
}

func UpdateNews(c *gin.Context) {
    newsid := c.Params.ByName("newsid")
    var form NewsForm
    c.Bind(&form)
    news := new(domain.News)
    news.Id = newsid
    news.Content = form.Content
    tagCode := strings.TrimSpace(form.TagCode)
    news.TagCode = tagCode
    news.Summary = form.Summary
    news.Title = form.Title
    newsservice.UpdateNews(news)
    c.JSON(200, nil)
}

func ExpandNews(c *gin.Context) {
    newsid := c.Params.ByName("newsid")
    news := newsservice.ExpandNewsWithId(newsid)
    jsonstr, err := json.Marshal(news)
    if err != nil {
        panic(err)
    }
    log.Debug("news json string: %s", string(jsonstr))
    c.JSON(200, news)
}

func DeleteNewsById(c *gin.Context) {
    newsid := c.Params.ByName("newsid")
    newsservice.DeleteNews(newsid)
    c.JSON(200, nil)
}


