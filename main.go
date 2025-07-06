package main

import (
	"flag"
	"fmt"
	"net/http"
)

var mainPort int

func main() {
	flag.IntVar(&mainPort, "n", 8080, "指定主端口")
	flag.Parse()

	// 静态文件模式（index.html）
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	// go template 模式
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.New("portcheck").Parse(portCheckHTML))
	// 	tmpl.Execute(w, nil)
	// })

	// 处理监听请求
	http.HandleFunc("/listen", listenHandler)

	addr := fmt.Sprintf(":%d", mainPort)
	fmt.Printf("主服务启动: http://localhost%s\n", addr)
	http.ListenAndServe(addr, nil)
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	port := r.FormValue("port")

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte("pong from port " + port))
		})

		server := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}
		fmt.Println("启动监听端口: ", port)
		_ = server.ListenAndServe()
	}()

	w.Write([]byte("监听已启动: " + port))
}

const portCheckHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Port Check</title>
  <style>
    body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
    h1 { color: #333; border-bottom: 1px solid #eee; padding-bottom: 10px; }
    form { margin: 15px 0; padding: 15px; background: #f5f5f5; border-radius: 5px; }
    label { display: inline-block; width: 100px; }
    input { padding: 8px; margin-right: 10px; }
    button { padding: 8px 15px; background: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer; }
    button:hover { background: #45a049; }
    #checkResult, #listenStatus { margin-top: 10px; padding: 10px; border-radius: 4px; }
    .success { background: #ddffdd; border-left: 6px solid #4CAF50; }
    .error { background: #ffdddd; border-left: 6px solid #f44336; }
  </style>
</head>
<body>
  <h1>PortCheck 工具</h1>

  <h2>① 启动服务端端口监听</h2>
  <form id="listenForm">
    <label>监听端口：</label>
    <input type="text" id="listenPort" required>
    <button type="submit">启动监听</button>
  </form>
  <p id="listenStatus"></p>

  <h2>② 检查该端口是否可达（浏览器请求）</h2>
  <form id="checkForm">
    <label>目标端口：</label>
    <input type="text" id="checkPort" required>
    <button type="submit">发起探测</button>
  </form>
  <p id="checkResult"></p>

  <script>
    document.getElementById('listenForm').onsubmit = async (e) => {
      e.preventDefault();
      const port = document.getElementById('listenPort').value;
      const statusEl = document.getElementById('listenStatus');
      
      statusEl.textContent = "启动中...";
      statusEl.className = "";
      
      try {
        const res = await fetch("/listen", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: "port=" + encodeURIComponent(port),
        });
        const text = await res.text();
        statusEl.textContent = text;
        statusEl.className = "success";
      } catch (err) {
        statusEl.textContent = "错误: " + err.message;
        statusEl.className = "error";
      }
    };

    document.getElementById('checkForm').onsubmit = async (e) => {
      e.preventDefault();
      const port = document.getElementById('checkPort').value;
      const resultEl = document.getElementById('checkResult');
      
      resultEl.textContent = "检测中...";
      resultEl.className = "";
      
      try {
        const url = location.protocol + '//' + location.hostname + ':' + port + '/';
        const res = await fetch(url);
        const text = await res.text();
        resultEl.textContent = "✅ 成功访问端口 " + port + ": " + text;
        resultEl.className = "success";
      } catch (err) {
        resultEl.textContent = "❌ 无法访问端口 " + port + ": " + err.message;
        resultEl.className = "error";
      }
    };
  </script>
</body>
</html>
`
