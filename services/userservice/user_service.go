package userservice
import (
	"prosnav.com/wxserver/utils"
	"fmt"
	"strings"
)

const (
	update_user = "update t_wc_user set tag = ARRAY_APPEND(tag, $1) where userid = $2"
	insert_user = "insert into t_wc_user (userid, product) values ($1, array[%s])"
)


func InsertUser(arr []string, tags []string, errMap *map[string]error){
	for _, v := range arr {
		_, err := utils.Engine.Exec(fmt.Sprintf(insert_user, strings.Join(tags, ", ")), v); if err != nil {
			(*errMap)[v] = err
		}
	}
}

func UpdateUser(arr []string, product string) error {
	for _, v := range arr {
		_, err := utils.Engine.Exec(update_user, product, v); if err != nil {
			fmt.Printf("update user error:%v", err);
			return err
		}
	}
	return nil
}
