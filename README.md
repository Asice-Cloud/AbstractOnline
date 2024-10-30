## Personal exercise project, an online chat system by Go
## Now I'm trying to build this project to a chat framework
### structure :

```
|-cache            //usual pages and query cache
|-config           //project config   
|-controller       //http operation in business logic,user and admin module
|  └─admin_module  //handle admin operation
|  └─authorization //some authorization operation
|  └─user_module   //handle user operation
|  └─verification  //verification of user, when logining or entering chat room...
|-docs             //swagger docs
|-middleware       //middleware
|  └─auth          //authorization ...
|  └─balance       //LoadBalance
|  └─blockIP       //ip blocked operation
|  └─log           //middleware of logger
|  └─safe          //defense of csrf, xss, sql injection...
|-model            //database model, hook
|-response         //customize response information
|-router           //router 
|-server           //designs for service: snowflake, jwt
|-service          //model operation in business service
|-session          //operations of multiply session: stroe, set, get ...
|-static           //resource of front-end 
|  └─img           //img
|  └─css           //css
|  └─js            //js
|-template         //html template
|-utils            //tool and mechanism ...
|-websocket_work   //websocket server and char room
```

#### what will be added:

msg quene, nginx, token validator... 

#### By the way:

F***k you CSDN
