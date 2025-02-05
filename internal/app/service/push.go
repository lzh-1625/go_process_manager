package service

import (
	"net/http"
	"strings"

	"github.com/lzh-1625/go_process_manager/internal/app/repository"

	"github.com/levigross/grequests"
)

type pushService struct{}

var PushService = new(pushService)

func (p *pushService)  Push(ids []int, placeholders map[string]string) {
	pl := repository.PushRepository.GetPushConfigByIds(ids)
	for _, v := range pl {
		if v.Enable {
			if v.Method == http.MethodGet {
				grequests.Get(p.getReplaceMessage(placeholders, v.Url), nil)
			}
			if v.Method == http.MethodPost {
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
