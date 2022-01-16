build:
	echo "building binary"
	go build -o bin/battery-measure

clean:
	echo "cleaning bin dir"
	rm -r bin

deploy:
	mv bin/battery-measure /data/cellar/

build_and_deploy: clean build deploy