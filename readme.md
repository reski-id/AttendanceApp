

<h3 align="center">Attendance Rest API <br>
<h5 align="center" >Golang Echo<h5>
<br>
</h4>
<p align="left">
<h2>
  Content <br></h2>
  • Key Features <br>
  • Installing Using Github<br>
  • End Point<br>
  • Technologi that i use<br>
  • Contact me<br>
</p>

## 📱 Features

* Auth
* Employee
* Clock In
* Clock Out
* Swagger OpenAPI

## ⚙️ Installing and Runing from Github

installing and running the app from github repository <br>
To clone and run this application, you'll need [Git](https://git-scm.com) and [Golang](https://go.dev/dl/) installed on your computer. From your command line:

```bash
# Clone this repository
$ git clone https://github.com/reski-id/AttendanceApp.git

# Go into the repository
$ cd AttendanceApp

# Install dependencies
$ go get

# Run the app
$ go run main.go

# if you have problem while running you can use bash cmd and type this..
$ source .env #then type 
$ go run main.go #again
```

> **Note**
> Make sure you allready create database mysql `attendancedb` for this app.more info in local `.env` and `utils/database.go` file.


## 📜 End Point  

Auth
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `POST`        | /api/v1/register            | Register
| `POST`        | /api/v1/login         | Login

Employee
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `GET`         | /api/v1/employees             | Get all employees      
| `GET`         | /api/v1/employees/:id          | Get One employees      
| `GET`         | /api/v1//employees/search       | Searching a employees      
| `POST`        | /api/v1/employees              | Insert employees 
| `PUT`         | /api/v1/employees/:id         | Update data employees
| `DELETE`      | /api/v1/employees/:id         | Delete employees  

Attendance
| Methode       | End Point      | used for            
| ------------- | -------------  | -----------                  
| `POST`        | /attendance/clock-in/:id             | Clock IN
| `POST`        | /attendance/clock-out/:id             | Clock OUT



## 📜 Swagger Open Api
after you running the app you can access swagger open api in this 
 http://localhost:8080/swagger/index.html

## 📜 Postman 
you can find postman testing in  `/screenshoot/` folder

## 🛠️ Technology

This software uses the following Tech:

- [Golang](https://go.dev/dl/)
- [Echo Framework](https://echo.labstack.com/)
- [Gorm](https://gorm.io/index.html)
- [OpenAPI Swaggo](https://github.com/swaggo/gin-swagger)
- [mysql](https://www.mysql.com/)
- [Linux](https://www.linux.com/)
- [Docker](https://www.docker.com/)
- [Dockerhub](https://hub.docker.com/u/programmerreski)
- [Git Repository](https://github.com/reski-id)
- [Trunk Base Development](https://trunkbaseddevelopment.com/)


## 📱 Contact me
feel free to contact me ... 
- Email programmer.reski@gmail.com 
- [Linkedin](https://www.linkedin.com/in/reski-id)
- [Github](https://github.com/reski-id)
- Whatsapp <a href="https://wa.me/+6281261478432?text=Hello">Send WhatsApp Message</a>
