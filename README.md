### just a user-group-permission wrapper of jenkins, not a ci system.


#### build
```
mkdir -p /usr/local/jenkins-wrapper-ci/{ui,bin,conf,log}

cd web/
npm install
npm run build
cp -rf dist/*  /usr/local/jenkins-wrapper-ci/ui/

cd ../server/
cp -rf resource /usr/local/jenkins-wrapper-ci/
cp -rf config.yaml /usr/local/jenkins-wrapper-ci/conf/
go build -o jenkins-wrapper-ci main.go
cp -rf jenkins-wrapper-ci /usr/local/jenkins-wrapper-ci/bin/
chmod a+x /usr/local/jenkins-wrapper-ci/bin/jenkins-wrapper-ci

```

#### run 
```
cd /usr/local/jenkins-wrapper-ci
./bin/jenkins-wrapper-ci -c ./conf/config.yaml
```


#### nginx conf
```
upstream backend {
        server localhost:8888;
}

server {
        listen   8080;

        location / {
                alias /usr/local/jenkins-wrapper-ci/ui/;
                index index.html;
                try_files $uri /index.html;
        }

        location ^~ /api/ {
                proxy_set_header Host $http_host;
                proxy_set_header  X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_pass http://backend/;
        }

        location ^~ /swagger/ {
                proxy_set_header Host $http_host;
                proxy_set_header  X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                proxy_pass http://backend;
        }

        error_page 404 /404.html;
        location = /404.html {
        }

        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
        }
}
```