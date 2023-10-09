package checker

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"golang.org/x/exp/rand"
	"image"
	"io"
	"net/http"
	"strings"
	"time"
)

var IcpNotForRecord = errors.New("域名未备案")

type IcpClient struct {
	core   *http.Client
	token  string
	agent  string
	cookie string
}

func NewIcpClient() *IcpClient {
	agent := RandAgent()
	return &IcpClient{agent: agent, core: http.DefaultClient}
}

func (i *IcpClient) GetCookies() error {

	req, err := http.NewRequest(http.MethodGet, "https://beian.miit.gov.cn/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", i.agent)
	res, err := i.core.Do(req)
	if err != nil {
		return err
	}
	fmt.Print(res.Cookies())

	for _, cookie := range res.Cookies() {
		if cookie.Name == "__jsluid_s" {
			i.cookie = cookie.Value
			return nil
		}
	}
	return errors.New("无法获取到cookie")
}

//	func (i *Icp) Query(domain string) (*request.DomainInfo, error) {
//		i.getUserAgent()
//		if err := i.auth(); err != nil {
//			return nil, errors.New("获取授权失败:" + err.Error())
//		}
//		return i.query(domain)
//	}
//
//	func (i *Icp) query(domain string) (*request.DomainInfo, error) {
//		queryRequest, _ := json.Marshal(&request.QueryRequest{
//			UnitName: domain,
//		})
//
//		result := &request.IcpResponse{Params: &request.QueryParams{}}
//		err := i.post("icpAbbreviateInfo/queryByCondition", bytes.NewReader(queryRequest), "application/json;charset=UTF-8", i.token, result)
//		if err != nil {
//			return nil, err
//		}
//
//		if !result.Success {
//			return nil, fmt.Errorf("查询：%s", result.Msg)
//		}
//
//		queryParams := result.Params.(*request.QueryParams)
//		if len(queryParams.List) == 0 {
//			return &request.DomainInfo{}, nil
//		}
//		return queryParams.List[0], nil
//	}
func (i *IcpClient) GetToken() error {
	timestamp := time.Now().Unix()
	data := []byte(fmt.Sprintf("testtest%d", timestamp))
	has := md5.Sum(data)
	authBody := fmt.Sprintf("authKey=%s&timeStamp=%d", fmt.Sprintf("%x", has), timestamp)

	req, err := http.NewRequest(http.MethodPost, "https://hlwicpfwc.miit.gov.cn/icpproject_query/api/auth", bytes.NewReader([]byte(authBody)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://beian.miit.gov.cn/")
	req.Header.Set("Referer", "https://beian.miit.gov.cn/")
	req.Header.Set("User-Agent", i.agent)
	req.Header.Set("Cookie", fmt.Sprintf("__jsluid_s=%s", i.cookie))

	resp, err := i.core.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("获取token的请求出错了")
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var tokenRet TokenRet
	err = json.Unmarshal(body, &tokenRet)
	if err != nil {
		return err
	}
	i.token = tokenRet.Params.Bussiness
	return nil
}

type TokenRet struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Params struct {
		Bussiness string `json:"bussiness"`
		Expire    int    `json:"expire"`
		Refresh   string `json:"refresh"`
	} `json:"params"`
	Success bool `json:"success"`
}

func (i *IcpClient) ImageVerify() error {
	req, err := http.NewRequest(http.MethodPost, "https://hlwicpfwc.miit.gov.cn/icpproject_query/api/image/getCheckImage", nil)
	if err != nil {
		return err
	}
	//req.Header.Set("Content-Type", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://beian.miit.gov.cn/")
	req.Header.Set("Referer", "https://beian.miit.gov.cn/")
	req.Header.Set("User-Agent", i.agent)
	req.Header.Set("Cookie", fmt.Sprintf("__jsluid_s=%s", i.cookie))
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("token", i.token)
	req.Header.Set("Content-Length", "0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("获取token的请求出错了")
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var imageRet ImageRet
	err = json.Unmarshal(body, &imageRet)
	if err != nil {
		return err
	}
	if imageRet.Code != 200 {
		return errors.New(imageRet.Msg)

	}

	println(imageRet.Params.BigImage)
	println(imageRet.Params.SmallImage)
	println(imageRet.Params.Height)
	println(imageRet.Params.Uuid)
	return nil

}

// 图片信息
type ImageRet struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Params struct {
		BigImage   string `json:"bigImage"`
		Height     string `json:"height"`
		SmallImage string `json:"smallImage"`
		Uuid       string `json:"uuid"`
	} `json:"params"`
	Success bool `json:"success"`
}

/*

{
    "code": 200,
    "msg": "操作成功",
    "params": {
        "bigImage": "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",
        "height": "78",
        "smallImage": "/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAA0JCgsKCA0LCgsODg0PEyAVExISEyccHhcgLikxMC4pLSwzOko+MzZGNywtQFdBRkxOUlNSMj5aYVpQYEpRUk//2wBDAQ4ODhMREyYVFSZPNS01T09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT0//wAARCABCAEIDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwCCN1ddyrt55FP2buaHkj42Kqk84X0q1GITtHmbWx0r2+fuee7X0K6wZp7QYXNXDF0+bd70hhyMVDnqIzRHk1ahXA2+9PNo2fl6VJFblRlutTKehothRHgUhTmpsYprbsfLHu96wbuUiLZRR+8/540UhmNNDLbsWi+4eT9aSBtzbW+8ealW5VFG0sfXd2NRT4YB+9dl7mLLEdzMhKBsL6/3fenW+oXMLFZ5PMUnIb1FVVkO0BqXI6N900+W4joLa6iuUwsm3npU5aMcbt2O9cwrMkoAO3uGrXs9SjbEVyuw9pPWuepTtqawlYv746UPHipAjYBQ5U9D60FDn5utc7Nr3Gb46KdsopAccTmlXbn5utMIwaVWwa35+xm43LAQOPl6UExkYb7w4pq/P8u38amDDG0tkVUarW5LgVs5bb+NSbmK5XoOKkaFCNydKhIwa6VUhJGbVmCysdwPLAZX/ZPrXQaVqsUtusVzzMOFPr71zx6U1WCcltvPBqKlNSWhUZcp3IRyMnrRXKi+vdoxd8dqK5fqxr7Yzx0pwUEe9IPLkXjoKaZBj5eg4rm9pctqxJgj71OHSoI7jPFPMg3AdyKpTJJRwc1Jv3cVBkj71IX5z/d5rRT0Jcbk5TioinNIJ1I3t1NQXDiLzCvR0H8zWiq2QchPsopU1DKL9KKftg5ChCT57fUfyNQN92T6D+VFFeUjeW46D78P0qeb77/UUUUylsWT/rh/uikBO8c9j/SiirQmU5/9Q/8Avf1NF0T5g5/gooq1sBkgnA5ooopgf//Z",
        "uuid": "13b22cf5-0162-4b62-a29a-983d902c48cb"
    },
    "success": true
}

*/

//	func (i *Icp) post(url string, data io.Reader, contentType string, token string, result interface{}) error {
//		postUrl := fmt.Sprintf("https://hlwicpfwc.miit.gov.cn/icpproject_query/api/%s", url)
//		queryReq, err := http.NewRequest(http.MethodPost, postUrl, data)
//		if err != nil {
//			return err
//		}
//		queryReq.Header.Set("Content-Type", contentType)
//		queryReq.Header.Set("Origin", "https://beian.miit.gov.cn/")
//		queryReq.Header.Set("Referer", "https://beian.miit.gov.cn/")
//		queryReq.Header.Set("token", token)
//		queryReq.Header.Set("User-Agent", i.agent)
//
//		//client := DefaultProxyClient()
//		client := http.DefaultClient
//		resp, err := client.Do(queryReq)
//
//		return GetHTTPResponse(resp, postUrl, err, result)
//	}
//
//	func (i *Icp) getUserAgent() {
//		i.agent = RandAgent()
//	}
//
//	func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
//		if err != nil {
//			return err
//		}
//		body, err := GetHTTPResponseOrg(resp, url, err)
//		if err == nil {
//			err = json.Unmarshal(body, &result)
//			if err != nil {
//				return err
//			}
//		}
//		return err
//	}
//
// func GetHTTPResponseOrg(resp *http.Response, url string, err error) ([]byte, error) {
//
//		defer resp.Body.Close()
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			return nil, err
//		}
//		// 300及以上状态码都算异常
//		if resp.StatusCode >= 300 {
//			errMsg := fmt.Sprintf("请求接口 %s 失败! 返回状态码: %d\n", url, resp.StatusCode)
//			err = fmt.Errorf(errMsg)
//		}
//		return body, err
//	}
func RandAgent() string {
	agents := Agents()
	rand.Seed(uint64(time.Now().UnixNano()))
	return agents[rand.Intn(len(agents)-1)]

}
func Agents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.70",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.62",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 Edg/92.0.902.67",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.55",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	}
}

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

func decode(b64img string) []byte {
	i := strings.IndexByte(b64img, ',')
	if i == -1 {
		return nil
	}
	b, err := base64.StdEncoding.DecodeString(b64img[i+1:])
	if err != nil {
		return nil
	}
	return b
}

func readBase64Image(b64Image string) (gocv.Mat, error) {
	origin, err := gocv.IMDecode(decode(b64Image), gocv.IMReadUnchanged)
	if err != nil {
		return gocv.Mat{}, err
	}
	return origin, nil
}

func preProcess(b64Image string) (alpha, processed gocv.Mat, err error) {
	origin, err := readBase64Image(b64Image)
	if err != nil {
		return gocv.Mat{}, gocv.Mat{}, err
	}
	//defer origin.Close()

	//resized := resize(origin, origin.Rows(), origin.Cols())
	grayed := gray(origin)
	//threshold := threshold(grayed)
	//defer resized.Close()
	//defer grayed.Close()
	//defer threshold.Close()

	//log.Debugf(origin.Cols(), origin.Rows(), resized.Cols(), resized.Rows())

	if origin.Channels() == 4 {
		return gocv.Split(origin)[3], grayed, nil
	}

	return gocv.Mat{}, grayed, nil
}

func resize(origin gocv.Mat, cols, rows int) gocv.Mat {
	resized := gocv.NewMatWithSize(cols, rows, origin.Type())
	gocv.Resize(origin, &resized, image.Pt(cols, rows), 0, 0, gocv.InterpolationNearestNeighbor)
	return resized
}

func gray(origin gocv.Mat) gocv.Mat {
	grayed := gocv.NewMat()
	gocv.CvtColor(origin, &grayed, gocv.ColorBGRToGray)
	return grayed
}

func Match(bg, block, mask gocv.Mat) image.Point {
	result := gocv.NewMatWithSize(
		bg.Rows()-block.Rows()+1,
		bg.Cols()-block.Cols()+1,
		gocv.MatTypeCV32FC1)
	defer result.Close()

	gocv.MatchTemplate(bg, block, &result, gocv.TmSqdiff, mask)
	gocv.Normalize(result, &result, 0, 1, gocv.NormMinMax)

	_, _, _, maxLoc := gocv.MinMaxLoc(result)

	return maxLoc
}
