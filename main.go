package main

import(
	"./logs"
	"net/http"
)

func main(){
	logs.Println("test_syslog")
	logs.Log(logs.Err,"Testting record error")
	logs.Log(logs.Warn, "Testing record warning")
	logs.Log(logs.Q_err,"Testting record error quietly")
	logs.Log(logs.Q_warn, "Testing record warning quietly")

	mux := http.NewServeMux()
	mux.HandleFunc("/log/getlogtext", logs.SendLogText)
	mux.HandleFunc("/log/getlogpage", logs.SendLogList)
	server := &http.Server{
		Addr : 			"localhost:8090",
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}
	logs.Println("begin to listen!")
	err := server.ListenAndServe()
	if err!=nil {
		logs.Log(logs.Err, err)
	}

}