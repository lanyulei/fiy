<p align="center">
  <img src="https://www.fdevops.com/wp-content/uploads/2020/09/1599039924-ferry_log.png">
</p>


<p align="center">
  <a href="https://github.com/lanyulei/ferry">
    <img src="https://www.fdevops.com/wp-content/uploads/2020/07/1595067271-badge.png">
  </a>
  <a href="https://github.com/lanyulei/ferry">
    <img src="https://www.fdevops.com/wp-content/uploads/2020/07/1595067272-apistatus.png" alt="license">
  </a>
    <a href="https://github.com/lanyulei/ferry">
    <img src="https://www.fdevops.com/wp-content/uploads/2020/07/1595067269-donate.png" alt="donate">
  </a>
</p>

# 基于Gin + Vue + Element UI前后端分离的资源管理系统，即CMDB

本项目使用数据模型进行数据校验规则的管理及维护，方便进行数据结构的变更维护及管理。

支持业务树结构进行数据的分类划分，多样化的资源搜索功能，结合ES实现的全局检索功能，数据的导入导出等等功能。

支持云厂商数据资源数据的同步。

操作审计的详细记录及数据调整前后的人性化对比。

等等，还有更多更详细的功能，可进行演示站点试用了解。

演示站点：http://fdevops.com:8060

账号密码：admin/123456

很多功能是参照蓝鲸CMDB，代码实现及数据结构设计完全自主，若是觉得跟蓝鲸CMDB有点类似，还请勿喷。

前端UI：https://github.com/lanyulei/fiy-ui

### 资源管理

* 统一搜索，使用 canal 同步MySQL数据到ES，通过ES进行全局数据搜索功能。
* 业务
    * 业务拓扑，对业务数据进行梳理，绑定及展示业务数据。
    * 服务模版，配置资源的服务模版，服务模版可配置进程参数信息，方便基于服务模版进行自动化任务，资源可选择是否绑定服务模版。
    * 集群模版，集群绑定服务模版，可基于集群进行自动化及批处理任务。
    * 服务分类，服务运行的是什么类型的服务，例如MySQL、Redis或者自定义的内部服务等。
* 模型
    * 模型管理，管理数据的模型及模型字段，数据会根据模型的特定规则写入数据，通过对模型字段定义规则即可在数据写入的时候，根据模型的字段规则进行数据的校验。
    * 模型管理，可视化查看模型的上下游关系。
    * 关联类型，模型关联的关联类型管理。
* 资源
    * 资源目录，模型所对应的数据入口。
    * 云账号，各家云厂商账号的管理。
    * 云资源同步，通过创建同步任务，并绑定云账号，进行云资源账号的数据同步。
* 运营分析
    * 操作审计，记录用户每次对资源数据写操作的详情，并进行前后数据的对比。

### 系统管理

基于casbin的RBAC权限控制，借鉴了go-admin项目的前端权限管理，可以在页面对API、菜单、页面按钮等操作，进行灵活且简单的配置。

### 系统工具

当前服务监控及系统配置。

# 打赏

> 如果您觉得这个项目帮助到了您，您可以请作者喝一杯咖啡表示鼓励:

<img class="no-margin" src="https://www.fdevops.com/wp-content/uploads/2020/07/1595075890-81595075871_.pic_hd.png"  height="200px" >

# License

开源不易，请尊重作者的付出，感谢。

[MIT License](https://github.com/lanyulei/fiy/blob/master/LICENSE)

Copyright (c) 2021 lanyulei
