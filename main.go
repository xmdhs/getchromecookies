package getchromecookies

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Chrome struct {
	Path    string
	DataDir string
	Web     string
}

func NewChrome() *Chrome {
	return &Chrome{
		Path: locateChrome(),
	}
}

func (c *Chrome) GetCookie() (s []CookieResultCooky, err error) {
	cxt := context.TODO()
	cxt, cancel := context.WithCancel(cxt)
	defer cancel()
	port := GetProt()
	_, err = c.runChrome(cxt, port)
	if err != nil {
		return nil, fmt.Errorf("GetCookie: %w", err)
	}
	ws, err := getDebugWs(port)
	if err != nil {
		return nil, fmt.Errorf("GetCookie: %w", err)
	}
	var cookie cookie
	u, err := url.Parse(c.Web)
	if err != nil {
		return nil, fmt.Errorf("GetCookie: %w", err)
	}
F:
	for {
		b, err := getCookie(ws)
		if err != nil {
			return nil, fmt.Errorf("GetCookie: %w", err)
		}
		err = json.Unmarshal(b, &cookie)
		if err != nil {
			return nil, fmt.Errorf("GetCookie: %w", err)
		}
		if len(cookie.Result.Cookies) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}
		for _, v := range cookie.Result.Cookies {
			if strings.HasSuffix(u.Hostname(), v.Domain) {
				if v.Value == "" && v.Size != 0 {
					time.Sleep(2 * time.Second)
					continue F
				}
			}
		}
		return cookie.Result.Cookies, nil
	}
}

func getCookie(ws string) ([]byte, error) {
	c, _, err := websocket.DefaultDialer.Dial(ws, nil)
	if err != nil {
		return nil, fmt.Errorf("getCookie: %w", err)
	}
	defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, []byte(`{"id": 1, "method": "Network.getCookies"}`))
	if err != nil {
		return nil, fmt.Errorf("getCookie: %w", err)
	}
	_, b, err := c.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("getCookie: %w", err)
	}
	return b, nil
}

func (c *Chrome) runChrome(cxt context.Context, port int) (*exec.Cmd, error) {
	cmd := exec.CommandContext(cxt, c.Path, "--user-data-dir="+c.DataDir, c.Web, "--remote-debugging-port="+strconv.Itoa(port))
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("runChrome: %w", err)
	}
	return cmd, nil
}

func getDebugWs(port int) (string, error) {
	url := "http://127.0.0.1:" + strconv.Itoa(port) + "/json"
	rep, err := http.Get(url)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("getDebugWs: %w", err)
	}
	b, err := io.ReadAll(rep.Body)
	w := []wsjson{}
	err = json.Unmarshal(b, &w)
	if err != nil {
		return "", fmt.Errorf("getDebugWs: %w", err)
	}
	return w[0].Ws, nil
}

type wsjson struct {
	Ws string `json:"webSocketDebuggerUrl"`
}

func GetProt() int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	list := strings.Split(l.Addr().String(), ":")
	i, err := strconv.ParseInt(list[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return int(i)
}
