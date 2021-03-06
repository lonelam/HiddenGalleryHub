echo 'Building...'
mkdir output
cd gallery-web-app
npm install
npm run build
cd ..
cp -r ./gallery-web-app/build/ output/
cd server
go mod tidy
go build -o ../output/server
cd ..
cd client
go mod tidy
go build -o ../output/client
cd ..
tar -zcvf output.tar.gz output
echo 'Done'