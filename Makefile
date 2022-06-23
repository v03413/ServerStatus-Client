include $(TOPDIR)/rules.mk

PKG_NAME:=ServerStatus-Client
PKG_VERSION:=0.15
PKG_RELEASE:=$(AUTORELEASE)

PKG_SOURCE:=$(PKG_NAME)-$(PKG_VERSION).tar.gz
PKG_SOURCE_URL:=https://codeload.github.com/v03413/ServerStatus-Client/tar.gz/v$(PKG_VERSION)?
PKG_HASH:=22ac41cdee2333dff9dcf82ab39b2cd118dfeb8b353342076257b48365259484

PKG_LICENSE:=GPLV3
PKG_LICENSE_FILES:=LICENSE
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

define Package/ServerStatus-Client
	TITLE:=ServerStatus Client
	SECTION:=net
	CATEGORY:=Network
	URL:=https://github.com/v03413
	DEPENDS:=$(GO_ARCH_DEPENDS) +ca-bundle
endef

define Package/ServerStatus-Client/install
	$(call GoPackage/Package/Install/Bin,$(PKG_INSTALL_DIR))

	$(INSTALL_DIR) $(1)/usr/bin/
	$(INSTALL_BIN) $(PKG_INSTALL_DIR)/usr/bin/client $(1)/usr/bin/client
endef

$(eval $(call BuildPackage,ServerStatus-Client))