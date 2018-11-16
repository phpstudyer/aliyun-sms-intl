package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	//ServeURL .
	ServeURL = "http://dysmsapi.ap-southeast-1.aliyuncs.com"
	// RegionID .
	RegionID = "ap-southeast-1"
)

//Sms .
type Sms struct {
	accessKeyID   string
	accessSercret string
	dict          map[string]string
}

//smsParam .
type smsParam struct {
	PhoneNumbers     string `json:"PhoneNumbers"`
	ContentCode      string `json:"ContentCode"`
	ContentParam     string `json:"ContentParam"`
	ExternalID       string `json:"ExternalId"`
	RegionID         string `json:"RegionId"`
	AccessKeyID      string `json:"AccessKeyId"`
	Format           string `json:"Format"`
	SignatureMethod  string `json:"SignatureMethod"`
	SignatureVersion string `json:"SignatureVersion"`
	SignatureNonce   string `json:"SignatureNonce"`
	Timestamp        string `json:"Timestamp"`
	Action           string `json:"Action"`
	Version          string `json:"Version"`
	SignName         string `json:"SignName"`
	Signature        string `json:"Signature"`
}

// smsResponse .
type smsResponse struct {
	RequestID     string `json:"RequestId"`
	ResultCode    string `json:"ResultCode"`
	ResultMessage string `json:"ResultMessage"`
	BizID         string `json:"BizId"`
}

//New .
func New(accessKeyID, accessSercret string) *Sms {
	return &Sms{
		accessKeyID,
		accessSercret,
		map[string]string{
			"OK":                              "请求成功",
			"isp.RAM_PERMISSION_DENY":         "RAM权限DENY",
			"isv.OUT_OF_SERVICE":              "业务停机",
			"isv.PRODUCT_UN_SUBSCRIPT":        "未开通云通信产品的阿里云客户",
			"isv.PRODUCT_UNSUBSCRIBE":         "产品未开通",
			"isv.ACCOUNT_NOT_EXISTS":          "账户不存在",
			"isv.ACCOUNT_ABNORMAL":            "账户异常",
			"isv.SMS_TEMPLATE_ILLEGAL":        "短信模板不合法",
			"isv.SMS_SIGNATURE_ILLEGAL":       "短信签名不合法",
			"isv.INVALID_PARAMETERS":          "参数异常",
			"isp.SYSTEM_ERROR":                "系统错误",
			"isv.MOBILE_NUMBER_ILLEGAL":       "非法手机号",
			"isv.MOBILE_COUNT_OVER_LIMIT":     "手机号码数量超过限制",
			"isv.TEMPLATE_MISSING_PARAMETERS": "模板缺少变量",
			"isv.BUSINESS_LIMIT_CONTROL":      "业务限流",
			"isv.INVALID_JSON_PARAM":          "JSON参数不合法，只接受字符串值",
			"isv.BLACK_KEY_CONTROL_LIMIT":     "黑名单管控",
			"isv.PARAM_LENGTH_LIMIT":          "参数超出长度限制",
			"isv.PARAM_NOT_SUPPORT_URL":       "不支持URL",
			"isv.AMOUNT_NOT_ENOUGH":           "账户余额不足",
		},
	}
}

//Send .
func (sms *Sms) Send(PhoneNumbers, ContentCode, ContentParam string) (*smsResponse, error) {
	sr := new(smsResponse)
	param := &smsParam{
		PhoneNumbers: PhoneNumbers,
		ContentCode:  ContentCode,
		ContentParam: ContentParam,
		SignName:     "",
		// ExternalID       string `json:"ExternalId"`
		RegionID:         "ap-southeast-1",
		AccessKeyID:      sms.accessKeyID,
		Format:           "JSON",
		SignatureMethod:  "HMAC-SHA1",
		SignatureVersion: "1.0",
		SignatureNonce:   time.Now().UTC().Format(time.RFC3339),
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		Action:           "SendSms",
		Version:          "2018-05-01",
	}
	paramStr, err := sms.urlEncode(param)
	if err != nil {
		return sr, err
	}

	param.Signature = url.QueryEscape(sms.sign(sms.accessSercret+"&", "GET&%2F&"+url.QueryEscape(sms.specialURLEncode(paramStr))))

	resp, err := http.Get(fmt.Sprintf("%s/?Signature=%s&%s", ServeURL, param.Signature, paramStr))
	if err != nil {
		return sr, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(sr)
	return sr, err
}

func (sms *Sms) sign(accessSercret string, stringToSign string) string {
	key := []byte(accessSercret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (sms *Sms) urlEncode(src interface{}) (string, error) {
	dest := make(map[string]string)
	srcSlice, err := json.Marshal(src)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(srcSlice, &dest); err != nil {
		return "", err
	}
	value := url.Values{}
	for k, v := range dest {
		if k == "Signature" {
			continue
		}
		value.Set(sms.specialURLEncode(k), sms.specialURLEncode(v))
	}
	str := value.Encode()
	return str, nil
}

//specialUrlEncode .
func (sms *Sms) specialURLEncode(value string) string {
	// value = url.QueryEscape(value)
	reg := regexp.MustCompile(`\+`)
	value = reg.ReplaceAllString(value, "%20")
	reg = regexp.MustCompile(`\*`)
	value = reg.ReplaceAllString(value, "%2A")
	reg = regexp.MustCompile(`%7E`)
	value = reg.ReplaceAllString(value, "~")
	// reg = regexp.MustCompile(`=$`)
	// value = reg.ReplaceAllString(value, "")
	return value
}
