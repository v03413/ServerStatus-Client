include $(TOPDIR)/rules.mk

PKG_NAME:=ServerStatus-Client
PKG_VERSION:=0.15
PKG_RELEASE:=$(AUTORELEASE)
PKG_BUILD_DIR := $(BUILD_DIR)/$(PKG_NAME)

PKG_LICENSE:=GPLV3
PKG_LICENSE_openwrt:=LICENSE
PKG_MAINTAINER:=V03413 <admin@qzone.work>

PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_USE_MIPS16:=0

GO_PKG:=github.com/v03413/ServerStatus-Client
GO_PKG_BUILD_PKG:=$(GO_PKG)/cmd/client
GO_PKG_LDFLAGS_X:= \
	$(GO_PKG).build=OpenWrt \
	$(GO_PKG).version=$(PKG_VERSION)

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

define Build/Prepare
	$(CP) ./* $(PKG_BUILD_DIR)/
endef

define Package/ServerStatus-Client
	TITLE:=ServerStatus Client
	SECTION:=net
	CATEGORY:=Network
	URL:=https://github.com/v03413
	DEPENDS:=$(GO_ARCH_DEPENDS) +ca-bundle
endef

define Package/ServerStatus-Client/install
	$(call GoPackage/Package/Install/Bin,$(PKG_INSTALL_DIR))

	# 安装执行文件
	$(INSTALL_DIR) $(1)/usr/bin/
	$(INSTALL_BIN) $(PKG_INSTALL_DIR)/usr/bin/client $(1)/usr/bin/client

	# 安装后台插件
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/model/cbi
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/controller

	$(INSTALL_CONF) ./openwrt/root/etc/config/client $(1)/etc/config/client
	$(INSTALL_BIN) ./openwrt/root/etc/init.d/client $(1)/etc/init.d/client
	$(INSTALL_DATA) ./openwrt/luci/model/cbi/client.lua $(1)/usr/lib/lua/luci/model/cbi/client.lua
	$(INSTALL_DATA) ./openwrt/luci/controller/client.lua $(1)/usr/lib/lua/luci/controller/client.lua
endef

$(eval $(call BuildPackage,ServerStatus-Client))