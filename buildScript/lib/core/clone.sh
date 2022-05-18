cd ..
[ -d v2ray-core ] && exit 0
rm -rf v2ray-core
pwd
ls
git clone https://github.com/matsuridayo/v2ray-core.git
cd v2ray
