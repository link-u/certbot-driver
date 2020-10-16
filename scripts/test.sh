#! /bin/bash -eux

set -eux

function shouldBeSingleBinary() {
  local path
  path="$(command -v "$1")"
  if ldd "${path}" 2> /dev/null; then
    echo "NG: ${path} is not fully statically linked."
    exit 255
  else
    echo "OK: ${path} is fully statically linked."
  fi
}

## git リポジトリ上の root のパスを取得
scripts_dir=$(cd "$(dirname "$(readlink -f $0)")" && pwd)
root_dir=$(cd "${scripts_dir}" && cd .. && pwd)
cd "${root_dir}"

dpkg-deb --contents ./artifact/*.deb

apt install -y ./artifact/*.deb
apt show certbot-driver

command -v certbot-driver
certbot-driver --version
certbot-driver --help
shouldBeSingleBinary "certbot-driver"
