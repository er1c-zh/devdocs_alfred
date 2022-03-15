# devdocs_alfred

一个搜索devdocs的alfred workflow。

## 优势

依赖golang提供的多平台编译能力，
编译之后不需要特殊的执行环境。

## Usage

```shell
./export.sh
```

生成的`output`文件夹就是workflow，引用到alfred中即可。

`sdoc`关键字开始输入要搜索的文档，
回车选中目标文档后，
键入要搜索的关键字，
回车选中备选项。

## TODO

- [] 本地文件缓存文档列表等数据。
- [] 在浏览器app中打开，而不是浏览器。
- [] 管理文档的缓存。
