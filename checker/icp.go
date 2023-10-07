package checker

//
//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"fmt"
//
//	"go.uber.org/zap"
//	"io"
//	"io/ioutil"
//	"math/rand"
//	"net/http"
//	"time"
//)
//
//var IcpNotForRecord = errors.New("域名未备案")
//
//type Icp struct {
//	token string
//	ip    string
//	agent string
//}
//
//func (i *Icp) Query(domain string) (*request.DomainInfo, error) {
//	i.getIp()
//	i.getUserAgent()
//	if err := i.auth(); err != nil {
//		return nil, errors.New("获取授权失败:" + err.Error())
//	}
//	return i.query(domain)
//}
//
//func (i *Icp) query(domain string) (*request.DomainInfo, error) {
//	queryRequest, _ := json.Marshal(&request.QueryRequest{
//		UnitName: domain,
//	})
//
//	result := &request.IcpResponse{Params: &request.QueryParams{}}
//	err := i.post("icpAbbreviateInfo/queryByCondition", bytes.NewReader(queryRequest), "application/json;charset=UTF-8", i.token, result)
//	if err != nil {
//		return nil, err
//	}
//
//	if !result.Success {
//		return nil, fmt.Errorf("查询：%s", result.Msg)
//	}
//
//	queryParams := result.Params.(*request.QueryParams)
//	if len(queryParams.List) == 0 {
//		return &request.DomainInfo{}, nil
//	}
//	return queryParams.List[0], nil
//}
//
//func (i *Icp) auth() error {
//	timestamp := time.Now().Unix()
//	authKey := utils.Md5(fmt.Sprintf("testtest%d", timestamp))
//	authBody := fmt.Sprintf("authKey=%s&timeStamp=%d", authKey, timestamp)
//
//	result := &request.IcpResponse{Params: &request.AuthParams{}}
//	err := i.post("auth", bytes.NewReader([]byte(authBody)), "application/x-www-form-urlencoded;charset=UTF-8", "0", result)
//	if err != nil {
//		return err
//	}
//
//	if !result.Success {
//		return fmt.Errorf("获取token失败：%s", result.Msg)
//	}
//
//	authParams := result.Params.(*request.AuthParams)
//	i.token = authParams.Bussiness
//	return nil
//}
//
//func (i *Icp) post(url string, data io.Reader, contentType string, token string, result interface{}) error {
//	postUrl := fmt.Sprintf("https://hlwicpfwc.miit.gov.cn/icpproject_query/api/%s", url)
//	queryReq, err := http.NewRequest(http.MethodPost, postUrl, data)
//	if err != nil {
//		return err
//	}
//	queryReq.Header.Set("Content-Type", contentType)
//	queryReq.Header.Set("Origin", "https://beian.miit.gov.cn/")
//	queryReq.Header.Set("Referer", "https://beian.miit.gov.cn/")
//	queryReq.Header.Set("token", token)
//	queryReq.Header.Set("User-Agent", i.agent)
//
//	client := DefaultProxyClient()
//	//client := http.DefaultClient
//	resp, err := client.Do(queryReq)
//
//	return GetHTTPResponse(resp, postUrl, err, result)
//}
//
//func (i *Icp) getIp() {
//	if i.ip != "" {
//		return
//	}
//	i.ip = RandIp()
//}
//
//func (i *Icp) getUserAgent() {
//	i.agent = RandAgent()
//}
//
//func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
//	if err != nil {
//		return err
//	}
//	body, err := GetHTTPResponseOrg(resp, url, err)
//	if err == nil {
//		err = json.Unmarshal(body, &result)
//		if err != nil {
//			return err
//		}
//	}
//	return err
//}
//
//func GetHTTPResponseOrg(resp *http.Response, url string, err error) ([]byte, error) {
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	// 300及以上状态码都算异常
//	if resp.StatusCode >= 300 {
//		errMsg := fmt.Sprintf("请求接口 %s 失败! 返回状态码: %d\n", url, resp.StatusCode)
//		err = fmt.Errorf(errMsg)
//	}
//	return body, err
//}
//
//func RandIp() string {
//	rand.Seed(time.Now().UnixNano())
//	return fmt.Sprintf("101.%d.%d.%d", 1+rand.Intn(254), 1+rand.Intn(254), 1+rand.Intn(254))
//}
//
//func RandAgent() string {
//	agents := Agents()
//	rand.Seed(time.Now().UnixNano())
//	return agents[rand.Intn(len(agents)-1)]
//
//}
//func Agents() []string {
//	return []string{
//		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.70",
//		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
//		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
//		"Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0",
//		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.62",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 Edg/92.0.902.67",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.55",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
//		"Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
//		"Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:90.0) Gecko/20100101 Firefox/90.0",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
//		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
//	}
//}
//
//func (i *Icp) CheckStatusName(domain string) (*request.DomainInfo, error) {
//	client := http.Client{}
//	resp, err := client.Get(fmt.Sprintf("%s/n5fms8aktb/%s", global.GVA_CONFIG.Robot.IcpHost, domain))
//	if err != nil {
//		global.GVA_LOG.Error("请求服务期出错", zap.Error(err))
//		return nil, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		global.GVA_LOG.Error("读取响应流出错", zap.Error(err))
//
//		return nil, err
//	}
//	var m DResult
//	err = json.Unmarshal(body, &m)
//	if err != nil {
//		global.GVA_LOG.Error("解析出错", zap.Error(err))
//		return nil, err
//	}
//	if resp.StatusCode != http.StatusOK {
//		global.GVA_LOG.Error("解析出错", zap.Int("code", http.StatusOK))
//		return nil, err
//	}
//	if m.Code == 0 && len(m.Info) > 0 {
//		info := m.Info[0]
//		return &info, nil
//	} else {
//		if m.Code == 0 {
//			return &request.DomainInfo{}, nil
//		} else {
//			return nil, errors.New(m.Msg)
//		}
//	}
//}
//
//type DResult struct {
//	Msg  string               `json:"msg"`
//	Info []request.DomainInfo `json:"info"`
//	Code int                  `json:"code"`
//}
