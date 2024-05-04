# Free chat to api

## 功能
- 免登chat to api
- 免登网关

## 安装
```bash
docker run -d -p 8484:8484 --name free-chat-to-api ghcr.io/zbronya/free-chat-to-api:latest
```

## 使用
```bash
curl 'http://127.0.0.1:8484/v1/chat/completions' \
--header 'Content-Type: application/json' \
--data '{
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "user",
            "content": "Say this is a test! "
        }
    ],
    "stream": true
}'
```

## 环境变量
- `PORT` - port to listen on (default: `8484`)
- `BIND` - address to bind to (default: `0.0.0.0`)
- `PROXY_URL` - proxy url
- `GATEWAY_URL` - gateway url (default: `https://chatgpt.com`)
