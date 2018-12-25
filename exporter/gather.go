package exporter

import (
	"bytes"
	"encoding/json"
	"eos_exporter/config"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/common/log"
)

var (
	errResponseStatusIsNotOK = errors.New("Status code is not equal to 200")
)

type accountGatherResult struct {
	accountInfo *AccountInfo
	err         error
}

type eosRpcClient struct {
	httpClient  http.Client
	baseAddress string
}

type getAccountRequest struct {
	AccountName string `json:"account_name"`
}

type getCurrencyBalanceRequest struct {
	Account string `json:"account"`
	Code    string `json:"code"`
	Symbol  string `json:"symbol"`
}

func parseAsset(asset string) (quantity float64, symbol string) {
	fmt.Sscanf(asset, "%f %s", &quantity, &symbol)
	return
}

func (rpc eosRpcClient) getCurrencyBalance(account, code, symbol string) (balance float64, err error) {
	url := fmt.Sprintf("%s/v1/chain/get_currency_balance", rpc.baseAddress)
	req := getCurrencyBalanceRequest{account, code, symbol}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}
	resp, err := rpc.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return 0, errResponseStatusIsNotOK
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	assets := []string{}
	err = json.Unmarshal(bytes, &assets)
	if err != nil {
		return
	}
	if len(assets) == 0 {
		return
	}
	balance, _ = parseAsset(assets[0])
	return balance, nil
}

func (rpc eosRpcClient) getAccountInfo(account string) (info *AccountInfo, err error) {
	url := fmt.Sprintf("%s/v1/chain/get_account", rpc.baseAddress)

	req := getAccountRequest{account}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}
	resp, err := rpc.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return nil, errResponseStatusIsNotOK
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return
	}
	return info, nil
}

func (rpc eosRpcClient) requestAccountInfos(accounts []string, tokens []config.TokenContract) []*AccountInfo {
	ch := make(chan *accountGatherResult)
	jobsDone := 0
	for _, account := range accounts {
		go func(acc string) {
			accInfo, err := rpc.getAccountInfo(acc)
			if err != nil {
				ch <- &accountGatherResult{err: err}
				return
			}
			accInfo.CurrencyBalances = make(map[string]Float64)
			for _, tok := range tokens {
				balance, err := rpc.getCurrencyBalance(acc, tok.Account, tok.Symbol)
				if err != nil {
					log.Debug(err)
				}
				accInfo.CurrencyBalances[tok.Symbol] = Float64(balance)
			}
			ch <- &accountGatherResult{
				accountInfo: accInfo,
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
	rpc := eosRpcClient{
		httpClient:  httpClient,
		baseAddress: e.Config.APIURL,
	}
	return rpc.requestAccountInfos(e.Config.Accounts, e.Config.Tokens), nil
}
