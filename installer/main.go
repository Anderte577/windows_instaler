package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	ID string `json:"id"`
	Server string `json:"server"`
	Interval int `json:"interval"`
}


const server =
"http://192.168.0.8:5000/status"


func main(){

	var id string

	fmt.Print("Podaj ID komputera: ")

	fmt.Scanln(&id)


	folder:=
	"C:\\HKL-Agent"


	os.MkdirAll(
		folder,
		0755,
	)


	exe,_:=os.Executable()


	copyFile(
		filepath.Join(filepath.Dir(exe),"HKL-Agent.exe"),
		filepath.Join(folder,"HKL-Agent.exe"),
	)


	config:=Config{

		ID:id,
		Server:server,
		Interval:30,
	}


	data,_:=json.MarshalIndent(
		config,
		"",
		"  ",
	)


	os.WriteFile(
		filepath.Join(folder,"config.json"),
		data,
		0644,
	)


	exec.Command(
		"reg",
		"add",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`,
		"/v",
		"HKL-Agent",
		"/t",
		"REG_SZ",
		"/d",
		folder+"\\HKL-Agent.exe",
		"/f",
	).Run()


	exec.Command(
		folder+"\\HKL-Agent.exe",
	).Start()


	fmt.Println("Gotowe")

}


func copyFile(src,dst string){

	data,err:=os.ReadFile(src)

	if err!=nil{
		return
	}

	os.WriteFile(
		dst,
		data,
		0755,
	)
}
