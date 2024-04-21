package services

import (
	"fmt"
	"github.com/captain686/go-warp2wireguard/util"
	"github.com/charmbracelet/log"
	"net"
	"sort"
	"sync"
	"time"
)

var alreadyList []string

func udpConnection(udpServer string) (bool, error) {
	conn, err := net.DialTimeout("udp", udpServer, time.Second*5)
	if err != nil {
		log.Error("UDP Connection Test Fail:", err)
		return false, err
	}
	log.Info("UDP Connection Test Passed")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	return true, nil
}

//func generateRandomHex() (*string, error) {
//	randomBytes := make([]byte, 16)
//
//	// 从加密随机源中读取随机字节
//	_, err := rand.Read(randomBytes)
//	if err != nil {
//		return nil, err
//	}
//
//	// 将随机字节转换为16进制字符串
//	randomHex := hex.EncodeToString(randomBytes)
//
//	return &randomHex, nil
//}

func handshake(ip string, port int, ch chan bool, res chan util.PassTest, wg *sync.WaitGroup) {
	serverAddr := fmt.Sprintf("%s:%d", ip, port)
	defer func() {
		wg.Done()
		<-ch
	}()
	response := util.PassTest{
		Address: serverAddr,
	}
	conn, err := net.DialTimeout("udp", serverAddr, time.Second*2)
	if err != nil {
		log.Error("Error connecting to server:", err)
		return
	}
	//connRtt := time.Now().Sub(coonStartTime).Milliseconds()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	go time.AfterFunc(time.Second*2, func() {
		res <- response
		err := conn.Close()
		if err != nil {
			return
		}
	})
	startTime := time.Now()
	//packet, err := generateRandomHex()
	//if err != nil {
	//	return
	//}
	handshakePacket, err := util.HandshakePacket()
	if err != nil {
		return
	}
	// 发送UDP请求
	_, err = conn.Write(handshakePacket)
	if err != nil {
		log.Error("Error sending UDP request:", err)
		return
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		//log.Error("Error receiving UDP response:", err)
		return
	}

	endTime := time.Now()
	rtt := endTime.Sub(startTime).Milliseconds()

	log.Info(fmt.Sprintf("Server %s -> RTT :%d ms", serverAddr, rtt))

	//fmt.Printf("Response (hex): %x\n", buf[:n])
	response.ReturnTime = &rtt
	res <- response
	alreadyList = append(alreadyList, ip)
}

func already(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func WarpSpeedTest() (*sync.Map, error) {
	var (
		wg      sync.WaitGroup
		coonMap sync.Map
	)
	udpPass, err := udpConnection("8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	if !udpPass {
		return nil, err
	}

	ch := make(chan bool, 600)
	res := make(chan util.PassTest, 600)
	go func() {
		for result := range res {
			if result.ReturnTime != nil {
				coonMap.Store(result.Address, *result.ReturnTime)
			}
			//_ = fmt.Sprintf(result.Address)
		}
	}()
	for _, cidrIp := range util.ServerIPs {
		for _, port := range util.CommonWarpPorts {
			for serverIp := range util.CheckCidrIPs(cidrIp) {
				if already(alreadyList, *serverIp) {
					continue
				}

				ch <- true
				wg.Add(1)
				go handshake(*serverIp, port, ch, res, &wg)
			}
		}
	}
	wg.Wait()
	close(ch)
	time.Sleep(time.Second * 5)
	close(res)
	//services.Handshake(udpServer)
	alreadyList = nil
	return &coonMap, nil
}

func SpeedSort(coonMap *sync.Map) []util.Speed {

	var keyValueList []util.Speed
	coonMap.Range(func(k, v interface{}) bool {
		keyValueList = append(keyValueList, util.Speed{Server: k.(string), TimeOut: v.(int64)})
		return true
	})

	// 自定义排序函数
	sort.Slice(keyValueList, func(i, j int) bool {
		return keyValueList[i].TimeOut < keyValueList[j].TimeOut
	})

	// 打印排序后的结果
	//for _, kv := range keyValueList {
	//	fmt.Printf("Key: %s, Value: %d\n", kv.Server, kv.TimeOut)
	//}
	coonMap = nil
	return keyValueList
}
