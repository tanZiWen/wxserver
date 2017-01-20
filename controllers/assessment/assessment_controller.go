package assessment

import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/services/assessmentservice"
    "prosnav.com/wxserver/modules/log"
	"strings"
)

type resultForm struct {
    Score int `json:"score"`
    Type string `json:"type"`
}

func Assess(c *gin.Context) {
    form := new(domain.Assessment)
    c.Bind(form)
    result, err := assessmentservice.Assess(form)
    if err != nil {
        log.Error(4, "Assess error %v\n", err)
        panic(err)
    }
    c.JSON(200, result)
}

func QueryPapper(c *gin.Context) {
    c.JSON(200, assessmentservice.Paper)
}

/*获取客户的KYC信息*/
func CustomersInfo(c *gin.Context)  {
	var customers []*domain.Assessment
	var err error
	name := c.Query("name")
	log.Debug("name", name)
	if (strings.TrimSpace(name) == "") {
		customers, err = assessmentservice.GetCustsInfo(); if err != nil {
			log.Error(4, "failed to get customers", err)
			panic(err)
		}
	} else {
		customers, err = assessmentservice.GetCustInfo(name); if err != nil {
			log.Error(4, "failed to get customers", err)
			panic(err)
		}
	}
	c.JSON(200, customers)
}

