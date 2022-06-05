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
wails init -n mytools -t react -ide vscode/goland

# 创建工程并配置goland工具，可以直接ide中调试go源码
wails init -n my-go-tools -t react-ts -ide goland

//可以参考的配置
wails init -n [Your Appname] -t https://github.com/AlienRecall/wails-react-template -ide goland

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

## 创建使用新的web工程，取代原有的frontend
* 修改main.go修改工程路径
```shell
//go:embed gui/build
var assets embed.FS
```
* 修改walls.json
```shell
{
  "name": "my-go-tools",
  "outputfilename": "my-go-tools",
  "gui:install": "npm install",
  "gui:build": "npm run build",
  "gui:dev:watcher": "npm run dev",
  "gui:dev:serverUrl": "http://localhost:3000",
  "author": {
    "name": "jujucom",
    "email": "jujucom@126.com"
  }
}
```

### 修改文件后没有reloading
默认情况下，wails使用build后的文件。但是create-react-app没有自动监控文件变化进行build文件输出的功能。
所以需要引入第三方包cra-build-watch.



