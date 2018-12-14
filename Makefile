

build :
		CGO_ENABLED=0 GOOS=linux go1.11 build -a -installsuffix cgo -o fast-storage cmd

docker-build :
		docker build -t fast-storage .

docker-run :
		docker run --name fast-storage --link some-redis:redis -d fast-storage
		docker logs fast-storage

docker-stop:
		docker stop fast-storage
		docker rm fast-storage

docker-run-redis:
		docker run --name some-redis -d redis