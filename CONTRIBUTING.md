## 为 CF 贡献代码

首先很高兴你看到这个，非常欢迎您和我一起完善 CF，让我们一起编写一个好用的、有价值的、有意义的工具。

### 是否发现了 bug?

如果你发现了程序中的 bug，建议可以先在 [issue](https://github.com/teamssix/cf/issues) 中进行搜索，看看以前有没有人提交过相同的 bug。

### 发现了 bug，并准备自己为 CF 打补丁

如果你发现了 bug，而且准备自己去修补它的话，则可以先把 CF Fork 到自己的仓库中，然后修复完后，提交 PR 到 CF 的 beta 分支下，另外记得在 PR 中详细说明 bug 的产生条件以及你是如何修复的，这可以帮助你更快的使这个 PR 通过。

### 为 CF 增加新功能

如果你发现 CF 现在的功能不能满足自己的需求，并打算自己去增加上这个功能，则可以先把 CF Fork 到自己的仓库中，然后编写完相应的功能后，提交 PR 到 CF 的 beta 分支下，另外记得在 PR 中详细说明增加的功能以及为什么要增加它，这可以帮助你更快的使这个 PR 通过。

建议师傅在增加新功能后，可以去完善对应的操作手册，操作手册项目地址：[github.com/teamssix/TWiki](https://github.com/teamssix/TWiki)，先 fork 操作手册项目，然后在 [github.com/teamssix/TWiki/tree/main/docs/CF](https://github.com/teamssix/TWiki/tree/main/docs/CF) 目录下编辑 CF 操作手册的文档，最后提 pr 到 T Wiki 项目的 beta 分支下就行。

## 使你的 PR 更规范

### 规范 Git commit

commit message 应尽可能的规范，为了保证和其他 commit 统一，建议 commit message 采用英文描述，并符合下面的格式：

```yaml
type: subject
```

type 是 commit 的类别，subject 是对这个 commit 的英文描述，例如下面的示例：

* feat: add object download function

* perf: optimize the display of update progress bar

关于 Git commit 的更多规范要求可以参见这篇文章：[如何规范你的Git commit？](https://zhuanlan.zhihu.com/p/182553920)

### 规范代码内容

为了面向不同的使用人群，CF 里的输出信息使用了双语，因此如果你打算为 CF 增加新功能，那么在编写输出信息时应尽可能符合下面的格式：

```yaml
中文 (English)
```

注意上面的括号是英文括号，并且英文开头应该首字母大写，例如下面的示例：

* 关于作者 (About me)
* 配置云服务商的访问密钥 (Configure cloud provider access key)

### 规范代码格式

代码在编写完后，应使用 `go fmt` 和 `goimports`对代码进行格式化，从而使代码看起来更加整洁。

在 goland 中可以设置当代码保存时自动格式化代码，配置的步骤可以参见这篇文章：[GoLand：设置gofmt与goimports，保存时自动格式化代码](https://blog.csdn.net/qq_32907195/article/details/116755338)