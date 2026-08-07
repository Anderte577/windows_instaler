package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type Config struct {
	ID       string `json:"id"`
	Server   string `json:"server"`
	Interval int    `json:"interval"`
}

type Status struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	System   string `json:"system"`
	Time     int64  `json:"time"`
}

func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func loadConfig() Config {

	data, err := os.ReadFile("config.json")

	if err != nil {
		panic("Brak config.json")
	}

	var c Config

	json.Unmarshal(data, &c)

	return c
}


func sendStatus(c Config){

	host,_:=os.Hostname()

	status:=Status{
		ID:c.ID,
		Hostname:host,
		IP:getIP(),
		System:runtime.GOOS,
		Time:time.Now().Unix(),
	}


	data,_:=json.Marshal(status)


	req,_:=http.NewRequest(
		"POST",
		c.Server,
		bytes.NewBuffer(data),
	)


	req.Header.Set(
		"Content-Type",
		"application/json",
	)


	client:=http.Client{
		Timeout:5*time.Second,
	}


	resp,err:=client.Do(req)

	if err==nil{
		io.Copy(io.Discard,resp.Body)
		resp.Body.Close()
	}

}


func main(){

	config:=loadConfig()


	for{

		sendStatus(config)

		time.Sleep(
			time.Duration(config.Interval)*time.Second,
		)
	}

}
