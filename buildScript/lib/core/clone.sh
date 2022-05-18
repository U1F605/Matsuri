cd ..
[ -d v2ray ] && exit 0
rm -rf v2ray-core
pwd
ls -l
git clone https://github.com/matsuridayo/v2ray-core.git
cd v2ray
