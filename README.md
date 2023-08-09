### just a user-group-permission wrapper of jenkins, not a ci system.

## concept
```
项目 = jenkins folder
应用 = jenkins job
发布 = jenkins build
```

## default role and permission

|       | 角色名  | 权限 |
| -----|  ----  | ---- | 
| admin |  管理员   | 默认拥有所有权限  |
| ops   |  运维人员  | 项目、应用、发布所有权限；设置项目管理员等  |
| dev   |  开发人员  | 项目、应用只读权限；发布相关权限 |
| project_manager  | 项目管理员  | 除dev角色权限外,可修改项目组、添加删除应用、设置项目成员、审核发布等  |


## build
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

## run 
```
cd /usr/local/jenkins-wrapper-ci
./bin/jenkins-wrapper-ci -c ./conf/config.yaml
```


## init
```
1. 访问localhost:8080/init/initdb, 根据提示初始化数据库
2. 修改config.yaml
3. 默认用户名密码: admin/123456
```

## nginx conf
see [jenkins-wrapper-ci.conf ](web/nginx/conf.d/jenkins-wrapper-ci.conf)
