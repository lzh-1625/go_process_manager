package logic

import (
	"io"
	"net/http"
	"strings"

	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/log"
)

type pushLogic struct {
	httpClient *http.Client
}

var PushLogic = &pushLogic{
	httpClient: &http.Client{
		Transport: http.DefaultTransport,
	},
}

func (p *pushLogic) Push(ids []int, placeholders map[string]string) {
	pl := repository.PushRepository.GetPushConfigByIds(ids)
	for _, v := range pl {
		if v.Enable {
			var resp *http.Response
			var reader io.Reader = nil
			var url string = v.Url
			if v.Method == http.MethodPost {
				reader = strings.NewReader(p.getReplaceMessage(placeholders, v.Body))
			}
			if v.Method == http.MethodGet {
				url = p.getReplaceMessage(placeholders, url)
			}
			req, err := http.NewRequest(v.Method, url, reader)
			if err != nil {
				log.Logger.Warnw("推送失败", "err", err, "remark", v.Remark)
				continue
			}
			req.Header.Add("content-type", "application/json")
			resp, err = p.httpClient.Do(req)
			if err != nil {
				log.Logger.Warnw("推送失败", "err", err, "remark", v.Remark)
				continue
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}
}

func (p *pushLogic) getReplaceMessage(placeholders map[string]string, message string) string {
	for k, v := range placeholders {
		message = strings.ReplaceAll(message, k, v)
	}
	return message
}
