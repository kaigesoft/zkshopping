package tools

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ILog Log 接口
type ILog interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{}) error
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Errorf(format string, params ...interface{}) error
}

var escaper = strings.NewReplacer("\n", "\\n", "\r", "", "|", "\\|")
var localIp string
var program string
func WriteApiLog(reqID, level, message string) error {
	conn, err := getConn33825()

	if err != nil {
		return seelog.Error("DialUDP ", UDPAddr2, " failed:", err)
	}

	_, file, line, _ := runtime.Caller(2)
	file = file[strings.LastIndex(file, "/")+1:]

	curTime := time.Now().Format("2006-01-02 15:04:05")

	messages := make([]string, 0, len(message)/4000+1)
	if len(message) < 4000 {
		messages = append(messages, escaper.Replace(message))
	} else {
		for i := 0; i < len(message); i += 4000 {
			if i+4000 > len(message) {
				messages = append(messages, escaper.Replace(message[i:]))
			} else {

				messages = append(messages, escaper.Replace(message[i:i+4000]))
			}
		}
	}

	segment := len(messages)
	msgPrefix := strings.Join([]string{"DataMoreApiLog", curTime, reqID, gin.Mode(), localIp, program, level, file + ":" + strconv.Itoa(line)}, "|")

	for i, msg := range messages {
		segmentDesc := " "
		if segment > 1 {
			segmentDesc = fmt.Sprintf(" (%d/%d) ", i+1, segment)
		}
		_, err = fmt.Fprintln(conn, msgPrefix, segmentDesc, msg)

		if err != nil {
			conn33825 = nil
			fmt.Println("udp log failed:", err)
			return err
		}
	}

	return nil
}

var UDPAddr2 *net.UDPAddr
var conn33825 *net.UDPConn

func getConn33825() (*net.UDPConn, error) {
	if conn33825 == nil {
		var err error
		conn33825, err = net.DialUDP("udp", nil, UDPAddr2)
		return conn33825, err
	}
	return conn33825, nil
}
