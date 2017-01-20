package assessmentservice

import (
    "prosnav.com/wxserver/domain"
    "gopkg.in/ini.v1"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/utils"
    "prosnav.com/wxserver/modules/idg"
)

const (
	ALL_CUSTOMERS = `select name, mobile, email, score, crt from p_assessment where del = false`
	SEARCH_CUSTOMER = `select name, mobile, email, score, crt from p_assessment where del = false and name like ?`
)

var (
    Paper = make(map[string]map[string]string)
    Assessment *ini.File
)

type resultForm struct {
    Score int `json:"score"`
    Type string `json:"type"`
    Endurance string `json:"endurance"`
    Venture string `json:"venture"`
}

func Assess(assess *domain.Assessment) (*resultForm, error) {
    result := assessResult(assess)
    assess.Id, _ = idg.Id()
    assess.Score = result.Score
    _, err := utils.Engine.InsertOne(assess)
    return result, err
}

//保守型(25-35)、稳健型（36-45），平衡型（46-60），成长型（61-80）、进取型（81-100）
func assessResult(assess *domain.Assessment) *resultForm {
    result := new(resultForm)
    score := Assessment.Section("q1").Key(assess.Q1).MustInt() +
    Assessment.Section("q2").Key(assess.Q2).MustInt() +
    Assessment.Section("q3").Key(assess.Q3).MustInt() +
    Assessment.Section("q4").Key(assess.Q4).MustInt() +
    Assessment.Section("q5").Key(assess.Q5).MustInt() +
    Assessment.Section("q6").Key(assess.Q6).MustInt() +
    Assessment.Section("q7").Key(assess.Q7).MustInt() +
    Assessment.Section("q8").Key(assess.Q8).MustInt() +
    Assessment.Section("q9").Key(assess.Q9).MustInt() +
    Assessment.Section("q10").Key(assess.Q10).MustInt() +
    Assessment.Section("q11").Key(assess.Q11).MustInt() +
    Assessment.Section("q12").Key(assess.Q12).MustInt() +
    Assessment.Section("q13").Key(assess.Q13).MustInt() +
    Assessment.Section("q14").Key(assess.Q14).MustInt() +
    Assessment.Section("q15").Key(assess.Q15).MustInt() +
    Assessment.Section("q16").Key(assess.Q16).MustInt() +
    Assessment.Section("q17").Key(assess.Q17).MustInt() +
    Assessment.Section("q18").Key(assess.Q18).MustInt()
    result.Score = score
    if score <= 26 {
        result.Type = "保守型"
    	result.Endurance = "弱"
        result.Venture = "低风险"
    } else if score <= 36{
        result.Type = "稳健型"
        result.Endurance = "较弱"
        result.Venture = "中低风险、低风险"
    } else if score <= 60 {
        result.Type = "平衡型"
        result.Endurance = "中等"
        result.Venture = "中高风险、中风险、中低风险、低风险"
    } else if score <= 75 {
        result.Type = "成长型"
        result.Endurance = "较强"
    	result.Venture = "中高风险、中风险、中低风险、低风险"
    } else {
        result.Type = "进取型"
        result.Endurance = "强"
    	result.Venture = "高风险、中高风险、中风险、中低风险、低风险"
    }
    return result
}


func loadQuestions(paper *ini.File) {
    for _, sec := range paper.Sections() {
        if sec.Name() == "DEFAULT" {
            continue
        }
        Paper[sec.Name()] = sec.KeysHash()
    }
}

func init() {
    paper, err := ini.Load("conf/paper.ini")
    loadQuestions(paper)
    if err != nil {
        log.Error(0, "Load configuration file paper.ini failed %v", err)
        panic(err)
    }
    Assessment, err = ini.Load("conf/assessment.ini")
    if err != nil {
        log.Error(0, "Load configuration file assessment.ini failed %v", err)
        panic(err)
    }
}

func GetCustsInfo() (customers []*domain.Assessment, err error)  {
	err = utils.Engine.Sql(ALL_CUSTOMERS).Find(&customers); if err != nil {
		log.Error(0, "select customers info error %v", err)
		return nil, err
	}
	return customers, nil
}

func GetCustInfo (name string) (customer []*domain.Assessment, err error)  {
	err = utils.Engine.Sql(SEARCH_CUSTOMER, name).Find(&customer); if err != nil {
		log.Error(0, "search customer info error %v", err)
		return nil, err
	}
	return customer, nil
}