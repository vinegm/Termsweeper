# Maintainer: vinegm <info AT vinegm D0T dev>
pkgname='termsweeper-bin'
pkgver=1.1.0
pkgrel=1
pkgdesc="Customizable minesweeper terminal UI"
arch=('x86_64' 'aarch64')
url="https://github.com/vinegm/termsweeper"
license=('MIT')
provides=('termsweeper')
conflicts=('termsweeper')
options=('!strip' '!debug')
source_x86_64=("termsweeper-${pkgver}-x86_64.tar.gz::${url}/releases/download/v${pkgver}/termsweeper_${pkgver}_linux_amd64.tar.gz")
source_aarch64=("termsweeper-${pkgver}-aarch64.tar.gz::${url}/releases/download/v${pkgver}/termsweeper_${pkgver}_linux_arm64.tar.gz")
sha256sums_x86_64=('b47846b418aa2d5a5c629fa0ce753dd3ba3299ffd47dc778286161dab37227a7')
sha256sums_aarch64=('40a84ee8aa02fd6c9b9ac0335f29c91604e87755567eb81669b6b03c458881eb')

package() {
    install -Dm755 termsweeper "$pkgdir/usr/bin/termsweeper"
    install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
    install -Dm644 -t "$pkgdir/usr/share/termsweeper/Themes" Themes/*
}
