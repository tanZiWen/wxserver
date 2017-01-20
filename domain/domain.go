package domain

import "time"

type User struct {
	Id         string   `xorm:"pk <-"`
	UserId     string   `xorm:"userid" json:"userid"`
	Name       string      `xorm:"name" json:"name"`
	Department string `xorm:"department" json:"department"`
	Mobile     string   `xorm:"mobile" json:"mobile"`
	Email      string   `xorm:"email" json:"email"`
	Gender     string   `xorm:"gender" json:"gender"`
	WeixinId   string   `xorm:"weixinid" json:"weixinid"`
	Tags       []string `xorm:"tags" json:"tags"`
	Status     string   `xorm:"status" json:"status"`
	avatar     string   `xorm:"avatar" json:"avatar"`
	Crt        time.Time   `xorm:"crt timestamp created"`
	Lut        time.Time   `xorm:"lut timestamp updated"`
	Del        bool    `xorm:"del"`
}

func (c *User) TableName() string {
	return "t_wc_user"
}

type Tag struct {
	Id      string   `xorm:"pk <-"`
	TagCode string `xorm:"tagcode"`
	Brief   string      `xorm:"brief"`
	Crt     time.Time   `xorm:"crt timestamp created"`
	Lut     time.Time   `xorm:"lut timestamp updated"`
	Del     bool    `xorm:"del"`
}

func (c *Tag) TableName() string {
	return "t_wc_tag"
}

type News struct {
	Id      string    `xorm:"pk <-" json:"id"`
	Title   string        `xorm:"title" json:"title"`
	Summary string        `xorm:"summary" json:"summary"`
	Content string        `xorm:"content" json:"content"`
	TagCode string    `xorm:"tagcode" json:"tagcode"`
	Crt     time.Time    `xorm:"crt timestamp created" json:"crt"`
	Del     bool    `xorm:"del"`
}

func (c *News) TableName() string {
	return "t_wc_news"
}

type Admin struct {
	Id     string   `xorm:"pk <-"`
	UserId string   `xorm:"userid"`
	Passwd string   `xorm:"passwd"`
	Name   string      `xorm:"name"`
	Crt    time.Time   `xorm:"crt timestamp created"`
	Lut    time.Time   `xorm:"lut timestamp updated"`
	Del    bool    `xorm:"del"`
}

func (c *Admin) TableName() string {
	return "t_wc_admin"
}

type Activity struct {
	Id           int64 `xorm:"id pk" json:"id,string"`
	Latitude     float64   `xorm:"latitude" json:"-"`
	Longitude    float64   `xorm:"longitude" json:"-"`
	Address      string `xorm:"address" json:"address"`
	ActivityName string `xorm:"activityname" json:"activityname"`
	Brief        string `xorm:"brief" json:"-"`
	StartDate    time.Time `xorm:"startdate timestamp" json:"-"`
	EndDate      time.Time `xorm:"enddate timestamp" json:"-"`
	TagId        string `xorm:"tagid" json:"-"`
	AgentId      int    `xorm:"agentid" json:"-"`
	MediaId      string `xorm:"mediaid"`
	Crt          time.Time `xorm:"crt timestamp created" json:"-"`
	Lut          time.Time `xorm:"lut timestamp updated" json:"-"`
	Del          bool      `xorm:"del" json:"-"`
}

func (c *Activity) TableName() string {
	return "p_activity"
}

type SignIn struct {
	Id         int64 `xorm:"id pk"`
	ActivityId int64 `xorm:"activityid" json:"activityid"`
	UserId     string `xorm:"userid" json:"userid"`
	Custname   string `xorm:"custname" json:"custname"`
	Mobile     string `xorm:"mobile" json:"mobile"`
	Crt        time.Time `xorm:"crt timestamp created"`
	Lut        time.Time `xorm:"lut timestamp updated"`
	Status     int `xorm:"status"`
}

func (c *SignIn) TableName() string {
	return "p_signin"
}

type Assessment struct {
	Id     int64 `xorm:"id pk"`
	Name   string `xorm:"name" form:"name" json:"name"`
	Mobile string `xorm:"mobile" form:"mobile" json:"mobile"`
	Email  string `xorm:"email" form:"email" json:"email"`
	Q1     string `xorm:"q1" form:"q1" json:"q1"`
	Q2     string `xorm:"q2" form:"q2" json:"q2"`
	Q3     string `xorm:"q3" form:"q3" json:"q3"`
	Q4     string `xorm:"q4" form:"q4" json:"q4"`
	Q5     string `xorm:"q5" form:"q5" json:"q5"`
	Q6     string `xorm:"q6" form:"q6" json:"q6"`
	Q7     string `xorm:"q7" form:"q7" json:"q7"`
	Q8     string `xorm:"q8" form:"q8" json:"q8"`
	Q9     string `xorm:"q9" form:"q9" json:"q9"`
	Q10    string `xorm:"q10" form:"q10" json:"q10"`
	Q11    string `xorm:"q11" form:"q11" json:"q11"`
	Q12    string `xorm:"q12" form:"q12" json:"q12"`
	Q13    string `xorm:"q13" form:"q13" json:"q13"`
	Q14    string `xorm:"q14" form:"q14" json:"q14"`
	Q15    string `xorm:"q15" form:"q15" json:"q15"`
	Q16    string `xorm:"q16" form:"q16" json:"q12"`
	Q17    string `xorm:"q17" form:"q17" json:"q17"`
	Q18    string `xorm:"q18" form:"q18" json:"q18"`
	Score  int `xorm:"score" json:"score"`
	Crt    time.Time `xorm:"crt timestamp created" json:"crt"`
	Del    bool      `xorm:"del" json:"-"`
}

func (c *Assessment) TableName() string {
	return "p_assessment"
}

type Appoint struct {
	Id     int64 `xorm:"id pk"`
	Mobile string `xorm:"mobile" json:"mobile"`
	Name   string `xorm:"name" json:"name"`
	Note   string `xorm:"note" json:"note"`
	Number int `xorm:"number" json:"number"`
	Crt    time.Time `xorm:"crt timestamp created"`
	Lut    time.Time `xorm:"lut timestamp updated"`
	Status int `xorm:"status"`
}

func (c *Appoint) TableName() string {
	return "app_user"
}

type Reservation struct {
	Id     int64 `xorm:"id pk"`
	Mobile string `xorm:"mobile" json:"mobile"`
	Name   string `xorm:"name" json:"name"`
	Doctor string `xorm:"doctor" json:"doctor"`
	Time   string `xorm:"time" json:"time"`
	Crt    time.Time `xorm:"crt timestamp created"`
	Lut    time.Time `xorm:"lut timestamp updated"`
	Status int `xorm:"status"`
}

type ReservationList struct {
	Mobile string `xorm:"mobile" json:"mobile"`
	Name   string `xorm:"name" json:"name"`
	Doctor string `xorm:"doctor" json:"doctor"`
	Time   string `xorm:"time" json:"time"`
}

func (c *Reservation) TableName() string {
	return "longhua"
}