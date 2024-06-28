## Personal exercise project, an online chat system by Go

### structure :

```

|-config           //project config   
|-controller       //http operation in business logic,user and admin module
|-docs             //swagger docs
|-logger           //record log by thread pool
|-middleware       //middleware
|  └─auth          //authorization ...
|  └─blockIP       //ip blocked operation
|  └─log           //middleware of logger
|-model            //database model, hook
|-pkg              //designs for service: snowflake, jwt
|-response         //customize response information
|-router           //router 
|-service          //model operation in business service
|-static           //webpage front-end 
|  └─assert        //assert of front-end
|  └─css           //css
|  └─js            //js
|-template         //html template
|-utils            //md5 ...
```

#### what will be added:

msg quene, nginx, redis, token validator...  (distributed system frame?)

#### other more:

Due to I am a student, I always delay to update by examination or other reasons.

#### By the way:

F***k you CSDN


