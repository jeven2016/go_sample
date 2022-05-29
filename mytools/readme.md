### 安装环境
* Ubuntu安装linux依赖
```shell
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev
```

* 安装CLI
Wails CLI 可以使用 go get 安装。安装之后，你应该使用 wails setup 命令进行设置
```shell
 go install github.com/wailsapp/wails/v2/cmd/wails@latest
 
 ln -s /home/jujucom/Desktop/workspace/install/go/bin/wails /usr/bin/wails

```

* 初始化工程
```shell
wails init -n mytools -t react

wails init -n mytools -t react-ts
```
## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. The frontend dev server will run on http://localhost:34115. Connect to this in your
browser and connect to your application.

## Building

To build a redistributable, production mode package, use `wails build`.  
参考地址：https://wails.io/zh-Hans/docs/reference/cli/
```shell
//编译windows版本， 自带webview2运行库
wails build -platform windows/amd64 -clean -webview2 embed

#编译mac
wails build -platform  darwin/amd64

```



