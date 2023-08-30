<p align="center">
  <a href="https://github.com/qishu321/cmdb-ops-flow">
    <img src="https://avatars.githubusercontent.com/u/95009146?s=400&u=0984e6a6a761fa007f6ad459abbb1ee9786424b8&v=4" alt="Logo" width="180" height="180">
  </a>

  <h1 align="center">cmdb-ops-flow</h1>
  <p align="center">
   本项目使用gin、gorm和ssh开发。旨在编写一个轻量，易用，多平台的运维项目。
    <br />
     <br />
  </p>
## 技术栈

#### 后端 Golang 1.19

- Gin 1.9.1             [(Web框架)](https://gin-gonic.com/zh-cn/)
- GORM v1.9.16     [(ORM)](https://gorm.io/zh_CN/)
- MySQL 5.7             [(数据库)](https://www.mysql.com/)

#### 前端 Vue.js 3

- vue-admin-template                     [(后台前端)](http://panjiachen.github.io/vue-admin-template)

#### 已实现的功能
- webssh
- cmdb
- 用户管理
- 批量命令或者脚本执行
- 简单的工作流执行
- etcd的备份和回档
- k8s的多集群管理
```bash
##k8s的多集群管理
- 目前实现的功能：
kubeconfig的管理，存储到数据库里，然后根据这个实现多集群的管理，多集群可以任意切换
web创建namespace、svc；web查看pod的日志、webssh登录pod、web获取集群监控汇总详情等
```

## 部署方法

###安装编译

```go
# clone
git clone https://github.com/qishu321/cmdb-ops-flow.git
##整体目录结构
├─api               api
├─conf              配置文件
├─middleware        中间件
├─models            数据库
├─router            路径
├─service           业务逻辑
├─test
└─utils             通用工具
    ├─common        加密解密
    ├─msg           状态码
    ├─result        状态码封装
    └─ssh           ssh



```


#### 如何参与开源项目

贡献使开源社区成为一个学习、激励和创造的绝佳场所。你所作的任何贡献都是**非常感谢**的。


1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### 版本控制

该项目使用Git进行版本管理。您可以在repository参看当前可用版本。

### 联系方式
## 预览
<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/kube-config.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/getpodyaml.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/webssh_pod.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/pod_log.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/kube-dashboard.png?raw=true" style="zoom: 25%;" />


<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/cmdb.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/etcdrestore.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/webssh2.png?raw=true" style="zoom:25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/%E5%91%BD%E4%BB%A4%E6%89%A7%E8%A1%8C.png?raw=true" style="zoom:25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/%E4%BD%9C%E4%B8%9A%E6%A8%A1%E6%9D%BF.png?raw=true" style="zoom:25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/%E4%BD%9C%E4%B8%9A%E6%89%A7%E8%A1%8C%E7%BB%93%E6%9E%9C.png?raw=true" style="zoom: 25%;" />

<img src="https://github.com/qishu321/cmdb-ops-flow/blob/main/doc/%E8%84%9A%E6%9C%AC%E5%BA%93.png?raw=true" style="zoom: 25%;" />


### 版权说明

该项目签署了MIT 授权许可

### 鸣谢

[vue-admin-template](http://panjiachen.github.io/vue-admin-template) 非常感谢vue-admin-template的开源项目

[webssh](https://github.com/widaT/webssh)               非常感谢widaT的webssh开源项目

[shaojintian README模板](https://github.com/shaojintian/Best_README_template)  非常感谢 shaojintian README模板

