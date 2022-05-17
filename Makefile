NAME=cronlocker
VERSION=0.6.0
ARCH=amd64
LICENSE="BSD 2-Clause License"
MAINTAINER="Tobias Schoknecht <tobias.schoknecht@viafintech.com>"
DESCRIPTION="Distributed sequential cronjob running tool."
HOMEPAGE="https://www.viafintech.com/"
GIT="git@github.com:viafintech/cronlocker.git"
GITBROWSER="https://github.com/viafintech/cronlocker"

.PHONY: test testci package cronlocker

all: cronlocker

test:
	go test -v `go list ./...|grep -v vendor`

testci:
	go test -v `go list ./...|grep -v vendor` -tags=docker


cronlocker: *.go */*.go
	go build -ldflags="-s -w" .

install: cronlocker
	install -d $(DESTDIR)/usr/bin
	install cronlocker $(DESTDIR)/usr/bin/cronlocker

DESTDIR=tmp

package: install cronlocker
	install -d $(DESTDIR)/usr/share/doc/$(NAME)
	install --mode=644 package/copyright $(DESTDIR)/usr/share/doc/$(NAME)/copyright
	fpm -s dir -t deb -n $(NAME) -v $(VERSION) -C $(DESTDIR) \
	-p $(NAME)_$(VERSION)_$(ARCH).deb \
	--license $(LICENSE) \
	--maintainer $(MAINTAINER) \
	--vendor $(MAINTAINER) \
	--description $(DESCRIPTION) \
	--url $(HOMEPAGE) \
	--deb-field 'Vcs-Git: $(GIT)' \
	--deb-field 'Vcs-Browser: $(GITBROWSER)' \
	--deb-upstream-changelog package/changelog \
	--deb-no-default-config-files \
	usr/bin usr/share/doc/

clean:
	rm -f cronlocker
	rm -f cronlocker*.deb
	rm -fr $(DESTDIR)
