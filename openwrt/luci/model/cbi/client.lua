require("luci.sys")

local desc = [[
<h4>探针配置<h4/>
<p>服务器搭建教程: https://github.com/cppla/ServerStatus <br/>本项目地址：https://github.com/v03413/ServerStatus-Client</p>
]]

m = Map("client", translate("ServerStatus Client"), desc)

-- # 读取配置文件
s = m:section(TypedSection, "server", "")
s.addremove = false
s.anonymous = true

enable = s:option(Flag, "enable", translate("Enable"))
enable.rmempty = false

server = s:option(Value, "server", translate("Server"))
server.datatype = "host"
server.rmempty = false

port = s:option(Value, "port", translate("Port"))
port.datatype = "port"
port.rmempty = false

user = s:option(Value, "user", translate("Username"))
user.datatype = "minlength(1)"
user.rmempty = false

password = s:option(Value, "password", translate("Password"))
password.password = true
password.rmempty = false
password.datatype = "minlength(1)"

local apply = luci.http.formvalue("cbi.apply")
if apply then
    io.popen("/etc/init.d/client restart > /dev/null &")
end

return m