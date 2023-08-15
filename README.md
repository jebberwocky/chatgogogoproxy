# GOGOGOProxy

## rewrite OpenAI proxy service in Golang

## build
Ubuntu
````
env GOOS=linux GOARCH=amd64 go build
````
Optimizing
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o echo-server server.go
```
## app config file
place apps.json to the root directory as executable.
```json
{
  "apps": [{
    "id": "2c094367-dc1e-4fb1-93cd-d7068e3e4278",
    "name": "monkey",
    "hostname": "colbt.cc:3131",
    "openai_key": "---MASKED---",
    "role":{
      "role":    "system",
      "content": "你是一个智能聊天助手"}
    },{
    "id": "0dce57ae-7412-4a78-85fc-67f88f114f0d",
    "name": "maas",
    "hostname": "maasdemo.colbt.cc",
    "openai_key": "---MASKED---",
    "role":{
      "role":    "system",
      "content": "你是一个智能销售pitch邮件生成器"}
  }]
}

```