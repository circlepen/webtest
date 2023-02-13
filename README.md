# webtest


### start

```
(build the image)
docker-compose up --build -d


(start container, no build)
docker-compose up -d
```

### ports

+ api: localhost:8000
+ mongodb: localhost:27017


### upload files example
```
curl -X POST http://localhost:8000/api/upload -F "file[]=@/Users/yijulai/Documents/test.txt" -F "file[]=@/Users/yijulai/Documents/test2.txt" -H "Content-Type: multipart/form-data"
```
