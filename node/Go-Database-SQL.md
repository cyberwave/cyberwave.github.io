# Go Database/SQL

[原文](http://mindbowser.com/golang-go-database-sql/)

在我们Golang教程的[第8章](http://mindbowser.com/golang-common-utilities/), 我们讨论了 "Common Utilities in Project Golang"。 在本章节，我们来探索 ’Go-database/SQL‘。

为了在Go中使用 `SQL` 或 `类似SQL`的数据库。我们需要使用`"database/sql"`包。它提供连接数据库的轻量级接口。

基本上，为了在Go中访问数据库，我们需要使用 `sql.DB` 。你需要使用它来创建语句，事务，执行查询和获取结果。但是要记住，`sql.DB` 不是一个数据库连接。根据Go规范，“它是一个接口的抽象和数据库的存在，它可能与本地文件一样多样化，通过网络连接访问，或者在内存中和进程中”。

`sql.DB`执行以下任务：

1. 打开和关闭与底层数据库驱动的连接。
2. 管理连接池。

连接池像这样被管理：当你要做某些操作的时候，该连接将被标记为使用中，如果不再使用，它被返回到池中。这样的后果是如果你释放连接到池中失败，将导致 `db.SQL` 打开过多的连接进而耗尽资源。

在创建 `sql.DB` 后，你可以使用它查询数据库，以及创建语句和创建交易。



- 导入数据库驱动

要使用 `database/sql` ，你将需要包本身以及与数据库相关的特定驱动程序。您通常不应该直接使用驱动程序包，尽管有些驱动程序会鼓励您这样做。相反，如果可能，您的代码应仅引用 `database/sql` 中定义的类型。这将有且于避免您的代码依赖于驱动程序，以便您可以使用最少的代码修改来更换底层驱动程序(以及您正在访问的数据库)。

在这里，我们将使用 @julienschmidt和@arnehormann提供的优秀 [MySQL](https://github.com/go-sql-driver/mysql) 驱动程序。

现在你需要像这样来导入访问 db 的包:

```golang
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
```

我们已经在这个第三方驱动前使用 `_` 标识，因此我们的代码中没有任何导出的可见名称。

要使用此第三方驱动程序，请使用 `go get` 命令从 GitHub 上下载。

```golang
go get "github.com/go-sql-driver/mysql"
```

现在，我们准备好访问数据库了。

由于我们已经导入了包，所以现在需要创建数据库对象 `sql.DB`，要创建 `sql.DB`，请使用 `sql.Open()` ，这将返回 `*sql.DB` 对象。

```golang
//main.go
package main
import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)
func main(){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
	if err != nil {
		fmt.Println(err)
     }else{
         fmt.Println("Connection Established")
     }
     defer db.Close()
}
```

现在，需要澄清一下：

1. `sql.Open()` 中的第一个参数是数据库的驱动程序名称。此字符串使用 `database/sql` 注册，并且通常与包名称相同。还有其他驱动程序，如 sqlite3 github.com/mattn/go-sqlite3 和 postgres 包是github.com/lib/pg

1. 第二个参数是驱动指定的语法，它告诉驱动如何访问底层的数据存储。在这个例子中，我们连接位于我们本地数据为中的 `employeedb` 数据库。

   ![root:root@tcp](http://mindbowser.com/wp-content/uploads/2017/09/13.png

2. 你应该始终检查并处理来自 `database/SQL` 操作的错误。

如果 `sql.DB` 的生命周期不应超出函数的范围，`defer db.Close()` 是惯用的。

正如已经说过的一样，`sql.Open()` 并不建立任何到数据库的连接，也不验证驱动的参数。相反，它只是简单的数据库抽象。数据存储的第一个实现连接将在第一次实际需要的时候被创建。如果你想要检查数据库是否可以访问，使用 `db.Ping()` 并记得检查它的错误。

```golang
err := db.Ping()
if err != nil {
    // do something here
}
```

即使在完成数据库对象时必须使用 `Close()` 数据库对象， `sql.DB` 对象也可以使用它。不要频繁使用 `db.Open()` 和 `db.Close()` 。而是为您需要访问的每个不同数据存储创建一个 `sql.DB` 对象，并保留它直到程序完成访问该数据存储。根据需要传递它，或者以某种方式在 全局范围提供它，但保持打开状态。

现在，在连接打开之后，我们将看到从数据存储区检索结果集的操作。



从数据库中获取数据

Go 的 `database/sql` 函数名是有意义的。如果一个函数名包括 *Query* ,它意味要求数据库返回行集，即使它是空的。不返回行的语句不应使用Query函数；他们应该使用 `Exec()`。

现在让我们看看如何查询数据库，使用结果。我们将在用户的表中查询 id 为 11 的用户并打印其 id 和名称。我们将使用 `rows.Scan()` 将结果分配给变量，一次一行。

```golang
// main.go
package main
import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
)
 
func main(){
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
	if err != nil {
		log.Fatal(err)
	}else{
		fmt.Println("Connection Established")
	}
	var (
		id int
		name string
	)
	rows,err:=db.Query("select id, username from user where id = ?", 1)
	if err!=nil{
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		err:=rows.Scan(&id,&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
```

- 这里我们使用 `db.Query()` 发送一个查询到数据库
- 我们延迟调用 `rows.Close()`
- 使用 `rows.Next()` 在行集上进行迭代
- 使用 `rows.Scan()` 读取每一行的列到变量中
- 在迭代行集后检查错误



一些需要采取的预防措施

1. 您应该始终在 `rows.Next()` 循环结束时检查错误。
2. 其次，只要有一个打开的结果集（由行表示），底层连接就会很忙，不能用于其他查询。这意味着它在连接池中是不可用的。如果用 `rows.Next()` 迭代所有行，最终你将读取最后一行同时 `rows.Next()` 将返回内部EOF错误并为您调用 `rows.Close()`。但由于某种原因，你退出该循环 - 提前返回等等，然后行不会关闭，并且连接保持打开状态。这很容易导致资源耗尽。

1. `rows.Close()` 是一个无害的，如果它已经关闭将是个空操作，，所以你可以多次调用它。但是请注意，我们首先检查错误，仅在没有错误时，调用 `rows.Close()`，以避免运行时出现混乱(runtime panic)。

1. 你应该始终执行 `defer rows.Close()`，即使您还调用 `rows.Close()`。即使没有错误，为了避免运行时恐慌。
2. 不要在循环中延迟调用

**准备查询**

您应始终准备可被多次使用的查询。这些预准备语句具有在执行语句时将传递的参数。这比串联字符串（避免SQL注入攻击）要好得多。

在 MySQL 中，参数占位符是 `?` 。而在 postgresql 中占位符是 `$N` ,其中 `N` 是一个数字。SQLite接受其中任何一个。在Oracle中，占位符以冒号和名称开头，如parameter1。我们在这里用 `?` 对于MySQL

```golang
//main.go
package main
import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
    "log"
)
 
func main(){
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
    if err != nil {
    	log.Fatal(err)
    }else{
    	fmt.Println("Connection Established")
    }
    var (
        id int
        name string
    )

    rows,err:=db.Query("select id, username from user where id = ?", 2)
    if err!=nil{
   		log.Fatal(err)
    }

    for rows.Next(){
    	err:=rows.Scan(&id,&name)
        if err != nil {
       		log.Fatal(err)
        }
        fmt.Println(id, name)
    }
}

```

这里 `db.Query()` 准备，执行和关闭已经准备的语句 。

**单行查询**

如果一个查询返回至多一行，则可以使用一些冗长的样板代码的快捷方式。

样例代码：

```golang
var name string
err:=db.QueryRow(“select username from user where id=?”)
if err!=nil{
	log.Fatal(err)
}
fmt.Println(err)
 
//Errors from the query are deferred until Scan() is called, and then are returned from that. You can also call QueryRow() on a prepared statement:
 
stmt, err := db.Prepare("select username from user where id = ?")
if err != nil {
	log.Fatal(err)
}
var name string
err = stmt.QueryRow(1).Scan(&name)
if err != nil {
	log.Fatal(err)
}
fmt.Println(name)
```

**使用事务修改数据**

**现在我们可以看下如何使用事务来修改数据。**

**修改数据**

使用 `Exec()` 与 INSERT,UPDATE,DELETE 或其他不先回行的语句一直使用

```golang
//main.go
package main
import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
    "log"
)
 
func main(){
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
    if err != nil {
    	log.Fatal(err)
    }else{
   		fmt.Println("Connection Established")
    }
    stmt, err := db.Prepare("INSERT INTO user(id,username) VALUES(?,?)")
    if err != nil {
    	log.Fatal(err)
    }
    res, err := stmt.Exec(33,"Ray Martin")
    if err != nil {
    	log.Fatal(err)
    }
    lastId, err := res.LastInsertId()
    if err != nil {
    	log.Fatal(err)
    }
    rowCnt, err := res.RowsAffected()
    if err != nil {
    	log.Fatal(err)
    }
    log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
}
```

此后执行一个语句，它给出了 `sql.Result`，它提供了对语句元数据的访问：受影响的最后插入的id和no.行数。 

如果您不想返回结果，请查看以下案例。

```golang
_, err := db.Exec("DELETE FROM users") // GOOD APPROACH
_, err := db.Query("DELETE FROM users") // BAD APPROACH
```

如果我们在这里使用 `db.Query()` ，这两个并不一样。然后它将返回 `sql.Rows` 并且它将保持连接打开直到连接关闭。

**使用事务**

在Go, 一个事务是一个存储连接数据存储的对象。它确保所有关系到同一个连接的操作会被执行。

这里，你需要调用 `db.Begin()`来开始一个事务，使用返回的 tx 变量的 `Commit()` 或 `Rollback()` 方法关闭它。在这个情况下，tx 从池中获取一个连接并存储它仅在事务中使用。

在事务中创建的准备语句被绑定到同样的事务中。

这里关于事务要记住的重要的事是：

Tx 对象可以保持打开状态，保存连接池并不返回它。

在事务工作中你应该关注不要使用 db 变量调用，确保所有的调用仅用你使用 `db.Begin()` 创建的 tx 变量，因为 db 不是一个事务，tx  才是事务。如果你试图使用 db 变量调用，则这些调用将不会在事务中起作用。

```golang
main.go
package main
import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "log"
)

func main(){
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
    if err != nil {
        log.Fatal(err)
    }
    tx,_:=db.Begin()
    stmt, err := tx.Prepare("INSERT INTO user(id,username) VALUES(?,?)")
    res,err:=stmt.Exec(4,"Ricky")
    res,err=stmt.Exec(5,"Peter")
    if err!=nil{
        tx.Rollback()
        log.Fatal(err)
    }
    tx.Commit()
    log.Println(res) 
}
```



## 使用 Prepared 语句

### Prepared 语句和连接

一个 prepared 语句是 一个带有参数占位符的语句，该参数占位符被发往数据库服务器可以被重复执行。这是一种性能优化以及安全措施；它可以防止攻击，例如SQL注入，攻击者劫持无人看守的字符串连接以产生恶意查询。

在MySQL以及大多数数据库中，首先将SQL发送到服务器并要求使用占位符来准备绑定参数。 服务器以语句ID响应。 然后，您将执行命令发送到服务器，并向其传递语句ID和参数。

在数据库级别，预准备语句被绑定到一个单独的 db 连接。典型的流程就像客户端将带有占位符的SQL语句发送到服务器准备，服务器以语句ID响应，然后客户端通过发送ID和语句来执行语句。

在 Go 中，连接并不是直接的在 `database/sql` 包中暴露。你需要在 db 或 tx 对象上准备一个语句，并不是直接在数据库连接上。

1. 当一个语句已经准备好，它在连接池中准备好。

1. Stmt 对象会记住使用的是哪一个连接
2. 当你执行 Stmt，它尝试使用连接。如果连接不可用，它会从池中获取其他连接并使用 db 再次准备语句。

由于重新编写语句，因此db的并发使用率会很高，这可能会使连接忙碌。

```golang
//main.go
package main
import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "log"
)

func main(){
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
    if err != nil {
   	 log.Fatal(err)
    }

    stmt, err := db.Prepare("INSERT INTO user(id,username) VALUES(?,?)")
    res,err:=stmt.Exec(99,"John")
    res,err=stmt.Exec(88,"Martin")
    if err!=nil{
    	log.Fatal(err)
    } 
    log.Println(res) 
}
```



### 在事务中准备语句

我已经说过准备语句在事务中被创建并仅绑定到同一个事务上。所以，当我们在 tx 对象上工作，这意味着我们的操作正在处理一个且唯一的连接。

这意味在 tx 中被创建的准备语句不能从它上分开使用。同样在 db 上创建的预准备语句不可以在事务中使用，因为它们绑定的是不同的连接。

要在tx中使用在事务外部准备的预准备语句，可以使用tx.Stmt()，它将从事务外部准备的语句创建新的特定于事务的语句

它通过获取现有的预准备语句，设置与事务的连接并在每次执行时重新准备所有语句来完成此操作。

```golang
tx, err := db.Begin()
if err != nil {
	log.Fatal(err)
}
defer tx.Rollback()
stmt, err := tx.Prepare("INSERT INTO user VALUES (?)")
if err != nil {
	log.Fatal(err)
}
defer stmt.Close() // danger!
for i := 0; i < 10; i++ {
    _, err = stmt.Exec(i)
    if err != nil {
    	log.Fatal(err)
	}
}
err = tx.Commit()
if err != nil {
	log.Fatal(err)
}
// stmt.Close() runs here!
```

以下是参数占位符语法，它们是特定于数据库的。考虑比较MySQL，PostgreSQL，Oracle。

--|--|--

MySQL | PostgreSQL | Oracle

--|--|----

WHERE  col=? | WHERE col=$1 | WHERE col=:col

--|--|--

VALUES=(?,?,?) | VALUES($1,$2,$3) | VALUES(:val1,:val2,:val3)



## 错误处理

几乎所有database/SQL类型的操作都返回错误作为最后一个值。你应该总是检查错误，永远不要忽略它们。您可能知道一些特殊的错误行为

### 迭代结果集错误

请考虑如下代码：

```golang
for rows.Next() {
    //...
}
if err = rows.Err(); err != nil {
    // handle the error here
}
```

### 关闭结果集错误

如果过早地结束循环，你应该总是明确地关闭 `sql.Rows` ，如果循环正常结束或抛出异常它会自动关闭，但是你可能意外地执行此操作

```golang
for rows.Next(){
    //...
    break //whoops, rows is not closed! memory leak...
}
// do the usual "if err = rows.Err()" [omitted here]...
//it's always safe to [re?]close here:
if err = rows.Close();err != nil {
    //but what should we do if there's an error?
    log.Println(err)
}
```

`rows.Close` 返回的错误是一般规则的唯一例外，它最好捕获并检查所有数据库操作中的错误。如果`rows.Close`返回错误，则不清楚应该怎么做。

### QueryRow 错误

考虑下面获取单行的代码

```golang
var name string
err = db.QueryRow("select username from user where id=?",1).Scan(&name)
if err != nil {
    log.Fatal(err)
}
fmt.Println(name)
```

如果没有 id=1 的用户怎么办？那将在结果中没有行，同时 `.Scan()` 不会扫描到结果到name变量中。接下来会发生什么？Go定义特殊的错误常量 `sql.ErrNoRows` ,当结果为空时，它由 `QueryRow()` 返回。这需要被处理。应用程序代码不会将空结果视为错误，如果您未检查错误是否为特殊常量，则会导致应用程序代码错误。

查询中的错误将被推迟，直到调用 `Scan()` ，然后从中返回。

```golang
var name string
err = db.QueryRow("select username from user where id = ?",1).Scan(&name)
if err != nil {
    if err == sql.ErrNoRows {
        // there were no rows, but otherwise no error occurred
    } else {
        log.Fatal(err)
    }
}
fmt.Println(name)
```



### NULLS 

可空列可导致许多丑陋的代码。如果可以，请避免使用它们。如果没有，那么你需要使用 `database/sql`包中的特殊类型来处理它们或定义你自己的类型

这里是 booleans, strings, integers 和 floats 的空类型，下面是如何使用它们。

```golang
for rows.Next(){
    vart s sql.NullString
    err := rows.Scan(&s)
    //check err
    if s.Valid {
        // use s.String
    } else {
        // NULL value
    }
}
```

但是有一些限制和理由避免空列.

1. 没有 `sql.NullUint64` 或 `sql.NullYouFavouriteType`。你需要自己定义这些。

1. 可空性可能很棘手，而且不会出现面向未来的问题。如果你认为某些东西不会为空，但你错了，你的程序就会崩溃。
2. 关于Go的一个好处是每个变量都有一个有用的默认零值。这不是可以为空的东西工作的方式。



### 连接池

连接池由 `database/SQL` 包提供。连接池是维护一个连接池并重用这些连接的机制。它用于增强在数据库上执行命令的性能。它有助于重用相同的连接对象来服务于许多客户端请求。

每次收到客户端请求时，都会在池中搜索可用连接，并且很可能它会获得空闲的连接。否则，传入的请求将排队或创建新连接并添加到池中(具体取决于池中已存在的连接数)。一旦请求完成使用连接，就会将其返回到分配给它的池中.

有关连接池的一些有用信息：

1. 连接池意味着在单个数据库上执行两个连续语句，可能会打开两个连接并分别执行它们。例如，`LOCK TABLES` 后跟一个可阻塞的 `INSERT`，因为 `INSERT` 是一个不保持表级锁的连接。
2. 连接在需要时被创建，并且池中没有空闲的连接。
3. 默认的，连接数没有限制。如果你试图一次做很多事，你可以创建任意数字的连接。这会导致数据库返回一个类似 "too many connections" 的错误。
4. 在 Go 1.1 或更高版本，你可以使用 `db.SetMaxIdleConns(N)` 来 限制池中空闲数目的连接。但是，这并不限制池的大小。

1. 在 Go 1.2.1 或更高版本，你可以使用 `db.SetMaxOpenConns(N)` 来限制数据库总共打开连接的个数，不幸的是，死锁(已修复)会阻止 `db.SetMaxOpenConns(N)` 在 1.2 中安全地使用。

2. 连接回收的速度非常快。

3. 保持一个 连接常时间空闲会导致问题。如果连接超时，请尝试使用 `db.SetMaxIdleConns(0)`，因为连接空闲时间过长。

   ```golang
   //main.go
   package main
   import (
       _ "github.com/go-sql-driver/mysql"
       "database/sql"
       "log"
   )
   
   func main(){
       db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/employeedb")
   
       // Connection Pooling methods
   
       db.SetConnMaxLifetime(500)
       db.SetMaxIdleConns(50)
       db.SetMaxOpenConns(10)
       db.Stats()
   
       if err != nil {
       	log.Fatal(err)
       }
       tx,_:=db.Begin()
       stmt, err := tx.Prepare("INSERT INTO user(id,username) VALUES(?,?)")
       res,err:=stmt.Exec(4,"Abhijit")
       res,err=stmt.Exec(5,"Yogesh")
       if err!=nil{
      		tx.Rollback()
       	log.Fatal(err)
       }
       tx.Commit()
       log.Println(res)
   }
   ```

   以下是我们发现有用的一些外部信息来源:

   √ http://golang.org/pkg/database/sql/database

   √ http://jmoiron.net/blog/gos-database-sql/

   √ http://jmoiron.net/blog/built-in-interfaces/
