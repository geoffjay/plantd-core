LAUNCHD_PATH = /Library/LaunchDaemons
SYSTEMD_PATH = /lib/systemd/system

M := $(shell printf "\033[34;1m▶\033[0m")

launchd: \
	install-launchd \
	enable-launchd \
	; $(info $(M) Installing launchd control files...)

install-launchd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[Broker|Identity|Proxy|State] make install-launchd"
else
	@install -m 644 "init/launchd/org.plantd.$(SERVICE).plist" "$(LAUNCHD_PATH)/"
endif

uninstall-launchd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[Broker|Identity|Proxy|State] make uninstall-launchd"
else
	@rm "$(LAUNCHD_PATH)/org.plantd.$(SERVICE).plist"
endif

enable-launchd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[Broker|Identity|Proxy|State] make enable-launchd"
else
	@launchctl load "$(LAUNCHD_PATH)/org.plantd.$(SERVICE).plist"
	@echo "start with: \`launchctl start org.plantd."$(SERVICE)`"
endif

disable-launchd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[Broker|Identity|Proxy|State] make disable-launchd"
else
	@launchctl stop "org.plantd.$(SERVICE)"
	@launchctl unload "$(LAUNCHD_PATH)/org.plantd.$(SERVICE).plist"
endif

systemd: \
	install-systemd \
	enable-systemd \
	; $(info $(M) Installing systemd control files...)

install-systemd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[broker|identity|proxy|state] make install-systemd"
else
	@mkdir -p /run/plantd
	@install -m 644 "init/systemd/plantd-$(SERVICE).service" "$(SYSTEMD_PATH)/plantd-$(SERVICE).service"
	@systemctl daemon-reload
endif

uninstall-systemd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[broker|identity|proxy|state] make uninstall-systemd"
else
	@rm "$(SYSTEMD_PATH)/plantd-$(SERVICE).service"
	@systemctl daemon-reload
endif

enable-systemd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[broker|identity|proxy|state] make enable-systemd"
else
	@systemctl enable "plantd-$(SERVICE)"
	@echo "start with: \`systemctl start plantd-"$(SERVICE)"\`"
endif

disable-systemd:
ifeq ($(SERVICE),)
	@echo "usage: SERVICE=[broker|identity|proxy|state] make disable-systemd"
else
	@systemctl stop "plantd-$(SERVICE)"
	@systemctl disable "plantd-$(SERVICE)"
endif

.PHONY: launchd systemd
