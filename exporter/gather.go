package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/common/log"
)

type Response struct {
	accountInfo *AccountInfo
	err         error
}

func requestAccountInfos(httpClient http.Client, url string, accounts []string) []*AccountInfo {
	ch := make(chan *Response)
	jobsDone := 0
	for _, account := range accounts {
		go func(acc string) {
			body := []byte(fmt.Sprintf("{\"account_name\":\"%s\"}", acc))
			resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(body))
			if err != nil {
				ch <- &Response{err: fmt.Errorf("Error sending request, Error %s", err)}
				return
			}
			if resp.StatusCode != 200 {
				ch <- &Response{err: fmt.Errorf("Error response status code, Status code is %d", resp.StatusCode)}
				return
			}
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ch <- &Response{err: fmt.Errorf("Error reading response body, Error %s", err)}
				return
			}
			accountInfo := &AccountInfo{}
			err = json.Unmarshal(bytes, &accountInfo)
			if err != nil {
				ch <- &Response{err: fmt.Errorf("Error unmarshalling response body, Error %s", err)}
				return
			}
			ch <- &Response{
				accountInfo: accountInfo,
				err:         nil,
			}
		}(account)
	}
	accountInfos := []*AccountInfo{}
	for {
		select {
		case r := <-ch:
			if r.err != nil {
				log.Errorf("Error scrapping API, Error %s", r.err)
			} else {
				accountInfos = append(accountInfos, r.accountInfo)
			}
			jobsDone++
			if jobsDone == len(accounts) {
				return accountInfos
			}
		}
	}
}

func (e *Exporter) gatherAccountData() ([]*AccountInfo, error) {
	if len(e.Config.Accounts) == 0 {
		return nil, nil
	}
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	getAccountURL := fmt.Sprintf("%s/v1/chain/get_account", e.Config.APIURL)
	return requestAccountInfos(httpClient, getAccountURL, e.Config.Accounts), nil
}
