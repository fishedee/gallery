.PHONY:dev
dev:
	-pkill server
	fishcmd run --watch server
release:
	go build .
	mv server fishindex
	-pkill fishindex
	nohup ./fishindex &
