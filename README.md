# IE-Backend-project
In this project, we are going to develop a service with Go programming language to provide monitoring of http endpoints in adjustable intervals (for example, every 30 seconds, 1 minute, 5 minutes). This service sends an http request to the endpoint and logs the status code of the response. Each URL must have an error threshold that indicates the maximum number of errors tolerable, after which our service must generate an alert to the user to whom the URL belongs. A successful http call is indicated by a status code of 2xx, and an unsuccessful call is indicated by a status code of other than 2xx.