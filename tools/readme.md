## 工具说明
当前工程包含了两类工具：  
* dist/move-repository   
  从Nexus下载包并上传至JFrog仓库中
* dist/iam-convert  
  将OA上导出的部门和用户数据转换成keycloak的数据格式，便于导入

## 工具使用
### move-repository使用
参见文档：

### 转换OA数据


## 编译工具
### 编译move-repository工具
编译成功后，在同级目录下会创建一个可执行文件move-repository
```shell
 go build move-repository.go
```
### 编译iamb-convert工具
编译成功后，在同级目录下会创建一个可执行文件move-repository
```shell
 go build move-repository.go
```


### 在linux下编译其他操作系统可执行文件
```shell
# Mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build move-repository.go

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build move-repository.go

```