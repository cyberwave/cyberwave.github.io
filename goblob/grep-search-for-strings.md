[原文] [How to use grep to search for strings in files on the shell](https://www.howtoforge.com/tutorial/linux-grep-command/)

# 1 GREP 命令 - 概览

grep 命令，意思是 **全局正则表达式输出(global regular expression print)**，仍然是Linux终端环境中功能最丰富的命令之一。它恰好是一个非常强大的程序，可以使用户能够根据对复杂的规则对输入进行排序，因此，它成为众多命令链中相当流行的一环。`grep` 命令主要用于搜索文本或在任何给定文件中搜索包含与所提供的单词/字符串匹配的行。`grep` 默认打印匹配的行，并且它可用于搜索与一个/多个正则表达式匹配的文本行，而且不会引起任何干扰，并且仅输出匹配的行。 

# 2 基本的 grep 命令语法

grep 命令的基本语法如下：

```shell
grep 'word' filename
grep 'word' file1 file2 file3
grep 'string1 string2' filename
cat otherfile | grep 'something'
command | grep 'something'
command option1 | grep 'data'
grep --color 'data' fileName
```

# 3 如何使用 grep 命令搜索文件

在第一个例子中，在 Linux 密码文件中搜索用户 “tom"：

```shell
grep tom /etc/passwd
```

下面是简单的输出：

```shell
tom:x:1000:1000:tom,,,:/home/tom:/bin/bash
```

可以使用选项来让 `grep` 忽略大小写，例如，使用 `-i` 选项匹配 abc, Abc, ABC 和所有可能的组合，输出如下：

```shell
grep -i "tom" /etc/passwd
```

# 4 递归使用 grep

如果目录层次中有一堆文本文件，例如，在 Apache 配置的文件 `/etc/apache2/` 中，你想要找到定义了指定文本的文件，使用 `grep` 命令的 `-r` 选项进行递归搜索。这将在目录 `/etc/apache2/` 及其所有子目录中为字符串“ 197.167.2.9”（如下所示）执行递归搜索操作：

```shell
grep -r "mydomain.com" /etc/apache2/
```

另外，还可以使用如下的命令：

```shell
grep -R "mydomain.com" /etc/apache2/
```

以下是在Nginx服务器上进行类似搜索的示例输出：

```shell
grep -r "mydomain.com" /etc/nginx/
/etc/nginx/sites-available/mydomain.com.vhost:        if ($http_host != "www.mydomain.com") {
```

这里，你会发现在该文件所在的文件名( 例如 /etc/nginx/sites-available/mydomain.vhost )的不同行上看到 mydomain.com 的结果。使用 `-h` 选项可以很容易地抑制在输出数据中包含文件名(如下所述)：`grep -h -R "mydomain.com" /etc/nginx/`。下面是样例输出

```shell
grep -r "mydomain.com" /etc/nginx/
if ($http_host != "www.mydomain.com") {
```

# 5 使用 grep 仅搜索单词

当搜索 abc，grep 会匹配所有包含 abc 的内容，即，kbcabc, abc123和其他其他的组合，不遵守单词边界。可以强制 `grep` 命令选择仅包含构成整个单词（这仅匹配 abc 单词）的那些，如下所示：

```shell
grep -w "abc" file
```

# 6 使用 `grep` 搜索两个不同的单词

为了搜索两个不同的单词，必须使用 `egrep` 命令，如下所示：

```shell
egrep -w 'word1 | word2' /path/to/file
```

# 7 统计匹配单词的行数

`grep` 命令可以使用 `-c` (count)选项报告每个文件中与特定模式匹配的次数：

```shell
grep -c 'word' /path/to/file
```

另外，用户可以使用 `-n` 选项在每个输出行之前获取该行在文件中的行数：

```shell
grep -n 'root' /etc/passwd
```

示例输出如下：

```shell
1:root:x:0:0:root:/root:/bin/bash
```

# 8 Grep 反向匹配

用户可以使用 `-v` 选项打印反向匹配，意思是它将仅匹配那些不包含给定单词的行。例如，使用下面的命令打印所有不包含单词 par 的行：

```shell
grep -v par /path/to/file
```

# 9 仅列出匹配文件的名字

你必须使用 `-l` 选项列出包含特定单词的文件名，例如单词 'primary'，使用下面的命令：

```shell
grep -l 'primary' *.c
```

最后，可以使用下面的命令强制 grep 以特定的颜色显示输出

```shell
grep --color  root /etc/passwd
```

# 10 如何使用 grep 命令处理多个搜索模式

有这种情况，你可能希望在给定的文件(或多个文件)中搜索多个模式。在这样的场景中，你应该使用 `grep` 提供的 `-e` 命令行选项。

例如，假设你想要在当前目录中的 所有文本文件中搜索单词 "how"，"to"，和 "forge"，你可以这样做：

```shell
grep -e how -e to -e forge *.txt
```

命令行选项 `-e` 在模式是以连字符号(-)开头的场景中也很帮助。例如，如果你想搜索 ”-how"，下面的命令无济于事：

```shell
grep -how *.txt
```

这就是 `-e` 命令行选项出场的时候了，该命令可以了解您在这种情况下到底要尝试搜索的内容：

```shell
grep -e -how *.txt
```

# 11 如何将grep输出限制为特定的行数

如果您想将grep输出限制为特定的行数，你可以使用 `-m` 选项来做这些，例如，假如你想在 *testfile1.txt* 文件中搜索单词 "how"，它包含下列的行：

```txt
myname: $ cat testfile1.txt
how are you?
how is your friend?
how have you been?
how is everyone else?
myname: $
```

但要求 `grep` 在找到包含搜索模式的3行后停止搜索，可以运行以下命令来做这些：

```shell
grep "how" -m3 testfile1.txt
```

下面是命令的 man 手册页的内容：

> If the input is standard input from a regular file, and NUM matching lines are output, grep ensures that the standard input is positioned to just after the last matching line before exiting, regardless of the presence of trailing context lines. **This enables a calling process to resume a search**.

因此，假如你有一个具有循环的 bash 脚本，并且想在每次循环迭代中获取一个匹配项，那么使用 `grep -m` 就足够了。

# 12 如何使用 `grep` 从文件中获取模式

如果你想，你也可以使 `grep` 命令从文件中获取模式。工具的 `-f` 选项可以让你做这些。

例如，假如你想要在当前目录中搜索所有 .txt 文件中的 “how" 和 ”to" 单词，但是想要通过一个名为 “input" 的文件来提供这些输入字符串：

```shell
grep -f input *.txt
```

# 13 如何使grep仅显示与搜索模式完全匹配的行

到目前为止，我们已经看到默认情况下grep会匹配并显示包含搜索模式的所有的行。但是如果需求是使 `grep` 仅显示完全匹配搜索模式的那些行呢。这可以使用 `-x` 选项来完成。

假如 *testfile1.txt* 包含如下内容

```txt
how are you?
He asked me how are you?
how are you? And how is your friend?
```

你想要搜索的是 “how are you?"。所以为了确保 grep 仅显示完全匹配这个模式的行，以下面的方式使用它：

```shell
grep -X "how are you?" *.txt
```

# 14 如何强制grep在输出中不显示任何内容

可能会有这种情况，你不需要 grep 命令输出任何内容。相反，您只是想知道根据命令的退出状态，是否发现有匹配项。可以通过 `-q` 选项来实现。

当 `-q` 选项没有任何输出，工具的结束状态可以通过 `echo $?` 命令来确实。在 grep 的例子中，如果成功(意味着找到匹配项)，命令的结束状态是 `0`，当没有发现匹配项时，它的结束状态是 `1`。

下面的内容显示了成功和未成功的场景：

```txt
hostname: -$ grep -q "how are you?" *.txt
hostname: -$ echo $?
0
hostname: -$ grep -q "I am file" *.txt
hostname: -$ echo $?
1
hostname: -$ echo $?
```

# 15 如何使grep显示不包含搜索模式的文件的名称

默认地，命令 `grep` 显示包含搜索模式(命令行也一样)的文件名。这很符合逻辑，正如这个工具的预期一样。然而，可能有需要获得不包含搜索模式的那些文件的情况。

这也是 grep 可以做到的 -- 使用 `-L` 选项可以做这些。所以，例如，在当前的目录中找到所有不包含单词 "how" 的文本文件，你可以使用下列命令：

```shell
grep -L "how" *.txt
```

# 16 如何抑制 grep 产生的错误消息

如果需要，还可以强制grep使输出中不显示的所有错误消息。可以通过使用 `-s` 选项实现。例如，请考虑以下情况，其中grep会产生与其遇到的目录有关的错误/**告警**：

```shell
hostname: -$ grep "how" *
grep: examples: Is a directory
```

所以在这场景中，`-s` 选项会有帮助 。

# 17 如何使grep递归搜索目录

从上一点中的示例可以清楚地看出，默认情况下grep命令不会进行递归搜索。为了确保 grep 搜索是递归的，使用 `-d` 选项并将值 `recurse` 传递给它。

```shell
grep -d recurse "how" *
```

**注意1：** 上一节讨论的与目录相关的错误/告警消息也可以使用 `-d` 选项忽略掉 -- 你所需要做的就是将 `skip` 传递给它

**注意2：** 使用 `--exclude-dir=[DIR]` 选项从递归搜索中将与模式匹配的 DIR 排除

# 18 如何使grep终止以NULL字符为文件名

正如我们所讨论的，当你仅希望工具在输出中展示文件名时，使用grep 命令的  `-l` 选项：

```shell
grep -l "how" *.txt
```

现在，你需要知道的是上面输出中的每个名字都是以换行符分隔/附上的。你可以这样验证：重写向输出到文件，打印文件内容：

```txt
grep -l "how" *.txt > output.txt
cat output.txt
testfile1.txt
testfile2.txt
```

所以， `cat` 命令的输出证明了两个文件名中有换行符的存在。

但是正如你可能已经知道，换行符也可以是文件名的一部分。所以当处理这类文件名包含换行，他们也被分隔/终止为新的一行，在 grep 输出中这将变得困难(特别是在脚本中访问这些输出)

如果分隔符/终止符不是换行符，那就更好了。好吧，您将很高兴知道 grep 提供了命令行选项 `-Z` ，以确保文件名后跟NULL字符而不是换行符。

所以，在我们的例子中，命令是：

```shell
grep -lZ "how" .txt
```

以下是您应该了解的相关命令行选项：

> **-z, --nul-data**
>
> Treat the input as a set of lines, each terminated by a zero byte (the ASCII NUL character) instead of a newline. Like the -Z or --null option, this option can be used with commands like sort -z to process arbitrary file names.

# 19 更多 GREP 命令的例子

在我们的第二个 GREP 命令教程中，你可以发现更多关于如何使用 Linux 命令的例子

* [How to perform pattern search in files using Grep](https://www.howtoforge.com/tutorial/how-to-perform-pattern-search-in-files-using-grep/)
