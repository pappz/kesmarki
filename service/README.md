# Compile to Raspberry pi

```
GOOS=linux GOARCH=arm GOARM=5 go build
```

# Setup users

Create a new file "/etc/kesmarki/users". Fill it with users in separated lines with syntax like below.

```
user1:password
user2:password
```

# Systemd service file

location: /etc/systemd/system/kesmarki.service
Setup step:
```
systemctl daemon-reload
systemctl enable kesmarki.service
systemctl start kesmarki.service
```
```
[Unit]
Description=kesmarki

[Service]
PIDFile=/run/kesmarki.pid
User=pi
Group=pi
StandardOutput=syslog
StandardError=syslog
LimitNOFILE=49152
ExecStart=/usr/local/bin/kesmarki
Restart=on-failure
EnvironmentFile=-/etc/kesmarki/config.env

[Install]
WantedBy=multi-user.target
```

# Setup environment variables
location: /etc/kesmarki/config.env
```
KM_WOL_BUDAFOKI=aa:bb:00:00:00:00
```