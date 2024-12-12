
docker pull namely/protoc-all

docker run -v $PWD:/defs namely/protoc-all -f service.proto -l go -o .