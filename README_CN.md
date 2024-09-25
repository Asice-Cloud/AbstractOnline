## 这里是简体中文文档
## 这个项目本来是要做一个简单的聊天室，但是后续增加了很多小内容，就直接做一个聊天的框架吧

[English](README.md) | 简体中文
```
|-cache            //缓存常用页面
|-config           //服务的配置   
|-controller       //涉及网路请求的处理和响应
|  └─admin_module  //admin模块的处理
|  └─authorization //认证模块
|  └─user_module   //用户模块
|  └─verification  //验证模块，图形验证码等内容
|  └─websocket_work //基于websocket进行的聊天室的搭建
|-docs             //swagger 文档
|-middleware       //中间件
|  └─auth          //验证中间件
|  └─balance       //负载均衡中间件
|  └─blockIP       //ip 黑名单
|  └─log           //日志记录中间件
|  └─safe          //csrf, xss, sql注入等防护中间件
|-model            //数据库模型以及钩子
|-pkg              //服务器的算法： 雪花、jwt等
|-response         //自定义http处理后的响应内容
|-router           //路由组
|-service          //业务代码，crud
|-session          //支持多种session的存储和多种方式
|-static           //前端的一些资源 
|  └─img           //img
|  └─css           //css
|  └─js            //js
|-template         //html页面
|-utils            //小工具
```
