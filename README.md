## Personal exercise project, an online chat system by Go

### structure :

```
|-cache            //usual pages and query cache
|-config           //project config   
|-controller       //http operation in business logic,user and admin module
|  └─admin_module  //handle admin operation
|  └─authorization //some authorization operation
|  └─user_module   //handle user operation
|  └─verification  //verification of user, when logining or entering chat room...
|  └─websocket_work //websocket server and char room
|-docs             //swagger docs
|-logger           //record log by thread pool
|-middleware       //middleware
|  └─auth          //authorization ...
|  └─balance       //LoadBalance
|  └─blockIP       //ip blocked operation
|  └─log           //middleware of logger
|  └─safe          //defense of csrf, xss, sql injection...
|-model            //database model, hook
|-pkg              //designs for service: snowflake, jwt
|-response         //customize response information
|-router           //router 
|-service          //model operation in business service
|-session          //operations of multiply session: stroe, set, get ...
|-static           //webpage front-end 
|  └─img           //img
|  └─css           //css
|  └─js            //js
|-template         //html template
|-utils            //tool and mechanism ...
```

#### what will be added:

msg quene, nginx, token validator... 

#### By the way:

F***k you CSDN
