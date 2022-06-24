### library

```shell

# rest client
go get github.com/go-resty/resty/v2

# uuid
go get github.com/google/uuid
```

### openconnect ID

https://www.jianshu.com/p/d453076e6433
https://blog.csdn.net/wdquan19851029/article/details/111887107
https://blog.csdn.net/weixin_45784983/article/details/106716433?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-1-106716433-blog-111887107.pc_relevant_default&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-1-106716433-blog-111887107.pc_relevant_default&utm_relevant_index=2
https://baeldung-cn.com/spring-boot-keycloak
Keycloak快速上手指南，只需10分钟即可接入Spring Boot/Vue前后端分离应用实现SSO单点登录

### access type

关于客户端的访问类型（Access Type）
上面创建的2个客户端的访问类型分别是public、bearer-only，那么为什么分别选择这种类型，实际不同的访问类型有什么区别呢？

事实上，Keycloak目前的访问类型共有3种：

confidential：适用于服务端应用，且需要浏览器登录以及需要通过密钥获取access token的场景。典型的使用场景就是服务端渲染的web系统。

public：适用于客户端应用，且需要浏览器登录的场景。典型的使用场景就是前端web系统，包括采用vue、react实现的前端项目等。

bearer-only：适用于服务端应用，不需要浏览器登录，只允许使用bearer token请求的场景。典型的使用场景就是restful api