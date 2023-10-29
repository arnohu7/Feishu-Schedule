# ECNC 自动化排班

## 介绍

原本打算前端部分写一个自动化插件，碍于时间不足，前端部分暂时不写。等本学期班表排好后再完善。

后端部分就是自动排班表的核心内容，该部分会自动读取已经排好的班表，并使用贪心算法自动化排班，之后将排班结果输出到另一张数据表中。

## 如何运行自动化排班程序

首先你的电脑要装 go 语言的开发环境，请在[Go语言下载](https://go.dev/dl/)下载对应平台的 go 语言，并将其下的 bin 文件夹添加至系统变量。

接着设置 GOPROXY：

```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

安装 `go` 必要插件：

```
go install -v github.com/rogpeppe/godef@latest
go install -v golang.org/x/tools/cmd/goimports@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v golang.org/x/tools/gopls@latest
```

之后再在当前文件夹下创建一个 `config.json` 文件，格式如下：

```
{
    "APP_ID": "",
    "APP_Secret": "",
    "APP_Token": "",
    "Access_Token": "",
    "Read_Table_ID": "",
    "Write_Table_ID": ""
}
```

上述各个字段的说明如下：

- `APP_ID` 和 `APP_Secret` 可以在[自动化排班后台](https://open.feishu.cn/app/cli_a5b7aac6333a5013/baseinfo)后台中找到。
- `APP_Token` 比较麻烦，打开[获取知识空间节点信息](https://open.feishu.cn/document/server-docs/docs/wiki-v2/space-node/get_node?appId=cli_a5b7aac6333a5013)，里面有调试工具，输入 `token`（`token` 指的是要读取的多维表格的 `url` 的 `?table` 前的那一串），在输出的结果中的 `obj_token` 就是 `APP_Token`。
- `Access_Token` 在上面的调试工具就可以获取到，也就是请求头中的 `Authorization`，记得选 `user_access_token`。
- `Read_Table_ID` 指的是**提交**的班表的 `table_id`，在 `url` 中就可以获取到。
- `Write_Table_ID` 指的是**输出**的班表的 `table_id`，在 `url` 中就可以获取到。

配置好 `config.json` 之后，在命令行中运行以下指令即可看到班表的输出结果：

```
go run main.go
```
