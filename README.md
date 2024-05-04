# About

This is just a simple SSH server (mostly for testing/diagnostic purposes). For example, if you have a rooted Android and you want to quickly get SSH access to it, this backdoor is an easy solution.

# Quick start

### Build for ARM64
```sh
git clone https://github.com/xaionaro-go/backdoor
cd backdoor
CGO_ENABLED=0 GOARCH=arm64 go build
```

### Upload to a smartphone
```sh
adb push backdoor /sdcard/Download/
adb shell 'su -c "mkdir -p /data/backdoor; mv /sdcard/Download/backdoor /data/backdoor/; chmod +x /data/backdoor/backdoor"'
```

### Launch the backdoor
```sh
adb push ~/.ssh/id_ed25519.pub /sdcard/Download/authorized_keys
adb shell 'su -c "/data/backdoor/backdoor /bin/sh 0.0.0.0:8022 /sdcard/Download/authorized_keys"' &
```

### Get the IP address
```sh
PHONE_ADDR="$(adb shell ip a show dev wlan0 | grep 'inet ' | tr "/" " " | awk '{print $2}')"
```

### Login

```sh
ssh -p 8022 "$PHONE_ADDR"
```