[原文接接](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

跨域资源共享(Cross-Origin Resource Sharing)

跨域资源共享(CORS)是一种机制，它使用额外的 HTTP 头来告诉浏览器使运行于一个源上的 web 应用程序，从一个不同的域上访问所选择的资源。当 web 应用程序请求不同来源(域，协议或端口)的资源时，将执行 **跨域 HTTP 请求**。

跨域请求例子：来自 `https://domain-a.com` 的JavaScript 后端代码使用 [XMLHttpRequest](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest) 请求 `https://domain-b.com/data.json`。 

出于安全的原因，浏览器限制了从脚本发起的 HTTP 跨域请求。例如，`XMLHttpRequest` 和 [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) 遵循 [同源策略](https://developer.mozilla.org/en-US/docs/Web/Security/Same-origin_policy)。这意味着，web 应用程序使用这些 API 仅可以访问应用程序加载时的同源资源，除非其他源的响应包含正确的 CORS 头。

![cross-origin requests(controlled by CORS)](https://mdn.mozillademos.org/files/14295/CORS_principle.png)

CORS 机制支持浏览器和服务器之间的安全跨域请求和数据转移。现代浏览器在 `XMLHttpRequest` 或 [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) 类的 API 中使用CORS，来减轻跨源 HTTP 请求的风险。

# 谁应该阅读该文章？

实际上，每个人。

进一步来说，本文章适用于 **web管理员，服务器开发人员和前端开发人员**，现代浏览器处理跨域共享的客户端，包括头部和策略实施。但是 CORS 标准意味着服务器必须新的请求和响应头。服务器开发人员讨论的[从服务器角度进行趺域共享(带有PHP代码段)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Server-Side_Access_Control)作为补充阅读。

# 什么请求使用 CORS?

[跨域共享标准](https://fetch.spec.whatwg.org/#http-cors-protocol) 可以为以下站点启用跨域 HTTP 请求：

* 如上所述，对 [XMLHttpRequest](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest) 或 [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) API。
* Web 字段（在CSS中 `@font-face` 使用的跨域字体），[这样服务器就可以部署 `TrueType` 字体，该字体只能跨站点加载并允许被允许的网站使用](https://www.w3.org/TR/css-fonts-3/#font-fetching-requirements)。

* [WebGL纹理](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial/Using_textures_in_WebGL)。
* 使用 [drawImage()](https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/drawImage) 绘制到画布上的 Images/video 帧。
* [图片的 CSS 形状](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Shapes/Shapes_From_Images)。

这篇文章只讨论跨域资源共享，包含一些必要的HTTP头。

# 功能概述

跨域资源共享（CORS）标准的工作原理是通过增加新的[HTTP 头](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers)。它允许服务器描述允许哪个来源从 Web 浏览器读取该信息。另外，对于可能会导致服务器数据产生副作用的HTTP请求方法(尤其是 [GET](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/GET) 或具有 [MIME types](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types) 类型的 [POST](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/POST) 以外的HTTP方法)，该规范要求浏览器“预检（preflight）" 请求，并使用 HTTP 的 [OPTIONS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS) 请求方法从服务器请求支持的方法。在服务器“认可(approval)”后，发送实际的请求。服务器也可以通知客户端是否应将“credentials"(例如 [Cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies)和 [HTTP认证](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication))和请求一起发送。

CORS 失败导致错误，但是出于安全的原因，关于错误的细节，对于*JavaScript*来说是*不可访问的*。代码所知道的是错误发生了。确定具体问题的唯一算途径是查看浏览器的控制台来获取详细信息。

后续的部分讨论了方案，并提供了所有HTTP头部使用的细分。

# 访问控制场景的示例

我们提供了三种方案来演示跨域资源共享的工作方式。这些示例都使用 [XMLHttpRequest](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest)，它可以在任何支持的浏览器中发出跨站点请求。

这些部分(服务器中这些跨域站点请求的正确处理的代码运行实例) 的 JavaScript 片断可以在 [http://arunranga.com/examples/access-control/](http://arunranga.com/examples/access-control/) 中的 ”in action" 中找到，将工作在支持跨域站点 `XMLHttpRequest` 的浏览器中。

从服务器角度讨论跨域资源共享（包含一些PHP片断）可以在 [服务器端访问控制（CORS）](https://developer.mozilla.org/en-US/docs/Web/HTTP/Server-Side_Access_Control)的文章中找到

## 简单请求(simple request)

一些请求不会触发[CORS preflight](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Preflighted_requests) 。这些在该文章中被称为 *简单请求(simple request)*，尽管 [Fetch](https://fetch.spec.whatwg.org/)规范(它定义了CORS)没有使用该术语。*简单请求* 是指**满足以下所有条件**的请求：

* 允许的方法之一：
  + [GET](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/GET)
  + [HEAD](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/HEAD)
  + [POST](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/POST)
* 除了由用户代理自动设置的头部(例如，[Connect](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Connection)，[User-Agent](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent)，或[在Fetch规范中定义为"禁止头部名称"外的其他头](https://fetch.spec.whatwg.org/#forbidden-header-name))外，唯一允许手动设置的头部是[Fetch 规定定义为”CORS-满足的请求头](https://fetch.spec.whatwg.org/#cors-safelisted-request-header)“，它们是:
  + [Accept](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept)
  + [Accept-Language](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language)
  + [Content-Language](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Language)
  + [Content-Type](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type) (但是注意以下附加要求)
  + [DPR](http://httpwg.org/http-extensions/client-hints.html#dpr)
  + [Downlink](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Downlink)
  + [Save-Data](http://httpwg.org/http-extensions/client-hints.html#save-data)
  + [Viewport-Width](http://httpwg.org/http-extensions/client-hints.html#viewport-width)
  + [Width](http://httpwg.org/http-extensions/client-hints.html#width)
* [Content-Type](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type)头部唯一允许的值是：
  + `application/x-www-form-urlencoded`
  + `multipart/form-data`
  + `text/plain`
* 在任何请求使用的[XMLHttpRequestUpload](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequestUpload)中没有注册事件监听，这些可以使用 [XMLHttpRequest.upload](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/upload)属性访问。
* 在请求中没有使用 [ReadableStream](https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream)

> **注意：**这些是 Web 内容已经可以发出的相同类型的跨站点请求，除非服务器发送适当的头部，否则不会向请求者释放响应数据。因此，附上跨站点请求伪造的站点不必担心 HTTP 访问控制。

> **注意：** WebKit Nightly 和 Safari 技术预览版对 [Accept](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept)，[Accept-Language](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language) 和 [Content-Language](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Language) 头中允许的值设置了其他限制。如果这些头的任何一个拥有 “nonstandard" 值，WebKit/Safari 并不认为请求是”简单请求“。除了以下 bug 外，没有文档记录WebKit/Sarafi 认为的 "nonstandard“ 值：
>
> - [Require preflight for non-standard CORS-safelisted request headers Accept, Accept-Language, and Content-Language](https://bugs.webkit.org/show_bug.cgi?id=165178)
> - [Allow commas in Accept, Accept-Language, and Content-Language request headers for simple CORS](https://bugs.webkit.org/show_bug.cgi?id=165566)
> - [Switch to a blacklist model for restricted Accept headers in simple CORS requests](https://bugs.webkit.org/show_bug.cgi?id=166363)
>
> 其他浏览器没有实现这些额外的限制，因为他们不是规范的一部分。

例如，假设位于 `https://foo.example` 的 web 内容希望调用在域名 `https://bar.other` 中的内容。这类代码可能被部署在 `foo.example` 中的 JavaScript 使用。

> ```javascript
> const xhr = new XMLHttpRequest();
> const url = 'https://bar.other/resources/public-data/';
>    
> xhr.open('GET', url);
> xhr.onreadystatechange = someHandler;
> xhr.send(); 
> ```

这个演示使用 CORS 头处理特权，在服务器和客户端进行简单的交换：

![client-server-pic](https://mdn.mozillademos.org/files/14293/simple_req.png)

我们来看一下这个例子中，浏览器将向服务器发送什么，以及服务器如何响应：

> GET /resources/public-data/ HTTP/1.1
> Host: bar.other
> User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
> Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
> Accept-Language: en-us,en;q=0.5
> Accept-Encoding: gzip,deflate
> Connection: keep-alive
> **Origin: https://foo.example**

需要注意的请求头是 [Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin)，它展示了从 `https://foo.example` 过来的调用。

> HTTP/1.1 200 OK
> Date: Mon, 01 Dec 2008 00:23:53 GMT
> Server: Apache/2
> **Access-Control-Allow-Origin: * **
> Keep-Alive: timeout=2, max=100
> Connection: Keep-Alive
> Transfer-Encoding: chunked
> Content-Type: application/xml
>
> […XML Data…]

在响应中，服务器返回一个 [https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin 头。[Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin) 和 [Access-Control-Allow-Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin) 的使用展示了访问控制协议最简单的使用。在这个例子中，服务器以`Access-Control-Allow-Origin: *` 响应，它意味着资源可以被**任何**域名访问。如果位于 `https://bar.other` 的资源所有者希望限制资源的访问，仅允许从 `https://foo.example` 的请求访问，它将发送：

> Access-Control-Allow-Origin: https://foo.example

现在除了 `https://foo.example` 外，不允许其他域名以跨域的方法访问资源。`Access-Control-Allow-Origin` 的头部应该包含请求的 `Origin` 头中发送的值。

## 预检请求(Preflighted requests)

不像 ”[简单请求(上面讨论的)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Simple_requests)，“preflighted" 请求首先使用 [OPTIONS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS) 方法发送一个到其他域资源的 HTTP 请求，来决定实际的请求的发送是否安全。跨站点请求这些被预检，因为它们可能会影响到用户数据。

以下是会被预检的请求示例：

> ```javascript
> const xhr = new XMLHttpRequest();
> xhr.open('POST', 'https://bar.other/resources/post-here/');
> xhr.setRequestHeader('Ping-Other', 'pingpong');
> xhr.setRequestHeader('Content-Type', 'application/xml');
> xhr.onreadystatechange = handler;
> xhr.send('<person><name>Arun</name></person>'); 
> ```

上面的例子创建一个 XML 正文，使用 `POST` 请求发送。同时，一个非标准的 HTTP `Ping-Other` 请求头被设置。这样的头不是 HTTP/1.1 的一部分，但是通常对 web 应用程序有用。由于请求的 Content-Type 设置为 `application/xml`，并且由于设置了自定义的头，所以该请求被预检。

![client-server-preflight.png](https://mdn.mozillademos.org/files/16753/preflight_correct.png)

(注意：正如上面描述的，实际的 `POST` 请求并不包含 `Access-Control-Request-*` 头；它们只被 `OPTIONS` 请求需要)

我们来看一下客户端和服务器间整个交换过程，第一个交换是 *预检请求/响应*:

> OPTIONS /resources/post-here/ HTTP/1.1
> Host: bar.other
> User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
> Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
> Accept-Language: en-us,en;q=0.5
> Accept-Encoding: gzip,deflate
> Connection: keep-alive
> Origin: http://foo.example
> Access-Control-Request-Method: POST
> Access-Control-Request-Headers: X-PINGOTHER, Content-Type
>
> HTTP/1.1 204 No Content
> Date: Mon, 01 Dec 2008 01:15:39 GMT
> Server: Apache/2
> Access-Control-Allow-Origin: https://foo.example
> Access-Control-Allow-Methods: POST, GET, OPTIONS
> Access-Control-Allow-Headers: X-PINGOTHER, Content-Type
> Access-Control-Max-Age: 86400
> Vary: Accept-Encoding, Origin
> Keep-Alive: timeout=2, max=100
> Connection: Keep-Alive

当预检请求完成，发送真正的请求：

> POST /resources/post-here/ HTTP/1.1
> Host: bar.other
> User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
> Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
> Accept-Language: en-us,en;q=0.5
> Accept-Encoding: gzip,deflate
> Connection: keep-alive
> X-PINGOTHER: pingpong
> Content-Type: text/xml; charset=UTF-8
> Referer: https://foo.example/examples/preflightInvocation.html
> Content-Length: 55
> Origin: https://foo.example
> Pragma: no-cache
> Cache-Control: no-cache
>
> <person><name>Arun</name></person>
>
> HTTP/1.1 200 OK
> Date: Mon, 01 Dec 2008 01:15:40 GMT
> Server: Apache/2
> Access-Control-Allow-Origin: https://foo.example
> Vary: Accept-Encoding, Origin
> Content-Encoding: <font color=red>gzip</font>
> Content-Length: 235
> Keep-Alive: timeout=2, max=99
> Connection: Keep-Alive
> Content-Type: text/plain
>
> [Some GZIP'd payload]

上面的第1行到11行代表使用 [OPTION](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS) 方法的预检请求。浏览器根据上面的JavaScript代码段所使用的请求参数确定需要发送此请求，服务器可以用实际的请求参数响应是否可以发送请求。OPTIONS 是 HTTP/1.1 方法，它被用来确定从服务器接受更多的信息，这是一个[安全](https://developer.mozilla.org/en-US/docs/Glossary/safe)的方法，意味着它不可以被用来修改资源。请注意，与OPTIONS请求一起，还发送了另外两个请求标头（分别是第10行和第11行）：

> Access-Control-Request-Method: POST
> Access-Control-Request-Headers: X-PINGOTHER, Content-Type

[Access-Control-Request-Method](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method) 头将作为预检请求的一部分通知服务器，当实际请求发送时，它将以 `POST` 请求方法发送。  [Access-Control-Request-Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers) 头通知服务器，当实际请求被发送时，它将与 `X-PINGOTHER` 和 `Content-Type` 自定义头一块儿发送。现在服务器有机会确定它是否接受这种情况下的请求。

上面的第 14-23 行，是服务器返回的响应，表明请求方法（`POST`）和请求头（`X-PINGOTHER`）是可接受的。我们来特别关注下第 17-20 行：

> Access-Control-Allow-Origin: http://foo.example
> Access-Control-Allow-Methods: POST, GET, OPTIONS
> Access-Control-Allow-Headers: X-PINGOTHER, Content-Type
> Access-Control-Max-Age: 86400

服务器使用 `Access-Control-Allow-Methods` 响应，并告诉 `POST` 和 `GET` 是查询问题中资源的可行的方法，注意这个头部和 [Allow](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Allow) 响应头是类似的，但是严格的在访问控制的上下文中使用。

服务器也发送了 `Access-Control-Allow-Headers` ，值为 "`X-PINGOTHER`,`Content-Type`"，确认了这些是允许的头，可以在实际的请求中被使用。像 `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers` 是逗号分隔的可接受头的列表。

最后，[Access-Control-Max-Age](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age) 以秒为单位给出了预检请求的响应在未发送另一个预检请求时可以缓存多长时间。在这个例子中，86400秒是24小时。注意，每个浏览器都有一个[最大间隔值](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age)，当 `Access-Control-Max-Age` 更大的时候，它有优先权。

## 预检请求和重定向

目前并不是所有的浏览器支持预检请求后的下列重定向。在预检请求后，如果发生了重定向，当前有的浏览器会报告一个类似下列的错误消息。

> The request was redirected to 'https://example.com/foo', which is disallowed for cross-origin requests that require preflight

> Request requires preflight, which is disallowed to follow cross-origin redirect

CORS 协议最初要求该行为，但[后来更改为不再需要它](https://github.com/whatwg/fetch/commit/0d9a4db8bc02251cc9e391543bb3c1322fb882f2)。然而，不是所有的浏览器实现这个更改，因此仍然表现出最初所需的行为。

因此，在所有浏览器都赶走规范之前，您可以通过执行以下一项或两项操作来解决此限制：

* 修改服务器端的行为以避免预检和/或避免重定向--如果你有待请求服务的控制权
* 修改请求，例如，它是一个 [简单请求](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Simple_requests)，不会导致预检。

但是，如果不可能做这些修改，另一个可能的办法是：

* 使用[简单请求](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Simple_requests)（对于 Fetch  API 用 [Response.url](https://developer.mozilla.org/en-US/docs/Web/API/Response/url) 或 [XMLHttpRequest.responseURL](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest/responseURL)）确定预检请求将以哪个 URL 结束。
* 使用第一步您从 `Response.url` or `XMLHttpRequest.responseURL` 中获取到的URL 构造另一个请求("真正"的请求)。

然而，如果由于在请求中存在 `Authorization` 头触发一个预检，使用上面的步骤你不能够解决这种限制。除非您可以控制请求的服务器，否则您将根本无法解决它。

## 带有凭证的请求(Requests with credentials)

[XMLHttpRequest](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest) 或 [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) 和CORS公开的最有趣的功能是能够发出知道 [HTTP cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies) 和HTTP身份验证信息的“凭证”请求的功能。默认地，跨站点 `XMLHttpRequest` 或 [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) 调用，浏览器将 **不** 发送凭证。当它调用时，必须在 `XMLHttpRequest` 对象或 [Request](https://developer.mozilla.org/en-US/docs/Web/API/Request) 创建者必须设置一个特定的标志。

在这个例子中，从 `http://foo.example` 加载的源内容发起了一个到 `http://bar.other` 资源的简单 GET 请求，设置了 Cookies。在 foo.example 上的内容可能包含像这样的 JavaScript ：

> ```javascript
> const invocation = new XMLHttpRequest();
> const url = 'http://bar.other/resources/credentialed-content/';
>     
> function callOtherDomain() {
>   if (invocation) {
>     invocation.open('GET', url, true);
>     invocation.withCredentials = true;
>     invocation.onreadystatechange = handler;
>     invocation.send(); 
>   }
> }
> ```

第 7 行展示了 [XMLHttpRequest](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest) 上的标志必须被设置，以便使用 Cookie 调用，即 `withCredentials` 的布尔值。默认地，调用不携带 Cookies。由于这是一个简单的 `GET` 请求，它是不预检的，但是浏览器将**拒绝**任何没有 [Access-Control-Allow-Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials) ： `true` 的头部，并**不**响应提供给调用的Web内容。

![client-server-credentials](https://mdn.mozillademos.org/files/14291/cred-req.png)



下面是服务端和客户端之间交互的示例：

> GET /resources/access-control-with-credentials/ HTTP/1.1
> Host: bar.other
> User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0
> Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
> Accept-Language: en-us,en;q=0.5
> Accept-Encoding: gzip,deflate
> Connection: keep-alive
> Referer: http://foo.example/examples/credential.html
> Origin: http://foo.example
> Cookie: pageAccess=2
>
> HTTP/1.1 200 OK
> Date: Mon, 01 Dec 2008 01:34:52 GMT
> Server: Apache/2
> Access-Control-Allow-Origin: https://foo.example
> Access-Control-Allow-Credentials: true
> Cache-Control: no-cache
> Pragma: no-cache
> Set-Cookie: pageAccess=3; expires=Wed, 31-Dec-2008 01:34:53 GMT
> Vary: Accept-Encoding, Origin
> Content-Encoding: <font color=red>gzip</font>
> Content-Length: 106
> Keep-Alive: timeout=2, max=100
> Connection: Keep-Alive
> Content-Type: text/plain
>
> [text/plain payload]

尽管第10行包含目的为 `http://bar.other` 上内容的 Cookie，如果 bar.otehr 不使用`[Access-Control-Allow-Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials): true (17行)响应，响应将被忽略并不返回 web 内容。

### 凭证请求和通配符

当响应一个 credentials 请求，服务器**必须**在头部 `Access-Control-Allow-Origin` 的值中指定一个源，而不是指定一个"`*`" 通配符。

由于上面示例中的请求头中包含一个 `Cookie` 头，如果 `Access-Control-Allow-Origin` 的值是 ”`*`“，请求将失败。但是它没有失败：因为 `Access-Control-Allow-Origin` 的值是 ”`http://foo.example`"（实际的源）而不是“`*`"通配符，凭据识别内容将返回到调用Web内容。

请注意，上面示例中的 `Set-Cookie` 响应标头还设置了另一个cookie。如果发生故障，则会根据所使用的API引发异常。

### 第三方cookies

注意，在 CORS 响应中设置的 cookies 必须遵守正常的第三方cookie策略。在上面的例子中，页从 `foo.example` 被加载，但是第22行的cookie是被 `bar.other` 发送的，因此，如果用户已将其浏览器配置为拒绝所有第三方Cookie，则将不会保存。

# HTTP响应头

这个部分列出了服务器会发回的在跨域资源共享规范中定义的访问控制请求中HTTP响应头。上一节概述了这些内容。

## Access-Control-Allow-Origin

返回的资源可能会有一个 [Access-Control-Allow-Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin) 头，语法如下：

> Access-Control-Allow-Origin: <origin> | *

`Access-Control-Allow-Origin` 要么指定一个单一的源，告诉浏览器允许源访问资源；或者--对于**没有**凭证的请求 -- ”`*`“通配符，告诉浏览器允许所有的源访问资源。

例如，为了允许来自源 `https://mozilla.org` 的代码访问资源，你会指定：

> Access-Control-Allow-Origin: https://mozilla.org

如果服务器指定一个单一的 origin 而不是 ”`*`" 通配符，服务器应该在 [Vary](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Vary) 响应头中也包含 `Origin` -- 来向客户端表明服务器基于 [Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin) 请求头的值而有所不同。

## Access-Control-Expose-Headers

[Access-Control-Expose-Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers) 让服务器将浏览器允许访问的白名单。

> Access-Control-Expose-Headers: <header-name>[, <header-name>]*

例如，下列

>Access-Control-Expose-Headers: X-My-Custom-Header, X-Another-Custom-Header

将允许 `X-My-Custom-Header` 和 `X-Another-Custom-Header` 头暴露给浏览器。

## Access-Control-Max-Age

[Access-Control-Max-Age](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age) 表明预检请求的结果可以被缓存多长时间

> Access-Control-Max-Age: <delta-seconds>

 `delta-seconds` 参数以秒为单位的数值表明结果可以被缓存的时间。

## Access-Control-Allow-Credentials

[Access-Control-Allow-Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials) 头表明当 `credentials` 标志为真是请求的响应是否可以暴露。当作为预检请求的一部分被使用时，表明实际的请求是否使用凭证。注意简单的 `GET` 请求是不预检的，因此，如果请求具有凭据的资源，则此标头如果未随资源一起返回，浏览器将忽略该响应，并且不将其返回到Web内容。

> Access-Control-Allow-Credentials: true

[Credentialed 请求](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Requests_with_credentials) 在上面被讨论过。

## Access-Control-Allow-Methods

[`Access-Control-Allow-Methods`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods) 头指定了方法或当访问资源时允许的方法。这在预检请求的响应中被使用。上面讨论了请求被预检的条件。

> Access-Control-Allow-Methods: <method>[, <method>]*

[预检请求](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Preflighted_requests)的示例在上面被讨论

## Access-Control-Allow-Headers

[`Access-Control-Allow-Headers`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers) 头被用于[preflight request](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Preflighted_requests) 的响应来表明当进行实际的请求时哪一个HTTP头可以被使用。

> Access-Control-Allow-Headers: <header-name>[, <header-name>]*

# HTTP 请求头

该部分列出了为了使用跨域资源共享特性发出HTTP请求，客户端可能会使用的头。注意，当调用服务器时，这些头将会为你设置，使用跨站点[`XMLHttpRequest`](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest)功能的开发者不必使用编程的方式设置任何跨域请求头。

## Origin

[`Origin`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin) 表示了跨站点访问请求或预检请求的源

> Origin: <origin>

源是一个URI，指示从中发起请求的服务器。它不包括任何路径信息，仅包含服务器的名字。

> **注意：** `origin` 的值可以是 `null`，或一个 URI。

注意，在任何访问控制请求中，[`Origin`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin) **总是**被发送。

## Access-Control-Request-Method

当发起预检请求时，[`Access-Control-Request-Method`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method) 被使用，使服务器知道当发起实际的请求时，将使用什么HTTP方法。

> Access-Control-Request-Method: <method>

这个的使用示例可以在[上面找到](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Preflighted_requests)

## Access-Control-Request-Headers

当发起预检请求时，使用[`Access-Control-Request-Headers`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers) 服务器知道当真正发起请求时，会使用什么HTTP头

> Access-Control-Request-Headers: <field-name>[, <field-name>]*

用法示例可以在[上面找到](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#Preflighted_requests)

# 相关阅读

- [CORS errors](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS/Errors)
- [Enable CORS: I want to add CORS support to my server](https://enable-cors.org/server.html)
- [`XMLHttpRequest`](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest)
- [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
- [Using CORS with All (Modern) Browsers](http://www.kendoui.com/blogs/teamblog/posts/11-10-03/using_cors_with_all_modern_browsers.aspx)
- [Using CORS - HTML5 Rocks](http://www.html5rocks.com/en/tutorials/cors/)
- [Code Samples Showing `XMLHttpRequest` and Cross-Origin Resource Sharing](https://arunranga.com/examples/access-control/)
- [Client-Side & Server-Side (Java) sample for Cross-Origin Resource Sharing (CORS)](https://github.com/jackblackevo/cors-jsonp-sample)
- [Cross-Origin Resource Sharing From a Server-Side Perspective (PHP, etc.)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Server-Side_Access_Control)
- [Stack Overflow answer with “how to” info for dealing with common problems](https://stackoverflow.com/questions/43871637/no-access-control-allow-origin-header-is-present-on-the-requested-resource-whe/43881141#43881141)：
  - 如何避免CORS预检
  - 如何使用CORS代理解决 *“No Access-Control-Allow-Origin header”* 问题
  - 如何修复 *“Access-Control-Allow-Origin header must not be the wildcard”*
