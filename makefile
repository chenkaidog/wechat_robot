docker_build:
	docker build -t openwechat:latest .

docker_run:
	docker run -itd -v ~/Documents/openwechat/log:/app/log -v ~/Documents/openwechat/storage:/app/storage --network=host openwechat:latest