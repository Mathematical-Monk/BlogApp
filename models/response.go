package models

type HttpResponse struct {
	Msg string `json:"msg"`
}


func CreateResStruct(msg string) HttpResponse {
	return HttpResponse{
		Msg : msg,
	}
	
}
