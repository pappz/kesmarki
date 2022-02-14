# Compile to Raspberry pi

```
GOARCH=arm64 go build
```

# Setup users

Create a new file "/etc/kesmarki/users". Fill it with users in separated lines with syntax like below.

```
user1:password
user2:password
```