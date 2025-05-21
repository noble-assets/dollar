cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r dollar.noble.xyz/v2/* ./
cp -r api/noble/dollar/* api/
find api/ -type f -name "*.go" -exec sed -i 's|dollar.noble.xyz/v2/api/noble/dollar|dollar.noble.xyz/v2/api|g' {} +

rm -rf dollar.noble.xyz/v2
rm -rf api/noble
rm -rf noble
