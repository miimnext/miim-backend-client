# 第一阶段：构建 Go 二进制
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o main .

# 第二阶段：创建精简运行环境
FROM alpine
WORKDIR /root/
COPY --from=builder /app/main .
#  确保将 .env文件复制到容器中
COPY .env /root/.env  
# 确保可执行文件
RUN chmod +x main 
CMD ["./main"]
