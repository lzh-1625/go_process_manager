package push

import (
	"msm/dao"
	"strings"

	"github.com/levigross/grequests"
)

type pushService struct{}

var PushService = new(pushService)

func (p *pushService) Push(placeholders map[string]string) {
	pl := dao.PushDao.GetPushList()
	for _, v := range pl {
		if v.Enable {
			if v.Method == "GET" {
				grequests.Get(p.getReplaceMessage(placeholders, v.Url), nil)
			}
			if v.Method == "POST" {
				grequests.Post(v.Url, &grequests.RequestOptions{
					JSON: p.getReplaceMessage(placeholders, v.Body),
				})
			}
		}
	}
}

func (p *pushService) getReplaceMessage(placeholders map[string]string, message string) string {
	for k, v := range placeholders {
		message = strings.ReplaceAll(message, k, v)
	}
	return message
}
