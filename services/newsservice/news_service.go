package newsservice
import (
    "time"
    "prosnav.com/wxserver/utils"
    "prosnav.com/wxserver/domain"
//    "prosnav.com/wxserver/modules/log"
)

const(
    queryAllByUserId = `select c.id, c.title, c.summary from t_wc_user a, t_wc_tag b, t_wc_news c
        where (array[b.tagcode] <@ a.tags and c.tagcode = b.tagcode  and a.userid = $1 or b.tagcode = 'all') and c.del != $2 order by c.crt desc `
    newsCountByUserId = `select count(c.*) count from t_wc_user a, t_wc_tag b, t_wc_news c
        where array[b.tagcode] <@ a.tags and c.tagcode = b.tagcode  and a.userid = $1 and c.del != $2`

    newsCount = `select count(*) count from t_wc_news where del != $1`
    queryAll = `select id, title, summary, tagcode from t_wc_news where del != $1 order by crt desc`
)

type News struct {
    Id       string     `xorm:"id" json:"id"`
    Title    string     `xorm:"title" json:"title"`
    ProductCode string  `xorm:"productcode" json:"productcode"`
    Content  string     `xorm:"content" json:"content"`
    Summary  string     `xorm:"summary" json:"summary"`
    Crt   time.Time     `xorm:"crt timestamp created"`
    Del   bool          `xorm:"del"`
}

func QueryNewsByUserid(userid string) []*News {
    var news []*News
    err := utils.Engine.Sql(queryAllByUserId, userid).Limit(3000).Find(&news)
    if err != nil {
        panic(err)
    }

    return news
}

type Count struct {
    Num int64 `xorm:"count"`
}

func CountNewsByUserid(userid string ) int64 {
    var c []*Count
    err:=utils.Engine.Sql(newsCountByUserId, userid, true).Find(&c)
    if err != nil {
        panic(err)
    }
    return c[0].Num
}

/*
*   news list without content
*/
func QueryShortNewsByUseridWithPage(userid string, currPage, pageCount int) []*News{
    var news []*News
    offset := (currPage - 1) * pageCount
    err := utils.Engine.Sql(queryAllByUserId + " limit $3 offset $4", userid, true, pageCount, offset).Find(&news)
    if err != nil {
        panic(err)
    }
    return news
}

func CountNews() int64 {
    var c []*Count
    err:=utils.Engine.Sql(newsCount, true).Find(&c)
    if err != nil {
        panic(err)
    }
    return c[0].Num
}

func QueryShortNewsWithPage(currPage, pageCount int) []*News {
    var news []*News
    offset := (currPage - 1) * pageCount
    err := utils.Engine.Sql(queryAll + " limit $2 offset $3", true, pageCount, offset).Find(&news)
    if err != nil {
        panic(err)
    }
    return news
}

func ExpandNewsWithId(id string) *domain.News {
    news := new(domain.News)
    news.Id = id
    utils.Engine.Get(news)
    return news
}

func UpdateNews(news *domain.News) {
    if _, err :=utils.Engine.Id(news.Id).Update(news); err != nil {
        panic(err)
    }
}

func InsertNews(news *domain.News) {
    if _, err := utils.Engine.InsertOne(news); err != nil {
        panic(err)
    }
}

func DeleteNews(newsid string) {
    news := new(domain.News)
    news.Id = newsid
    utils.Engine.Id(newsid).UseBool("del").Update(&domain.News{Del:true})
}