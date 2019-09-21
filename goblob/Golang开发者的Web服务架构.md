Web服务架构是构建每个项目的第一阶段，这就像是你建造一个房屋前，先有计划图纸一样。

本文将向你显现，在使用Golang创建一个简单的web服务时，我是如何构架项目的。保持一个简单而直观的架构是非常重要的，我们知道，在 golang 中你可以通过引用包名来调用方法。

接下来，我将展示一个简单，但传统的Web服务体系结构模型，在我参与的大多数项目中，我使用它，处理每个Web服务的组件

# 一、/API

API 包是个文件夹，你可以按其服务目的分组把所有 API 端点放置到其子包中。这意味着，我更喜欢有一个特殊的包，它的主要范围是解决一个特定的问题。

例如，所有的登录，注册，忘记密码，重置密码处理，我喜欢把字们定义在名为 **registration** 的包中。

该 `registration` 包看起来如下：

```javascript
.
├── api
│   ├── auth
│   │   ├── principal.middleware.go
│   │   └── jwt.helper.go
│   ├── cmd
│   │   └── main.go
│   ├── registration
│   │   ├── login.handler.go
│   │   ├── social_login.handler.go
│   │   ├── register.handler.go
│   │   ├── social_register.handler.go
│   │   ├── reset.handler.go
│   │   ├── helper.go
│   │   └── adapter.go
├── cmd
│   └── main.go
├── config
│   ├── config.dev.json
│   ├── config.local.json
│   ├── config.prod.json
│   ├── config.test.json
│   └── config.go
├── db
│   ├── handlers
│   ├── models
│   ├── tests
│   ├── db.go
│   └── service.go
├── locales
│   ├── en.json
│   └── fr.json
├── public
├── vendor
├── Makefile
..........................
```



## handler.go

正如你所看到的，它们的文件名是以 **handler.go** 结尾的。在里面，你可以有效地写代码，它将处理请求，即请求数据将从数据库中检索，处理，最后组成响。

下面是一个可以很好解析的例子。

```go
http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request){
  //handle the request
}
```

## helper.go

有时候，在发送响应之前，你需要从多个地方搜索数据并处理它们，然后，当所有的细节都集合到一块，响应可以被发送到客户端 app。但是代码处理必须保持尽可能的简单，所以所有处理部分的额外代码可以放到这儿。

## adaper.go

在客户端和web服务间的交互中，它们发送和接收数据，但是同时，可能有一个第三主的 [API]() 参与，另一个应用，或其他数据库。请记住这些，在将数据从一个应用程序发送到另一个中，在新的 app 接收它之前，我们需要转换数据格式。这个转换函数可以写在 *adapter.go* 文件中。

例如，如果我需要将结构体 `A` 转换为结构体 `B`，我需要一个适配函数，它看起来像：

```go
type A struct {
  FirstName string
  LastName string
  Email string
}

type B struct {
  Name string
  Email string
}

func ConvertAToB(obj A) B {
  return B{
    Name: obj.FirstName + obj.LastName,
    Email: obj.Email,
  }
```

## /api/auth

大多数web服务必须有至少一个认证方法被实现，像：

* [OAuth](https://en.wikipedia.org/wiki/OAuth) - 开放认证(Open Authentication)
* 基本认证
* Token 认证 ( 我喜欢JWT - [JSON Web Token](https://jwt.io/) )
* OpenID 

我个人使用 [JWT](https://jwt.io/)，因为我为我们的客户端([ATNM](https://www.airtouchmedia.com/))写web服务，大部分明为手机app 或 [CMS](https://en.wikipedia.org/wiki/Content_management_system) ，如果你想了解更多关于 [Web 认证 API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API)，Mozilla 有一篇不错的文件，它解释得非常好。

![KeyAU_47qONLdkwW](https://miro.medium.com/max/1500/0*KeyAU_47qONLdkwW.png)

什么是 JWT?

> JSON Web Token 是一种开放的行业标准 **[RFC 7519](https://tools.ietf.org/html/rfc7519)** 方法，用于在双方之间安全地表明应用。

你为什么要使用 JWT?

> **Authorization**：*这是使用 JWT 的一个非常普遍的场景。当用户登录，它的每一个子请求将包含JWT，允许用户来访问路由，服务以及该 token 允许访问的资源。单点登录(Single Sign On) 是现在广泛使用 JWT 的一个特点，因为它的负载比较小，并且它有轻松跨域的能力。*
>
> **信息交互(Information Exchange)**：JSON Web Token 是在不同部分之间安全地交换数据的一个很好的方式。因为 JWT 可以被签名 -- 例如，使用公、私钥对 -- 你可以确定发送者的身份。另外，签名是使用头部(header)和负荷(payload)计算出来的，你也可以认证内容没有被篡改。

所以，你必须验证签名，通过编码和反编码消息体，或者构建 JWT 体。对于这类的处理，我创建一个文件 *jwt.helper.go*，来保持一致性，并在包 *auth* 的一个位置找出所有与 JWT 相关的代码。

我们来讨论一个文件，*auth* 包中的 *principal.middleware.go*。这个文件之所以叫这个名字，是因为它是与任何API交互的第一个中间件，所以所有的请求都要通过它。在这个文件中，我写了一个函数用于阻止任何未通过规则的请求，会发送一个 **401状态码** 的响应。现在，如果想知道这些规则，我们已经谈论过 JWT，所以附加到任何请求(除了登录，注册这些不需要authorizatoin外)，客户端必须发送它的[HTTP header](https://en.wikipedia.org/wiki/List_of_HTTP_header_fields)， **Authorization**，必须包含 JWT token。

总结一下，如果客户端 app 不发送一个 token，或 token 损坏，或不合法，web 服务将阻止这个请求。

从哪儿获得 token?

在阅读之前段落时，这可能是你考虑的问题。我提到过，登录和注册不需要发送 token，因为你实际上是从这些请求中获得 token。所以填写上你的认证(Credentials)，如果它们是正确的，你将在响应中获得一个 token，登录时，将为每个它的请求附加上。

## /cmd

我总是喜欢将 *main.go* 放到这个包中，它包含项目的所有子包。它就像是一个封装所有子模块的一个包装器，使得它们可以一块工作。

**为什么叫这个名字？**，因为 *cmd* 比 command 简短。

**通过命令能理解什么？** 命令代表任务，它是某事的一部分，调用其他任务或独立支行。*main.go* 文件是一个命令，通常在一个文件中封装 web 服务的所有函数和包，并且只调用任何包的主函数。在任何时候，如果你想移除一个功能，你可以简单地通过在 main 文件中注释掉它的实例来移除它。

# 二、/config

我认为这个包是非常重要的，因为我发现将所有配置放到一个位置，而不是在项目的所有文件中搜索是非常有用的。在这个包中，我通常写一个名为 *config.go* 的文件，它包含 配置的模型。这个模型只是一个[结构体](https://gobyexample.com/structs)，例如：

```go
type JWT struct {
 Secret string `required:"true"`
}
type Database struct {
 Dialect  string `default:"postgres"`
 Debug    bool   `default:"false"`
 Username string `required:"true"`
 Password string `required:"true"`
 Host     string `required:"true"`
 Port     int
 SSLMode bool
}
type Configuration struct {
 Database Database `required:"true"`
 JWT      JWT      `required:"true"`
}
```

这只是一个结构体定义，我们需要将实际的数据放置到相应的位置。对于这部分，我喜欢用多个 [JSON文件](https://json.org/)，根据环境变量，它们的命名为类似 `config.ENV.json` 。对于之前的结构体定义，一个伪 JSON 示例如下：

```json
{
    "Database": {
        "Dialect": "postgres",
        "Debug": true,
        "Username": "postgres",
        "Password": "pass",
        "Host": "example.com",
        "Port": 5432,
        "SSLMode": true,
    },
    "JWT": {
        "Secret": "abcdefghijklmnopqrstuvwxyz"
    }
}
```

我们来讨论一下 business，因为这部分对我来说非常特别，投入时间寻找最佳答案非常重要。对于你来说可能不是问题，但我遇到了一些问题，试图以一种好的方式导入配置。有很多可能性，但我必须面对两者之间选择的困境：

* 将配置对象作为变量从 main.go 中传递给我需要使用的最终函数。这当然是一个好主意，因为我仅为那些需要它的实例传递变量，所以这种方式我不会有速度质量的妥协。但这对于开发和重构来说非常耗时，因为我需要将配置一直从一个到另一个函数，所以，结果，你可能会疯掉。
* 定义一个全局变量，在需要的地方使用这个实例。但在我看来，这并不是最好的选择，因为我必须要声明一个变量，比如，在 main.go 文件中，然后在 main 函数中我需要 `Unmarshal()` JSON 文件，将它的内容放到全局声明的变量对象中。但是，可能我在该变量初始化就绪前调用它，这样我将获得一个空的对象，没有任何实际值，这种情况下我的 app 可能崩溃掉。
* 直接在需要的地方注入配置对象，这是我的最佳选择，它非常适合我。在 *config.go* 文件的末尾几千 ，我声明了以下几行：

```go
var Main = (func() Configuration {
var conf Configuration
 if err := configor.Load(&conf, "PATH_TO_CONFIG_FILE"); err != nil {
  panic(err.Error())
 }
 return conf
})()
```

关于这个实现，你需要知道的是我使用了一个叫做 [Configor](https://github.com/jinzhu/configor) 的库，它 unmarshal 一个文件，在这个例子中是个 JSON，把它加载到变量 `conf` 中并返回。

任何时候，当你需要使用配置文件的内容时，只需要键入包名，即 `config`，并调用变量 `Main` ，正如下面检索数据库配置的例子所示：

```
var myDBConf = config.Main.Database
```

**！！！注意：** 正如你看到的，你必须在插入配置文件的路径，但是由于你需要为不同的环境使用不同的文件，或许你可以设置一个叫 `CONFIG_PATH` 的环境变量。将其定义为一个环境变量，或在运行之前设置它：

```
$ CONFIG_PATH=/home/username/.../config.local.json go run cmd/main.go
```

使用 `os.Getenv("CONFIG_PATH")` 替换 `PATH_TO_CONFIG_FILE`。这样，它就不在你配置文件的路径位置，你就可以跳过一些操作系统错误。

# 三、/db

**db** 包你 web 服务最重要的包之一，您真的必须花费大量的时间来思考架构并开发该包，因为这是网络服务的目的之一 -- 收集和存储数据。在下面几行中，我将介绍我自己的版本，该版本大适合构建 Web 服务的大多数情况，请继续关注...

在深入这个文件结构之前，我必须向你坦白，我更喜欢使用 [ORM](https://en.wikipedia.org/wiki/Object-relational_mapping)，与使用 SQL 查询并将数据转换为数组，并尝试调试简单的查询相比，使用对象更容易并且提供了一种很好的方法来处理对象。我使用 [GORM](https://github.com/jinzhu/gorm)，因为它满足了我所有的需求：拥有所有基本的 ORM 函数(Find, Update, Delete,等等...)，接受联合（Has One, Has Many, 属于，多对多，多态），接受事务，可以构建 sql，具有自动迁移和其他很酷的特性。

## /db.go

我用些文件来保持  GORM 所有重要的配置。所以在这个文件中，我定义一个函数，它将数据库连接作为对象返回，这个函数在 *main.go* 中被调用，传递给所有需要与数据库交互的 API。

```go
// NewDatabase returns a new Database Client Connection
func NewDatabase(config *config.Database) (*gorm.DB, error) {
 db, err := gorm.Open(config.Dialect, config.Source)
 if err != nil {
  return nil, err
 }
if err := InitDatabase(db, config); err != nil {
  return nil, err
 }
return db, nil
}
```

你可以在 NewDatabase 中看到，在第 8 行，它调用函数 `InitDatabase()`，这个函数定义了我的 ORM 的行为并操作 AutoMigration。

```go
// InitDatabase initializes the database
func InitDatabase(db *gorm.DB, config *config.Database) error {
 db.LogMode(config.Debug)
// auto migrate
 models := []interface{}{
  &models.Account{},
  &models.PersonalInfo{},
  &models.Category{},
  &models.Subcategory{},
 }
 if err := db.AutoMigrate(models...).Error; err != nil {
  return err
 }
// Personal info
 if err := db.Model(&models.PersonalInfo{}).AddForeignKey("account_id", fmt.Sprintf("%s(id)", models.AccountTableName), "CASCADE", "CASCADE").Error; err != nil {
  return err
 }
// Subcategories
 if err := db.Model(&models.Subcategory{}).AddForeignKey("category_id", fmt.Sprintf("%s(id)", models.CategoryTableName), "CASCADE", "CASCADE").Error; err != nil {
  return err
 }
return nil
}
```

自动迁移(Auto Migration)验证表是否存在，如果不存在或模型不同，则会尝试进行同步。

除了自动迁移，我手动设置了外键，或如果需要，索引或其他 sql 约束。

Main.go 中一个简单的实例化示例看起来类似：

```go
// setup the database dbc
if err := db.NewDatabase(&config.Config.Database); err != nil {
  panic(err.Error())
}
```

## /service.go

这个文件的目的是为所有的处理程序保留一个结构，而不是在多个位置导入处理程序，或保持不一致，以便将一个对象从 *main.go* 传递给所有的 API 处理程序，该 API 处理程序包含一个到所有数据库处理程序的一个引用。所以这个文件看起来像：

```go
package db
import (
 "github.com/jinzhu/gorm"
 "PROJECT_FOLDER/db/handlers"
)
type Service struct {
 Account  *handlers.AccountHandler
 Category *handlers.CategoryHandler
}
func NewService(db *gorm.DB) Service {
 return Service{
  Account:  handlers.NewAccountHandler(db),
  Category: handlers.NewCategoryHandler(db),
 }
}
```

可以看到，我只有两个处理程序，并不包含 PersonalInfo 和 Subcategory，因为这些并不需要，所以从主处理程序中分开。例如，在不知道配置给它的账户时，你不需要个人信息，所以这两个将被封装在一个对象中。

这在 *main.go* 中 可以简单地被调用，就下面这样：

```go
// create the database service 
dbService := db.NewService(dbc）
```

## /db/models

> [Models](http://gorm.io/docs/models.html) 通常只是普通的Golang结构，基本的Go类型或它们的指针。

我在自动迁移函数中放了四个模型：Account, PersonalInfo, Category 和 Subcategory。我喜欢将每个模型定义不同的文件中，选择一个直观的名字，比如 *account.go*, *personalInfo.go*, *category.go* 和 *subcategory.go*。

作为一个例子，account.go 的内容如下：

```go
package models
import (
 "crypto/md5"
 "encoding/hex"
 "PROJECT_FOLDER/utils"
"github.com/jinzhu/gorm"
 "golang.org/x/crypto/bcrypt"
)
type Role string
const (
 RoleAdmin Role = "admin"
 RoleUser  Role = "user"
)
const (
 AccountTableName = "accounts"
)
type Account struct {
 gorm.Model
Username string `gorm:"type:varchar(100); unique; not null"`
 Password string `gorm:"not null"`
 Role     Role   `gorm:"type:varchar(5); not null"`
 Active   bool   `gorm:"not null"`
 Token    string `gorm:"not null"`
IsSocial     bool `gorm:"not null"`
 Provider     string
 PersonalInfo PersonalInfo
}
func (account *Account) BeforeCreate() error {
 password, err := HashPassword(account.Password)
 if err != nil {
  return err
 }
 account.Password = *password
 account.Token = GenerateToken()
return nil
}
func HashPassword(password string) (*string, error) {
 hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
 if err != nil {
  return nil, err
 }
pass := string(hash)
 return &pass, nil
}
func GenerateToken() string {
 hasher := md5.New()
    
 // you can check the utils.RandStr() in /utils chapter of this article
 hasher.Write([]byte(utils.RandStr(32)))
return hex.EncodeToString(hasher.Sum(nil))
}
```

## /db/handlers

数据库处理函数代表样板代码，它在多个位置重复，所以与其调用 GORM 函数，调用一个在 API 处理程序中使用的准备响应函数更好。

```go
package handlers
import (
 "PROJECT_FOLDER/db/models"
"fmt"
"github.com/jinzhu/gorm"
)
type AccountHandler struct {
 db *gorm.DB
}
func NewAccountHandler(db *gorm.DB) *AccountHandler {
 return &AccountHandler{
  db,
 }
}
func (h *AccountHandler) Find(accountID uint) (*models.Account, error) {
 var res models.Account
if err := h.db.Find(&res, "id = ?", accountID).Error; err != nil {
  return nil, err
 }
return &res, nil
}
func (h *AccountHandler) FindBy(cond *models.Account) (*models.Account, error) {
 var res models.Account
if err := h.db.Find(&res, cond).Error; err != nil {
  return nil, err
 }
return &res, nil
}
func (h *AccountHandler) Create(account *models.Account) error {
 return h.db.Create(account).Error
}
func (h *AccountHandler) UpdateProfile(profile *models.PersonalInfo, accountID uint) error {
 var personalInfo models.PersonalInfo
 cond := &models.PersonalInfo{
  AccountID: accountID,
 }
// only create it if it doesn't already exist
 if h.db.First(&personalInfo, cond).RecordNotFound() {
  profile.AccountID = accountID
  return h.db.Create(profile).Error
 }
return h.db.Model(models.PersonalInfo{}).Where(cond).Update(profile).Error
}
```

如果你想要一个 SQL 包处理程序的示例，我推荐 [XZYA](https://github.com/Xzya) 的 [SQLHandler](https://github.com/Xzya/sqlhandler)。

## /db/tests

我知道你的项目管理者可以跳过这一步，因为项目必须尽快地准备好，但是相信我，最好写测试样例。所以在我看来， Go 中编写测试单元的一个不错的包是 [Testify](https://github.com/stretchr/testify)， 它使用起来非常简单并且功能强大。

## /gen

Gen 文件夹是放置所有第三方主加生成的代码的文件夹。保持所有代码在一个位置非常简单，因为在生成新版本，你可能需要时时的清理它，你可以使用 [Makefile](https://en.wikipedia.org/wiki/Makefile) 简单地处理这些任务，我们将在后面讨论它。

工作中，我们通常使用有 [Swagger](https://swagger.io/)，一个使我们更轻松并帮助我们维护一个文件的工具，它可以作为 [API声明](https://editor.swagger.io/?_ga=2.19720637.408907328.1539089399-1459746301.1538552745)，[代码生成](https://github.com/swagger-api/swagger-codegen) 和 [文档](https://rebilly.github.io/ReDoc/) 的基础。因为 Swagger 非常的酷，而且我们只是人类，因此使用图形化界面比编写 [YAML](https://en.wikipedia.org/wiki/YAML) 和 JSON 规范文件要简单得多。对于这类的工作，我使用 [StopLight](https://stoplight.io/)，它的工作基本是提供一个图形化界面并在工作完成之后，导入你期望的规范文件。

## /locals

在大多数时候，翻译工作都是在客户端 app 中实现的，但是有时候你需要发送一些自定义错误或一些翻译后的电子邮件模板，你将会面对一个问题。在我的 Go 开发者的职业生涯中，我遇到了这样的困境，我选择了一个非常简单，非常基本的构建一些东西，或浏览可能适合我的现有的包。我发现了一个简单而又是我所需要的包。我需要一个可以读取 JSON 文件的简单的包，因为客户端 app 已经拥有了这些翻译，这样我不需要创建额外的东西，所以我发现了 [GoTrans](https://github.com/bykovme/gotrans)。这个包最酷的地方是，例如，你可以在 *cmd/main.go* 中声明它，然后可以在项目中的任何地方调用翻译函数。

### 如何初始化 Gotrans?

```go
// initialize locals
if err := gotrans.InitLocals("locals"); err != nil {
  panic(err)
}
```

### 如何使用有 Gotrans?

JSON 文件如下

```json
{
  "hello_world": "Hello World",
  "find_more": "Find more information about the project on the website %s"
}
```

> JSON 文件名应该使用浏览器支持的标准语言码或语言国家码。localization 文件夹里至少有一个 “en.json" 文件。

```go
gotrans.Tr("fr", "hello_world")
```

## 四、/public

你可能会问 *什么？！web服务 中会有一个 public 文件夹？！* 不错，或许并不是每次都需要它，但我正在尽可能地解释web服务的一般架构，有时你需要有类似 *条款和条件* 页或 *隐私策略* 或 HTML 邮箱模板或任何可以公开并作为公共 API 可以导出的资源。

# 五、/utils

构建一个大项目，有时需要额外的工具或者说是辅助解决小问题。但是这些[辅助](https://boobo94.xyz/ai/5-questions-build-custom-alexa-skill/)只j 一小段代码，所以不需要单独为这些简单和小段代码创建包。是的，**utils**包可以解决这个问题，因为你可以在此处将不同的代码放入单独的不同的文件中，例如：

```go
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandStr(n int) string {
 b := make([]rune, n)
 for i := range b {
  b[i] = letters[rand.Intn(len(letters))]
 }
 return string(b)
}
```

* 生成哈希密码
* 创建处理程序上传到云
* 创建处理程序发送 email
* 日志管理
* 等等...

基本上这里是存放所有无法分类的混乱的地方，但是这不是一团糟，它可以让你更轻松并帮助你节省时间，不需要在多个地方重复写同样的代码。

# 六、/vendor

这个文件夹是唯一一个你不需要修改任何东西的地方，所有的外部依赖或工程中导入的包都被下载并存储在这个地方，以便你的构建工作正常。这在 `build` 或 `run` 任务中自动被创建，因为在工程被编译前，它验证是否所有的导入都在 vendor 文件夹中。

## 如何下载包？

这时有多个正确的答案，我不想争论，但是我可以告诉默认的是 `go get PACKAGE` ，它将 [依赖](https://stackoverflow.com/questions/37237036/how-should-i-use-vendor-in-go-1-6) 放到 `$GOPATH/src` 或 `go install PACKAGE` 将二进制放到 `$GOPATH/bin` ，将包放到 `$GOPATH/pkg` 中。

## 如何管理包

使用管理依赖工具，你可以完成基本的任务，并且可以节省一些时间

我喜欢 **[DEP](https://github.com/golang/dep)** ，Golang 默认的一个。

在 MAC 中，它可以简单地使用 brew 安装

```shell
$ brew install dep
$ brew upgrade dep
```

或使用 CURL

```shell
$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

# 七、Makefile

我使用 `Makefile`，因为它简单并且可以自动化一些我不得不经常重复的一些任务，和一些在构建之前，我必须先执行的一此步骤，并且我需要在几个月甚至几年后执行些过程。可能我需要花费一些时间弄清楚该如何做。但是与花费所有的时间再次发现应该如何构建项目比起来，我可以趁热做这些，然后只需要查看 `makefile` 并选择我需要运行的任务。

```makefile
IP = "XXX.XXX.XXX.XXX"
PEM_FILE = "...PATH_TO_PEM/key.pem"

# Run the server
run:
	CONFIG=config/config.local.json go run cmd/main.go --port 8000
	
# Remove generated files under gen/ folder
clean:
	rm -r gen

serve:
	realize start
	
#Build
build:
	go build cmd/main.go
	
# Build for linux
build-linux:
	env GOOS=linux go build cmd/main.go
	
# Deploy just the code to dev
deploy-dev-code: build-linux copy-project
# Full deploy to dev
# With tests and swagger generation
deploy-dev: gen deploy-dev-code
# Copy the project build and dependencies to server
copy-project:
 ssh -i $(PEM_FILE) ubuntu@$(IP) 'sudo service api stop'
 scp -i $(PEM_FILE) -r locales/ ubuntu@$(IP):/home/ubuntu/project_name
 scp -i $(PEM_FILE) main ubuntu@$(IP):/home/ubuntu/project_name
 ssh -i $(PEM_FILE) ubuntu@$(IP) 'sudo service api start'
 rm main
```

你可以从 [GUN.org](https://www.gnu.org/) 找到关于 [makefile](https://www.gnu.org/software/make/manual/make.html#toc-Overview-of-make) 和如何使用它的精彩文章。

在这个文章中，你已经学习了关于 API 和如何构建架构，从 web 服务中如何与数据库交互，如何制作配置文件，在[客户端和服务器](https://www.boobo94.xyz/tutorials/setup-apache-host-as-proxy/) 如何使用 JWT 关注安全和权限，如何使用其他包使得工作更轻松，在结尾时学习了如何使用 `Makefile` 运行多个任务。























## 