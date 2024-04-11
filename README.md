## Personal exercise project, a online chat system by Go



### structure :

```

|-config           //project config   
|-controller       //http operation in business logic 
|-docs             //swagger docs
|-logger           //record log
|-middleware       //middleware
|  └─auth          //verify ip, login ...
|  └─log           //middleware of logger
|-model            //database
|-pkg              //algorithm for service
|-response         //customize response information
|-router           //router
|-service          //model operation in business service
|  └─validator     //check for data
|-static           //webpage front-end 
|  └─assert        //assert of front-end
|  └─css           //css
|  └─js            //js
|-template         //html template
```

#### what will be added:

msg quene, nginx, redis cache, token validator, grpc 
