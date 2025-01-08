cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r dollar.noble.xyz/* ./
cp -r api/noble/dollar/* api/
find api/ -type f -name "*.go" -exec sed -i 's|dollar.noble.xyz/api/noble/dollar|dollar.noble.xyz/api|g' {} +

rm -rf dollar.noble.xyz
rm -rf api/noble
rm -rf noble
