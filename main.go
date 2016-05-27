package main

import (
	"os"
	//"path/filepath"
	"fmt"
	"bufio"
	"strings"

	"io"
	"net/url"
	"encoding/json"
	"net/http"
	"bytes"
	"net"
	"time"
	"crypto/tls"
	"io/ioutil"
	"axon"
	"path/filepath"
)

var (
	mob []string
	num int64
err error
)

func main() {
	mob=append(mob,"./INFO-20160401.log")
	mob=append(mob,"./INFO-20160402.log")
	mob=append(mob,"./INFO-20160403.log")
	mob=append(mob,"./INFO-20160404.log")
	mob=append(mob,"./INFO-20160405.log")
	mob=append(mob,"./INFO-20160406.log")
	mob=append(mob,"./INFO-20160407.log")
	mob=append(mob,"./INFO-20160408.log")
	mob=append(mob,"./INFO-20160409.log")
	mob=append(mob,"./INFO-20160410.log")
	mob=append(mob,"./INFO-20160411.log")
	mob=append(mob,"./INFO-20160412.log")
	mob=append(mob,"./INFO-20160413.log")
	mob=append(mob,"./INFO-20160414.log")
	mob=append(mob,"./INFO-20160415.log")
	mob=append(mob,"./INFO-20160416.log")
	mob=append(mob,"./INFO-20160417.log")
	//mob=append(mob,"./INFO-20160418.log")
	mob=append(mob,"./INFO-20160426.log")
	mob=append(mob,"./INFO-20160427.log")
	mob=append(mob,"./INFO-20160428.log")
	mob=append(mob,"./INFO-20160429.log")
	mob=append(mob,"./INFO-20160430.log")
	mob=append(mob,"./")
	mob=append(mob,"./")
	for _, value := range mob {
		num=0
		readFileMobs(value)
		fmt.Printf("%s send http %d number \n",value,num)
		err = filepath.Walk(value, func(path string, f os.FileInfo, err error) error {
			if ( f == nil ) {return err}
			if f.IsDir() {return nil}


			return nil
		})
		if err != nil {
			fmt.Printf("filepath.Walk() returned %v\n", err)
		}
	}

	writeFileMobs();
}

func getMyUserJumpUrl(mob, price, province, channel string) (string,error) {
	dataValue:=url.Values{}
	dataValue.Add("phone",mob)
	dataValue.Add("price",price)
	dataValue.Add("province",province)
	dataValue.Add("channel",channel)
	dataValueJson,err:=json.Marshal(dataValue)
	if err != nil {
		fmt.Printf("getMyUserJumpUrl data Marshal error! error: %s \n",err.Error())
		return "",err
	}
	resp,err:=http.NewRequest("POST","https://api.pinzhi.xin/sb/addflow",bytes.NewReader(dataValueJson))
	if err != nil {
		fmt.Printf("getMyUserJumpUrl NewRequest error! error: %s sendUrl: %s dataValue: %s \n",err.Error(),"https://api.pinzhi.xin/sb/addflow",string(dataValueJson))
		return "",err
	}
	resp.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	httpClient.Transport=&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	respone, err := httpClient.Do(resp)
	if err != nil {
		fmt.Printf("getMyUserJumpUrl send http post error! err: %s sendUrl: %s dataValue: %s \n",err.Error(),"https://api.pinzhi.xin/sb/addflow",string(dataValueJson))
		return "",err
	}
	requestData,err:=ioutil.ReadAll(respone.Body)
	defer respone.Body.Close()
	if err != nil {
		fmt.Printf("getMyUserJumpUrl send http read request data error! error: %s sendUrl: %s dataValue: %s \n",err.Error(),"https://api.pinzhi.xin/sb/addflow",string(dataValueJson))
		return "",err
	}
	//var model PointXinBaoApiData
	//err=json.Unmarshal(requestData,&model)
	//if err != nil {
	//	fmt.Printf("getMyUserJumpUrl request data unmarshal error! err: %s sendUrl: %s dataValue: %v \n",err.Error(),"https://api.pinzhi.xin/sb/addflow",string(dataValueJson))
	//	return nil,err
	//}
	return string(requestData),nil
}

func writeFileMobs() {
	fs,err:=os.OpenFile("./fileListMob.txt",os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("文件读取失败")
		return
	}
	//fmt.Println(mob)
	_,err=fs.WriteString(strings.Join(mob,"\n"))
	if err != nil {
		fmt.Println("写入失败,err: %s ",err.Error())
		return
	}
	fs.WriteString(strings.Join(mob,"\n"))
	fs.Close()
}

func readFileMobs(filePath string) {
	fs,err:=os.Open(filePath)

	if err != nil {
		fmt.Println("文件读取失败!")
		return
	}
	buf := bufio.NewReader(fs)

	for {

		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.Index(line, "nationwideDiscntInfo") > 0 {
			num++
		}
		ff:=axon.AxonEncrypt(line)
		ff=fmt.Sprintf("delete from userblacklist where Mob='%s';",ff)
			//mob=append(mob,ff)

		//if len(ff)==11 {
		//	//fmt.Println(axon.AxonEncrypt(ff))
		//}
		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
}