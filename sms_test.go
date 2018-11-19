package sms

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {

	resp, err := New("your_accessid", "your_accesskey").Send("1828761***", "SMS_10****", "{\"code\":\"1234\"}")
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	if resp.ResultCode !="OK"{
		fmt.Println("err:", resp.ResultMessage)
		return
	}
	fmt.Println(resp.RequestID)
}
