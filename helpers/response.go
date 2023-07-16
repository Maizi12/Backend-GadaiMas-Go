package helpers

import(
	"net/http"
	"encoding/json"
	"fmt"
)

type ResponseWithData struct{
	//Status string `json:"status"`
	//Message string `json:"message"`
	Data any `json:"data"`
}
type ResponseWithToken struct{
	Data any `json:"data"`
	Token any `json:"token"`
}
type ResponseWithoutData struct{
	Status string `json:"-"`
	Message string `json:"message"`
}

func Response(w http.ResponseWriter, code int, message string,payload interface{},payload2 interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var response any
	status:="success"
	fmt.Println(status)
	if code>=400{
		status="failed"
	}
	if payload!=nil&&payload2!=nil{
		response =&ResponseWithToken{
		Data:payload,
		Token:payload2,
		}
	}
	if payload!=nil&&payload2==nil{
		response =&ResponseWithData{
			//Status:status,
			//Message:message,
			Data:payload,
		}
	}
	if payload==nil && payload2==nil{
		response=&ResponseWithoutData{
			Message:message,
		}
	}
	res,_:=json.Marshal(response)
	w.Write(res)
}