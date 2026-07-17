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
StandardOutput=journal
StandardError=journal
LimitNOFILE=49152
ExecStart=/usr/local/bin/kesmarki
Restart=on-failure
EnvironmentFile=-/etc/kesmarki/config.env

[Install]
WantedBy=multi-user.target
```

# View logs

The service logs to the systemd journal. Follow live output with:
```
sudo journalctl -u kesmarki -f
```
Note the `-u` (unit) flag — `journalctl kesmarki` fails with "Invalid argument".

# Setup environment variables
location: /etc/kesmarki/config.env
```
KM_WOL_BUDAFOKI=aa:bb:00:00:00:00
# IP address pinged to report whether the budafoki PC is online.
# Defaults to 192.168.0.10 when unset.
KM_WOL_BUDAFOKI_IP=192.168.0.10

# Dynamic DNS (optional). When both KM_DDNS_* below are set, the service keeps
# the given Route 53 A record pointed at the machine's current public IP.
KM_DDNS_HOSTED_ZONE_ID=Z0123456789ABCDEFGHIJ
KM_DDNS_RECORD_NAME=kesmarki.godevltd.com
KM_DDNS_INTERVAL=1h

# AWS credentials for the Route 53 update (standard AWS SDK env vars).
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...
```

## Dynamic DNS

If the machine sits behind a dynamic public IP, set the `KM_DDNS_*` variables
above. On startup and then every `KM_DDNS_INTERVAL` (default `1h`) the service
looks up its public IP and, when it changed, UPSERTs the A record via the
Route 53 API. Leave `KM_DDNS_HOSTED_ZONE_ID` / `KM_DDNS_RECORD_NAME` unset to
disable the feature — the rest of the service is unaffected.

The IAM user/role behind the AWS credentials only needs permission to change
records in the one hosted zone:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "route53:ChangeResourceRecordSets",
      "Resource": "arn:aws:route53:::hostedzone/Z0123456789ABCDEFGHIJ"
    }
  ]
}
```

Find the hosted zone ID in the Route 53 console (or with
`aws route53 list-hosted-zones-by-name --dns-name godevltd.com`).