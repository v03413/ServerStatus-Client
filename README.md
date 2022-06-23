# Golang 版本的探针客户端

云探针 https://github.com/cppla/ServerStatus ，Golang版客户端

## 开发原因

- 学习Golang，自我提升
- 疫情原因，闲置在家
- 原版客户端均为Python脚本，运行依赖系统Python环境，开发此版本的目的就是为了摆脱环境，同时对`OpenWrt`进行适配。

Ps：经过半个月的测试，未发现什么大问题，可正常使用。

## OpenWRT 插件编译

假设你的Lean openwrt 在 lede 目录下

```bash
# 进入包目录
cd lede/package/  

# 拉取编译文件
git clone https://github.com/jerrykuku/luci-app-vssr.git 

# 返回编译根目录
cd ../

# 开始编译插件
make package/ServerStatus-Client/compile -j1 V=99
```

如果没看到error相关信息，说明编译成功；最后的ipk插件文件应该在：`{lede}/bin/packages/x86_64/base`目录，得到文件：`ServerStatus-Client_**.ipk`

## 使用教程

普通X86 Arm-v8架构机器，Release下载对应版本即可直接运行，参数同原版一样。