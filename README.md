# IoT Server with Golang Fiber
IoT Server for fulfilling my thesis at IPB University. In This research, researcher successfully developed a high performance IoT back-end server based on REST API with Golang Fiber framework, then compare its performance and memory usage effectiveness to application developed by previous research which use Python Falcon and Python Sanic framework. The method used for developing the application is waterfall method. Based on the result of the performance test, Golang Fiber application have transaction per second (TPS) performance 73 times faster than Python Falcon and 2 times faster than Python Sanic. Golang Fiber application also have memory usage 5 times more effective than Python Falcon and 1,6 times more effective than Python Sanic.

## Background
The importance of an IoT server lies in its ability to provide a platform for collecting, processing, analyzing, and storing large amounts of data generated by IoT devices. It enables real-time communication between devices and applications, allowing for automation, monitoring, and control of various processes. By using an IoT server, organizations can achieve greater efficiency, cost savings, and improved decision-making capabilities.

## Why
In this study, researchers will compare the REST API framework in the Golang programming language to create an IoT server. Researchers will use the Fiber framework to create a backend server to analyze and receive IoT data. Researchers chose the Golang programming language because this programming language has very good performance in the effectiveness of CPU and memory allocation (Effendy et al. 2021), this shows that this language is very suitable for making IoT servers that get very large data input. The Fiber framework was chosen because this framework has high performance compared to other Golang frameworks, Fiber is listed as a framework that has the least processing time and the highest number of responses per second compared to other Golang frameworks (TechEmpower 2022), this framework is the 5th most popular Golang framework on Github (Kwon 2022). For the database, researchers will use the PostgreSQL database to match other studies so that the database becomes the control variable. Researchers will analyze the performance and effectiveness of the server using transaction per second metrics and the percentage of memory usage, then compare the results with previous research by Hanin (2021) using the Falcon framework and Alvin (2023) using the Sanic framework. This is done to analyze the performance of the Fiber framework compared to the Falcon and Sanic frameworks.

## Problem
1. Does Golang with the Fiber framework have good performance in making IoT servers?
2. How does the performance of the Fiber framework compare to the Falcon and Sanic frameworks in making IoT servers?

## Goals
The purpose of this research is to develop an IoT server application using the Golang programming language and the Fiber framework with the PostgreSQL database, then analyze the performance of the IoT server and compare the results with previous research by Hanin (2021) which uses the Falcon and Alvin (2023) framework which uses Sanic framework.


## Summary
In this research, an IoT data server application has been successfully developed using the Fiber framework with the Postgresql database. The development was carried out using the waterfall method for three iterations. The application that has been created can meet the defined user needs and succeed in the defined functional test. This research also added user interface features to the application. Applications with the Fiber framework have better performance than applications with the Sanic framework in Alvin's research (2023) and applications with the Falcon framework in Hanin's research (2021). Fiber applications have an average TPS of 2252.4, while Falcon has an average TPS of 30.5 and Sanic has an average TPS of 1106.8. Based on these TPS figures, the Fiber application has a performance of 73 times faster than the Falcon application and 2 times faster than the Sanic application. In addition, the Fiber application also has lower memory usage than the Falcon and Sanic applications. Fiber application has a memory usage percentage of 7.79%, Falcon 38.98%, and Sanic 12.85%. Based on these results, the Fiber framework has better performance and optimization of memory usage when used as an IoT back-end server compared to the Falcon and Sanic frameworks.

## Version
There are 2 version of this application
### Version 1
The data used is based on the research of Hanin (2021). He build the REST API IoT server based on ThingSpeak service. You can see it in the branch v1

### Version 1 Optimized
An Optimized application using Go-json library for JSON encoding and decoding. See this documentation from [Fiber](https://docs.gofiber.io/guide/faster-fiber/)

### Version 2
The data used is based on the modification of the previous version and the research of Alvin (2023). You can see it in the branch v2

## Paper
To be published

## Documentation
This API is documented using Postman, you can see the documentation here:
1. Version 1: https://documenter.getpostman.com/view/14947205/2s93CGSbrP
2. Version 2: https://documenter.getpostman.com/view/14947205/2s93RRxZh9

## Testing
The testing script can be found here:
1. Version 1: https://documenter.getpostman.com/view/14947205/2s93JzMLy5
2. Version 2: https://documenter.getpostman.com/view/14947205/2s93RRxZhA

You can run the test by cloning the Postman and run it yourself, you must first create the database. You can also run the server by using the test script by running 
```
./script/test.sh
```


## Running the application
1. Clone the repository
2. Make sure you have installed Golang > 1.19 
3. There are two ways to start the server

### In Development
1. Install [Golang Air](https://github.com/cosmtrek/air) library, this will auto build and run the server for you. 
2. Run `go install github.com/cosmtrek/air@latest`
3. Open the terminal, go to your project directory
4. type the command `air`, make sure the golang binary is added to your path variable. 
  a. In linux and mac, it is usually stored at /home/{username}/go/bin
5. The server is up and running, you can go to http://localhost:3000 for accessing the app


### In Production
1. This section of the tutorial assume you have a Linux Ubuntu for the deployment. Please adjust if you use different OS.
2. Make sure the script is excecutable 
```
chmod -R +x script/
```
3. Build the application `./script/build.sh`
4. Create a service file, [reference](https://stackoverflow.com/questions/58022141/pm2-like-process-management-solution-for-golang-applications)
5. Copy the service unit file from iot.service at this repository
  ```
[Unit]

[Install]
WantedBy=multi-user.target

[Service]
ExecStart=/root/iot-server/build/server-iot
WorkingDirectory=/root/iot-server
User=root
Restart=always
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=%n
  ```
6. Create this file at `/etc/systemd/system/iot.service`
7. Run `systemctl start iot.service`

#### Usual Operations
To have it always on when the machine starts:
```
systemctl enable iot.service
```

#### If you change your unit file after the first start or enable, you need to run:
```
systemctl daemon-reload
```

#### To see the status of the process, run:
```
systemctl status iot.service
```

#### To see the STDOUT of the process, run:
```
journalctl -f -u iot.service
```

For further help, please read the [manual page](https://www.freedesktop.org/software/systemd/man/systemd.unit.html).
