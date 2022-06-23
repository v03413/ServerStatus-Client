module("luci.controller.client", package.seeall)

function index()
    entry({"admin", "services", "client"}, cbi("client"), _("ServerStatus Client"), 1)
end