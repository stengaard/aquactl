ifndef PI_HOST
	echo "Please set PI_HOST env var" $(PI_HOST)
	exit 1
endif

TARGET=aquactl

$(TARGET).pi : *.go
	GOARM=6 GOOS=linux GOARCH=arm go build -o $@

clean:
	rm *.pi

upload: $(TARGET).pi
	ssh $(PI_HOST) sudo pkill $(TARGET).pi || true
	scp $(TARGET).pi $(PI_HOST):.

devrun: upload
	ssh $(PI_HOST) sudo ./$(TARGET).pi

install: upload
	scp $(TARGET).boot.sh $(PI_HOST):.
	ssh $(PI_HOST) -- sudo mv ./$(TARGET).pi /usr/sbin/$(TARGET).pi \&\& sudo mv ./$(TARGET).boot.sh /etc/init.d \&\& sudo update-rc.d $(TARGET).boot.sh defaults


run: install
	ssh $(PI_HOST) -- sudo invoke-rc.d $(TARGET).boot.sh restart


.PHONY: upload clean devrun install run
