[只需要五步，为你的Golang应用程序增加GUI](https://medium.com/@social_57971/how-to-add-a-gui-to-your-golang-app-in-5-easy-steps-c25c99d4d8e0)
构建一个Golang应用程序是简单而有趣的。但是有时你想要锦上添花：GUI
在这个故事中，我将介绍如何使用 [astilectron](https://github.com/asticode/go-astilectron)和它的 [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) 以及 [bundler](https://github.com/asticode/go-astilectron-bundler) 来给一个简单的Golang添加一个GUI

我们的 Golang GUI 应用程序将探索一个文件夹并展示其内容的宝贵信息。

你可以在[这儿](https://github.com/asticode/go-astilectron-demo)找到最终的代码

![Astilectron demo](https://cdn-images-1.medium.com/max/800/0*rQRTObJN7IG1O5dd.png)

# 步骤1：组织项目
这个文件的结构如下：
```javascript
|--+ resources
   |--+ app
      |--+ static
         |--+ css
            |--+ base.css
         |--+ js
            |--+ index.js
         |--+ lib
            |--+ ... (all the css/js libs we need)
      |--+ index.html
   |--+ icon.icns
   |--+ icon.ico
   |--+ icon.png
|--+ bundler.json
|--+ main.go
|--+ message.go
```
正如您所看到的，我们需要3种不同格式的图标来实现跨平台格式化：
`darwin` 的 `.icns`，`windows` 的 `.ico` 和 `linux` 的 `.png`。

我们将使用下列的 CSS/JS 库：
* [astiloader](https://github.com/asticode/js-toolbox)
* [astimodaler](https://github.com/asticode/js-toolbox)
* [astinotifier](https://github.com/asticode/js-toolbox)
* [chartjs](http://www.chartjs.org/)
* [fontAwesome](http://fontawesome.io/)

# 步骤2：实现架构
## GO
第一我们需要在 `main.go` 中建立 [astilectron](https://github.com/asticode/go-astilectron)的 [bootstrap](https://github.com/asticode/go-astilectron-bootstrap):
```golang
package main

import (
	"flag"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
)

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:    *debug,
		Homepage: "index.html",
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("About")},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = iw
			return nil
		},
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(700),
			Width:           astilectron.PtrInt(700),
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
```
2 全局变量 `AppName` 和 `BuiltAt` 将使用 [bundler](https://github.com/asticode/go-astilectron-bundler)(通过ldflags)自动被填充。

你可以看到我们的主页为 `index.html`，我们将有一个漂亮的菜单，该菜单拥有2个项目(`about`和`close`)，我们的主窗口将 `700x700`,`centered` 并且有一个 `#333` 的背景。

我们也要根据GO 标志添加一个 `debug` 选项，以防我们想要使用 HTML/JS/CSS 开发工具。

最终我们保存一个指向 `astilectron.Window` 的全局变量指针 `w`，以防稍后使用选项　`OnWait` 使用它，`OnWait` 包含一个回调，该回调执行一次window，菜单和所有已经被创建的对象。

## HTML
现在我们需要创建我们的HTML主页，在 `resources/app/index.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="static/css/base.css"/>
    <link rel="stylesheet" href="static/lib/astiloader/astiloader.css">
    <link rel="stylesheet" href="static/lib/astimodaler/astimodaler.css">
    <link rel="stylesheet" href="static/lib/astinotifier/astinotifier.css">
    <link rel="stylesheet" href="static/lib/font-awesome-4.7.0/css/font-awesome.min.css">
</head>
<body>
    <div class="left" id="dirs"></div>
    <div class="right">
        <div class="title"><span id="path"></span></div>
        <div class="panel"><span class="stat" id="files_count"></span> file(s)</div>
        <div class="panel"><span class="stat" id="files_size"></span> of file(s)</div>
        <div class="panel" id="files_panel">
            <div class="chart_title">Files repartition</div>
            <div id="files"></div>
        </div>
    </div>
    <script src="static/js/index.js"></script>
    <script src="static/lib/astiloader/astiloader.js"></script>
    <script src="static/lib/astimodaler/astimodaler.js"></script>
    <script src="static/lib/astinotifier/astinotifier.js"></script>
    <script src="static/lib/chart/chart.min.js"></script>
    <script type="text/javascript">
        index.init();
    </script>
</body>
</html>
```
这儿没有什么神秘的：我们声明我们的 `css` 和 `js` 文件，我们安装 `html` 结构并且我们确保我们的 `js` 脚本是已经通过 `index.init()` 初始化的。

## CSS
现在我们需要创建我们CSS格式在文件 `resources/app/static/css/base.cs`:
```css
* {
    box-sizing:  border-box;
}

html, body {
    background-color: #fff;
    color: #333;
    height: 100%;
    margin: 0;
    width: 100%;
}

.left {
    background-color: #333;
    color: #fff;
    float: left;
    height: 100%;
    overflow: auto;
    padding: 15px;
    width: 40%;
}

.dir {
    cursor: pointer;
    padding: 3px;
}

.dir .fa {
    margin-right: 5px;
}

.right {
    float: right;
    height: 100%;
    overflow: auto;
    padding: 15px;
    width: 60%;
}

.title {
    font-size: 1.5em;
    text-align: center;
    word-wrap: break-word;
}

.panel {
    background-color: #f1f1f1;
    border: solid 1px #e1e1e1;
    border-radius: 4px;
    margin-top: 15px;
    padding: 15px;
    text-align: center;
}

.stat {
    font-weight: bold;
}

.chart_title {
    margin-bottom: 5px;
}
```

## JS
最后我们创建我们自己的JS脚本在 `resources/app/static/js/index.js`:
```javascript
let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
    }
};
```
该 `init` 方法适当地初始化库。

# 步骤3：设置GO和Javascript之间的通信
一切都开始落实到位，但我们仍然缺少一个关键组件：GO和Javascript之间的通信。

## 从Javascript到GO的通信
为了从Javascript到GO进行通信，我们首先需要从Javascript发送消息到GO，并在接收到响应后执行回调：
```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will send a message to GO
    astilectron.sendMessage({name: "event.name", payload: "hello"}, function(message) {
        console.log("received " + message.payload)
    });
})
```
同时，我们需要在GO中监听消息，并通过 `MessageHandler` 引导选项向Javascript发送可选消息：
```golang
func main() {
	bootstrap.Run(bootstrap.Options{
		MessageHandler: handleMessages,	
	})
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
		    payload = err.Error()
		    return
		}
		payload = s + " world"
	}
	return
}
```
这个简单的示例将会在Javascript输出中打印 `received hello world`。

在我们的例子中，我们会添加更多的逻辑，由于我们希望允许浏览文件夹并显示有关其内容的有价值的信息。
因此，我们添加以下内容到 `resources/app/static/js/index.js`:
```javascript
let index = {
    addFolder(name, path) {
        let div = document.createElement("div");
        div.className = "dir";
        div.onclick = function() { index.explore(path) };
        div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
        document.getElementById("dirs").appendChild(div)
    },
    init: function() {
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Explore default path
            index.explore();
        })
    },
    explore: function(path) {
        // Create message
        let message = {"name": "explore"};
        if (typeof path !== "undefined") {
            message.payload = path
        }

        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
            }

            // Process path
            document.getElementById("path").innerHTML = message.payload.path;

            // Process dirs
            document.getElementById("dirs").innerHTML = ""
            for (let i = 0; i < message.payload.dirs.length; i++) {
                index.addFolder(message.payload.dirs[i].name, message.payload.dirs[i].path);
            }

            // Process files
            document.getElementById("files_count").innerHTML = message.payload.files_count;
            document.getElementById("files_size").innerHTML = message.payload.files_size;
            document.getElementById("files").innerHTML = "";
            if (typeof message.payload.files !== "undefined") {
                document.getElementById("files_panel").style.display = "block";
                let canvas = document.createElement("canvas");
                document.getElementById("files").append(canvas);
                new Chart(canvas, message.payload.files);
            } else {
                document.getElementById("files_panel").style.display = "none";
            }
        })
    }
};
```
一旦Javascript `astilectron` 命名空间准备好，它将执行新的 `explore` 方法，它发送一个消息到 GO ，同时接收响应，根据响应更新HTML。

随后我们添加下列代码到　`message.go`:
```go
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/asticode/go-astichartjs"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "explore":
		// Unmarshal payload
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}

		// Explore
		if payload, err = explore(path); err != nil {
			payload = err.Error()
			return
		}
	}
	return
}

// Exploration represents the results of an exploration
type Exploration struct {
	Dirs       []Dir              `json:"dirs"`
	Files      *astichartjs.Chart `json:"files,omitempty"`
	FilesCount int                `json:"files_count"`
	FilesSize  string             `json:"files_size"`
	Path       string             `json:"path"`
}

// PayloadDir represents a dir payload
type Dir struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// explore explores a path.
// If path is empty, it explores the user's home directory
func explore(path string) (e Exploration, err error) {
	// If no path is provided, use the user's home dir
	if len(path) == 0 {
		var u *user.User
		if u, err = user.Current(); err != nil {
			return
		}
		path = u.HomeDir
	}

	// Read dir
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(path); err != nil {
		return
	}

	// Init exploration
	e = Exploration{
		Dirs: []Dir{},
		Path: path,
	}

	// Add previous dir
	if filepath.Dir(path) != path {
		e.Dirs = append(e.Dirs, Dir{
			Name: "..",
			Path: filepath.Dir(path),
		})
	}

	// Loop through files
	var sizes []int
	var sizesMap = make(map[int][]string)
	var filesSize int64
	for _, f := range files {
		if f.IsDir() {
			e.Dirs = append(e.Dirs, Dir{
				Name: f.Name(),
				Path: filepath.Join(path, f.Name()),
			})
		} else {
			var s = int(f.Size())
			sizes = append(sizes, s)
			sizesMap[s] = append(sizesMap[s], f.Name())
			e.FilesCount++
			filesSize += f.Size()
		}
	}

	// Prepare files size
	if filesSize < 1e3 {
		e.FilesSize = strconv.Itoa(int(filesSize)) + "b"
	} else if filesSize < 1e6 {
		e.FilesSize = strconv.FormatFloat(float64(filesSize)/float64(1024), 'f', 0, 64) + "kb"
	} else if filesSize < 1e9 {
		e.FilesSize = strconv.FormatFloat(float64(filesSize)/float64(1024*1024), 'f', 0, 64) + "Mb"
	} else {
		e.FilesSize = strconv.FormatFloat(float64(filesSize)/float64(1024*1024*1024), 'f', 0, 64) + "Gb"
	}

	// Prepare files chart
	sort.Ints(sizes)
	if len(sizes) > 0 {
		e.Files = &astichartjs.Chart{
			Data: &astichartjs.Data{Datasets: []astichartjs.Dataset{{
				BackgroundColor: []string{
					astichartjs.ChartBackgroundColorYellow,
					astichartjs.ChartBackgroundColorGreen,
					astichartjs.ChartBackgroundColorRed,
					astichartjs.ChartBackgroundColorBlue,
					astichartjs.ChartBackgroundColorPurple,
				},
				BorderColor: []string{
					astichartjs.ChartBorderColorYellow,
					astichartjs.ChartBorderColorGreen,
					astichartjs.ChartBorderColorRed,
					astichartjs.ChartBorderColorBlue,
					astichartjs.ChartBorderColorPurple,
				},
			}}},
			Type: astichartjs.ChartTypePie,
		}
		var sizeOther int
		for i := len(sizes) - 1; i >= 0; i-- {
			for _, l := range sizesMap[sizes[i]] {
				if len(e.Files.Data.Labels) < 4 {
					e.Files.Data.Datasets[0].Data = append(e.Files.Data.Datasets[0].Data, sizes[i])
					e.Files.Data.Labels = append(e.Files.Data.Labels, l)
				} else {
					sizeOther += sizes[i]
				}
			}
		}
		if sizeOther > 0 {
			e.Files.Data.Datasets[0].Data = append(e.Files.Data.Datasets[0].Data, sizeOther)
			e.Files.Data.Labels = append(e.Files.Data.Labels, "other")
		}
	}
	return
}
```
当接收正确的消息，它执行新的 `explore` 方法，它根据一个路径返回有价值的信息。
最后，如简化示例中所示，我们不要忘记添加适当的 `MessageHandler` [bootstrap](https://github.com/asticode/go-astilectron-bootstrap)选项。

从 GO 到 Javascript 的通信
为了从 GO 到 Javascript 通信，我们首先需要发送一个从 GO 到 Javascript 的消息，当接收到响应的时候执行一个回调：
```golang
// This will send a message and execute a callback
// Callbacks are optional
bootstrap.SendMessage(w, "event.name", "hello", func(m *bootstrap.MessageIn) {
    // Unmarshal payload
    var s string
    json.Unmarshal(m.Payload, &s)

    // Process message
    log.Infof("received %s", s)
})
```
同时，我们需要在 Javascript 中监听消息，返回一个选项消息到 GO:
```javascript

// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        if (message.name === "event.name") {
            return {payload: message.message + " world"};
        }
    });
})
```
这个简单的例子将在GO输出中打印 `received hello world`。
在我们的例子中首先添加下列代码到 `main.go`:
```javascript
func main() {
	bootstrap.Run(bootstrap.Options{
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astilectron.PtrStr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
								return
							}
							astilog.Infof("About modal has been displayed and payload is %s!", s)
						}); err != nil {
							astilog.Error(errors.Wrap(err, "sending about event failed"))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = iw
			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
			}()
			return nil
		},
	})
}
```
它使得 `about` 菜单列表可点击并显示一个正确的内容模式。在 GO 应用程序已经被初始化后5秒请求显示一个通知。
最后我们添加下列代码到 `resources/app/static/js/index.js`：
```javascript
let index = {
    about: function(html) {
        let c = document.createElement("div");
        c.innerHTML = html;
        asticode.modaler.setContent(c);
        asticode.modaler.show();
    },
    init: function() {
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
        })
    },
    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    }
};
```
它监听GO消息并作出相应的反应。

# 步骤4：绑定应用程序
现在代码已经到位，我们需要确保我们将我们的Golang GUI应用程序以最好的方式呈现给我们的用户：
* 为 `darwin` 用户提供的一个 MacOSX 应用程序
* 为 `windows` 提供一个有漂亮图标的 `.exe`
* 为 `linux` 用户提供一个简单的二进制
我们很幸运，我们有 [astilectron](https://github.com/asticode/go-astilectron) 的 [bundler](https://github.com/asticode/go-astilectron-bundler) 做这这些！
首先我们运行下列代码安装它：
```shell
$ go get -u github.com/asticode/go-astilectron-bundler/...
```
然后我们添加正确的 [bootstrap](https://github.com/asticode/go-astilectron-bootstrap)选项到 `main.go`：
```golang
func main() {
	bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		RestoreAssets:  RestoreAssets,
	})
}
```
接下来我们创建配置文件 `bundler.json`：
```golang
{
  "app_name": "Astilectron demo",
  "icon_path_darwin": "resources/icon.icns",
  "icon_path_linux": "resources/icon.png",
  "icon_path_windows": "resources/icon.ico",
  "output_path": "output"
}
```
最后我们在项目路径下运行下列命令(确保 `$GOPATH/bin` 在你的 `$PATH` 中)：
```shell
$ astilectron-bundler -v
```
# 步骤5：在实际中看效果
瞧！结果是在 `output/<os>-<arch>` 文件夹中，准备测试:)  
![Astilectron demo](https://cdn-images-1.medium.com/max/800/0*V2GBx5FPJVpw2SeZ.png)

你当然可以绑定你的 Golang GUI 应用程序到其他环境，检查 [bundler](https://github.com/asticode/go-astilectron-bundler) 文档看怎样获取这些。

# 结论
借助一些组织和结构，在Golang应用程序中添加GUI从未如此简单，这要归功于 [astilectron](https://github.com/asticode/go-astilectron)，它的 [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) 和它的[bundler](https://github.com/asticode/go-astilectron-bundler)。
注意到使用它的两个主要缺点是有趣的：
* 二进制文件的大小至少为50MB，第一次执行后，包含二进制文件的文件夹的大小至少为200MB
* 内存消耗可能会有点忙乱，因为Electron（在引擎盖下运行）已知不能很好地管理它

但是，如果您准备好了，那么您将很快为您的Golang应用程序添加一个GUI！

快乐的GUI编码！
