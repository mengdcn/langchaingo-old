package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/tgis/thenextleg"
	"html/template"
	"os"
	"reflect"
	"strconv"
	"time"
)

func main() {
	//	var ssp thenextleg.MessageResponse
	//	s := `{"progress":100,
	//"response":{"accountId":"2HNhG5DksW4OLPtjG3ZP","createdAt":"2023-10-24T11:34:09.667Z","originatingMessageId":"Hk1aJbDvdV2MrW9qolRN","buttons":["U1","U2","U3","U4","🔄","V1","V2","V3","V4"],"imageUrl":"https://cdn.discordapp.com/attachments/1152602889924124682/1166338842316509254/freedomlink._Carp_in_the_Lotus_Pond_f76c6db0-a743-49f3-9347-0593a3091b23.png?ex=654a20b1&is=6537abb1&hm=efda793c5ba9188c185725ca7cb8a52d972a73f0eada15a476f6e8f266b3ee5a&","imageUrls":["https://cdn.midjourney.com/f76c6db0-a743-49f3-9347-0593a3091b23/0_0.png","https://cdn.midjourney.com/f76c6db0-a743-49f3-9347-0593a3091b23/0_1.png","https://cdn.midjourney.com/f76c6db0-a743-49f3-9347-0593a3091b23/0_2.png","https://cdn.midjourney.com/f76c6db0-a743-49f3-9347-0593a3091b23/0_3.png"],"responseAt":"2023-10-24T11:34:10.190Z","description":"","type":"imagine","content":"Carp in the Lotus Pond","buttonMessageId":"INNNWP6Nkfo2cixjkWx5"}}`
	//	err := json.Unmarshal([]byte(s), &ssp)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println(ssp)
	//	return

	token := os.Getenv("mid_token")
	// 生成图片任务和查询进度  ----start
	c, err := thenextleg.New(thenextleg.WithAuthToken(token))
	if err != nil {
		fmt.Println(err.Error())
	}
	payload := &thenextleg.ImagineRequest{
		Msg: "Carp in the Lotus Pond",
		//Ref:             "",
		//WebhookOverride: "",
		//IgnorePrefilter: "",
	}
	resp, err := c.Imagine(context.Background(), payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", resp)

	// get message

	msgId := resp.MessageId

	if msgId == "" {
		fmt.Println("message id is empty")
		return
	}
	var respp *thenextleg.MessageResponse
	for i := 0; i < 30; i++ {
		respp, err = c.Message(context.Background(), msgId)
		if err != nil {
			fmt.Println(err)
			break
		}
		a, _ := json.Marshal(respp)

		fmt.Printf("%#v\n", string(a))
		str, err := ToStringE(respp.Progress)
		if err != nil {
			fmt.Println(err)
			break
		}
		// 完成
		if str == "100" {
			break
		}
		// 失败
		if str == "incomplete" {
			break
		}
		time.Sleep(time.Second * 2)
	}

	if respp.Response.ButtonMessageId == "" {
		fmt.Println("fail")
		return
	}
	// 生成图片任务和查询进度  -------end

	// 对生成的图片微调任务，和微调任务进度查询   -----start
	// buttons  INNNWP6Nkfo2cixjkWx5
	btnMessageid := respp.Response.ButtonMessageId
	//btnMessageid := "INNNWP6Nkfo2cixjkWx5"

	butPayload := &thenextleg.ButtonRequest{
		ButtonMessageId: btnMessageid,
		Button:          "U1",
	}
	btnResp, err := c.Button(context.Background(), butPayload)

	bgtMsgId := btnResp.MessageId
	var resppp *thenextleg.MessageResponse
	for i := 0; i < 30; i++ {
		resppp, err = c.Message(context.Background(), bgtMsgId)
		if err != nil {
			fmt.Println(err)
			break
		}
		a, _ := json.Marshal(resppp)

		fmt.Printf("%#v\n", string(a))
		str, err := ToStringE(resppp.Progress)
		if err != nil {
			fmt.Println(err)
			break
		}
		// 完成
		if str == "100" {
			break
		}
		// 失败
		if str == "incomplete" {
			break
		}
		time.Sleep(time.Second * 2)
	}
	// 对生成的图片微调任务，和微调任务进度查询   -----end
}

func ToStringE(i any) (string, error) {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

func indirectToStringerOrError(a any) any {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
