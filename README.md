# 简单web 应用与服务——Golang

本次选择使用较轻量级组件negroni + gorilla/mux实现简单web应用，该应用实现以下功能：

- 支持静态文件服务
- 支持简单 js 访问
- 提交表单，并输出一个表格
- 对 `/unknown` 给出开发中的提示，返回码 `5xx`

## 前言

实现简单web应用于服务的基础是golang官方提供的http包，可以通过[Golang web 应用开发](https://github.com/astaxie/build-web-application-with-golang)进行学习。

## 框架选择

- [negroni ](http://github.com/codegangsta/negroni)
- [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux)

## 项目结构

- controller
- public
  - css
  - img
  - js
  - template
- router

`controller`用于管理handler

`public`为静态文件托管的目录

`router`为路由控制

## 功能实现

接下来介绍功能的实现

### 支持静态文件服务

支持静态文件服务功能的方法：

- 通过negroni内置中间件实现

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/urfave/negroni"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome to the home page!")
  })

  // http.FileServer的使用例子, 若你预期要"像伺服器"而非"中间件"的行为
  // mux.Handle("/public", http.FileServer(http.Dir("/home/public")))

  n := negroni.New()
  n.Use(negroni.NewStatic(http.Dir("/public")))
  n.UseHandler(mux)

  http.ListenAndServe(":3002", n)
}
```

- 通过mux路由router的Handler

```go
r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(dir))))
```

实现了静态文件服务之后，可以通过`localhost:port/public/`访问。

### 支持简单js访问

访问首页时首先使用ajax向服务发送一个路径为`/time`的请求，并且将返回的数据显示到页面中。

```html
<!-- index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <link href="../css/index.css" rel="stylesheet" type="text/css" />
    <script src="http://code.jquery.com/jquery-latest.js"></script>
    <script src="../js/hello.js"></script>
    <title>CloudIO-GOLANG</title>
</head>
<body>
    <div>
        <img id="image" src="../img/background.jpg">
    </div>
    <div id="title">Welcome to CookiesChen's MiaoMiao House</div>
    <p class="time">Now time is {{.Time}}</p>
    <form action="/form" method="post">
        <div id="username">
            username <input type="text" name="username">
        </div>
        <div id="password">
            password <input type="password" name="password">
        </div>
        <div id="login">
            <input type="submit" value="登录">
        </div>
    </form>
</body>
</html>
```

```js
/* hello.js */
$(document).ready(function() {
    $.ajax({
        url: "/time"
    }).then(function(data) {
        $('.time').append(data.time);
    });
});
```

在服务器端，首先注册路由

```go
R.HandleFunc("/time", controller.TimeHandler).Methods("GET")
```

然后是Handler的实现，是处理请求的时候，获取当前系统的时间，然后使用`render`将JSON写入到`ResponseWriter`

```java
func TimeHandler(w http.ResponseWriter, r *http.Request)  {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	formatter.JSON(w, http.StatusOK, struct {
		Time string `json:"time"`
	}{Time: time.Now().String()})
}
```

### 处理表单输入

表单提交如下

```html
<form action="/form" method="post">
    <div id="username">
        username <input type="text" name="username">
    </div>
    <div id="password">
        password <input type="password" name="password">
    </div>
    <div id="login">
        <input type="submit" value="登录">
    </div>
</form>
```

```go
if r.Method == "POST" {
    r.ParseForm()
    fmt.Println("username:", r.Form["username"])
    fmt.Println("password:", r.Form["password"])
    t, _ := template.ParseFiles("./public/template/form.html")
    log.Println(t.Execute(w, struct {
        Name      string `json:"name"`
        Password string `json:"password"`
    }{Name: r.Form["username"][0], Password: r.Form["password"][0]}))
}
```

处理表单的时候，需要显式解析表单`r.ParseForm()`，默认不会解析，解析之后通过键值可以获取到对应的`slice`。最后使用模板进行输出即可。

### 对于未注册路径

在handler中简单的输出错误即可。

```go
http.Error(w, "no such directory", 500)
```

### Negroni中间件源码阅读

本质上来说`Negroni`是一个HTTP Handler,因为他实现了HTTP Handler接口，所以他可以被`http.ListenAndServe`使用，其次`Negroni`本身内部又有一套自己的Handler处理链，通过他们可以达到处理http请求的目的，这些Handler处理链中的处理器，就是一个个中间件。

```go
func (n *Negroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n.middleware.ServeHTTP(NewResponseWriter(rw), r)
}
```

再来看看Negroni的Handler，可以看到，与http.Handler唯一不同的就是多了一个`next`，而这就是中间件实现的基础。

```go
type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
```

构建中间件处理链，Negroni有一个`handler`的`slice`，通过调用`Use`方法，可以将handler依次保存在该slice中。

```go
type Negroni struct {
	middleware middleware
	handlers   []Handler
}

func (n *Negroni) Use(handler Handler) {
	if handler == nil {
		panic("handler cannot be nil")
	}

	n.handlers = append(n.handlers, handler)
	n.middleware = build(n.handlers)
}
```

当handler的数量大于1时， 通过递归构造，这样先添加的handler会排在前面执行。当等于0时，返回空的中间件。当等于1时，没有next，将next设置为空的中间件即可。

```go
type middleware struct {
	handler Handler
	next    *middleware
}

func build(handlers []Handler) middleware {
	var next middleware

	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middleware{},
	}
}
```

编写自己的中间件

```go
package main

import (
	"fmt"
	"github.com/CookiesChen/cloudgo-io/router"
	"github.com/urfave/negroni"
	"net/http"
)

const port = "9090"

func main() {
	r := router.R

	n := negroni.Classic()

	n.UseFunc(sayhi)
	n.UseFunc(sayhello)

	n.UseHandler(r)
	n.Run(": " + port)
}

func sayhi(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	fmt.Println("Hi")
	next(rw,r)
}

func sayhello(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
	fmt.Println("Hello")
	next(rw,r)
}
```

以上每次访问服务器时，都会链式执行中间件，在本次实验中并没有进行有效的中间件的编写，在实际应用时，可以通过中间件起到filter的作用，用于白名单和用户身份验证等方面，具体实现还需要根据应用要求。