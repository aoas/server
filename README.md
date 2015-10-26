# AOAS 通用服务端
AOAS是一个快速开发系统模板, 基于角色授权体系. 

整个系统以JSON格式来交互数据, 前台传过来的内容也是JSON格式字符串, 而非传统的form字段. 如在js端传过来的, 可用JSON.stringify(xxx)来格式化后再传过来.

为了方便移动端调用, 未用session功能, 而是用了`jwt token`

### 新增功能模块
若要增加模块功能. 可按如下步骤来操作.

* 在`models`目录增加新的model. 并在models/base.go下的`SyncTables`方法里增加新model的名字, 以便同步数据结构到数据库.
* 在`controllers`目录增加对应的操作控制器. 可继承Base这个控制器. 里面会带几个可能会用到的对象. 如logger, dbengine, config. 有某些情况下可能要用到config中的某些值. 另可在新的controller里写init方法. 把可能用到的权限写进去, 以便后续做授权操作.` 具体可参考controllers/user.go`
* 在routers/router.go里增对应的路由连接. 

调试时可用建议用[gin](https://github.com/codegangsta/gin)这类的控件, 以便实时刷新变动. `gin`一般用到3000端口, 但我们实际app的端口不是3000时, 可跟参数 `-a xxxx` 即可. 如我们在config中设置的app port是`8080`,  可用如下命令去启动gin: `gin -a 8080`. 此时程序gin以3000端口来启动.

### API调用
从客户端访问时, 有权限检查的API需要传如下header. 下面的token在用户调用`login`登陆后会得到.

```http
Authorization: Bearer DvjoEd6sKbHBLtMvrWWT
```
> Bearer后跟的是你调用login时拿到的token. 注意token和Bearer中间有个空格

### 基础API
默认API前缀地址为: `http://localhost:3000/api/...`, 如登陆操作URL为: `http://localhost:3000/api/login`.

##### 登陆/注册
URL|Method|Description|Permission
---|------|-----------|----------
login|POST|登陆
register|POST|注册账号

##### 用户相关
URL|Method|Description|Permission
---|------|-----------|----------
users|GET|查询用户列表|user.list
users/{id}|GET|查看指定用户信息|user.get
users/{id}/active|POST|禁用/启用户用户账号|user.active
users/{id}/roles|GET|用户角色列表|user.roles

##### 角色/权限相关
URL|Method|Description|Permission
---|------|-----------|----------
roles|GET|查询角色列表|role.list
roles/{id}|GET|查看指定用户信息|user.get
roles/{id}/users|GET|角色用户列表|role.users
roles/{id}/users|POST|增加用户到指定角色|role.adduser
roles/{id}/users|DELETE|从指定角色删除用户|role.deleteuser
roles/{id}/permissions|GET|角色可操作的权限列表|role.permissions
roles/{id}/permissions |POST|增加权限到指定角色|role.addpermission
roles/{id}/permissions |DELETE|从指定角色删除权限|role.deletepermissions
> 上面的增加/删除用户时, 需要传用户id列表. 如{"user_ids":[1,22]}. 增/删权限时, 需要传权限id列表, 如: {"permission_ids": ["user.active", "user.roles"]}

##### 文件操作相关
URL|Method|Description|Permission
---|------|-----------|----------
files|GET|查询用户上传文件记录列表|file.list
files|POST|上传文件(file对象名为`file`)|file.upload



### 用到的库   
* [gin](https://github.com/gin-gonic/gin) http framework
* [xorm](https://github.com/go-xorm/xorm) 数据库orm   
* [toml](https://github.com/BurntSushi/toml) 解析config
* [jwt-go](https://github.com/dgrijalva/jwt-go) 生成token相关

### 工具推荐
调试调用API时, 我推荐 [Insomnia](http://insomnia.rest/), 整个用下来非常不错. 尤其支持变量及组功能相对有用. 


