package util

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func Connect(host string, port int, id_rsa string) (*ssh.Session, chan<- string, <-chan string, error) {
	// 读取配置文件

	// 读取私钥
	privateKey, err := ioutil.ReadFile(id_rsa)
	if err != nil {
		log.Errorf("get private key error : %s", err.Error())
		panic(err)
	}

	// 连接配置
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Errorf("parse private key error : %s", err.Error())
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 开始连接
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		println("login err :" + err.Error())
	}

	session, err := client.NewSession()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("vt100", 80, 80, modes); err != nil {
		log.Fatal(err)
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	e, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}

	in, out := MuxShell(w, r, e)
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	return session, in, out, nil
}

func MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 3)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			t += n
			result := string(buf[:t])
			if strings.Contains(result, "[sudo] password") ||
				strings.Contains(result, "#") ||
				strings.Contains(result, "$") {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}
