# Golang 版本的探针客户端

云探针 https://github.com/cppla/ServerStatus ，Golang版客户端；同时增加了对`OpenWRT`的适配、**也就是说利用此项目，您可以给自己的路由器挂个探针。**  

## 开发原因  

- 学习Golang，自我提升
- 疫情原因，闲置在家
- 原版客户端均为Python脚本，运行依赖系统Python环境，开发此版本的目的就是为了摆脱环境，同时对`OpenWrt`进行适配。

## 使用教程

普通服务器，使用方法同原版Python脚本一样，只是由脚本换成了本项目的可执行文件，具体方法不再阐述。  

---  

## OpenWRT  

### Luci & 截图
![ServerStatus Client](https://raw.githubusercontent.com/v03413/ServerStatus-Client/main/openwrt/images/luci.png)

### 编译 & 使用

假设你的Lean openwrt 在 lede 目录下

```bash
# 进入包目录
cd lede/package/  

# 拉取编译文件
git clone https://github.com/v03413/ServerStatus-Client.git

# 返回编译根目录
cd ../

# 开始编译插件
make package/ServerStatus-Client/compile -j1 V=99
```

如果没看到error相关信息，说明编译成功；最后的ipk插件文件应该在：`{lede}/bin/packages/x86_64/base`目录，得到文件：`ServerStatus-Client_**.ipk`
