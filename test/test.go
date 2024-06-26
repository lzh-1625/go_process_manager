package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	// 启动 Minecraft 服务器
	cmd := exec.Command("java", "-jar", "launcher-airplane.jar")

	// 创建伪终端
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println("Start pty:", err)
		return
	}
	defer func() { _ = ptmx.Close() }() // 最后关闭伪终端

	// 创建通道用于读取子程序的输出
	outputChan := make(chan string)

	go func() {
		defer close(outputChan)
		readOutput(ptmx, outputChan)
	}()

	// 将子程序的输出发送到 WebSocket 客户端
	go func() {
		for output := range outputChan {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(output)); err != nil {
				log.Println("WriteMessage:", err)
				break
			}
		}
	}()

	// 从 WebSocket 客户端读取消息并发送到子程序的标准输入
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage:", err)
			break
		}
		// 检查是否是发送 Tab 键的命令
		if string(message) == "SEND_TAB" {
			_, err := ptmx.Write([]byte{9}) // Tab 键的 ASCII 码是 9
			if err != nil {
				log.Println("Write to stdin:", err)
			}
		} else if string(message) == "READ_INPUT" {
			// 读取伪终端当前输入管道内已经存在的内容
			// 注意：伪终端没有单独的输入缓冲区，输入会立即被处理
			// 这里假设你想获取当前输出内容
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Currently no direct way to fetch unsent input. Consider monitoring the terminal buffer.")); err != nil {
				log.Println("WriteMessage:", err)
			}
		} else {
			_, err := ptmx.Write(append(message, '\n')) // 确保命令后有换行符
			if err != nil {
				log.Println("Write to stdin:", err)
			}
		}
	}

	// 等待子程序结束
	if err := cmd.Wait(); err != nil {
		log.Println("Wait:", err)
		return
	}
}

func readOutput(reader io.Reader, outputChan chan<- string) {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			outputChan <- string(buf[:n])
		}
		if err != nil {
			if err != io.EOF {
				log.Println("Read:", err)
			}
			break
		}
	}
}
