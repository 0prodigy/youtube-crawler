# Youtube video scrapper 

> This is a simple web scraper that scrapes youtube videos from a query and expose view api for the same.


## Built With

- GO
- MYSQL
- Fiber
- GORM


## Getting Started

**Just fork it and fell free to use it.**

To get a local copy up and running follow these simple example steps.

### Prerequisites

- Docker

### Install

- create a .env file in the root directory and add the following variables

```
YOUTUBE_API_KEY 

DATABASE_URL = "user:password@tcp(host:port)/youtube?charset=utf8&parseTime=True&loc=Local"

```

- run the following command to start mysql docker container

```
docker run --name mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=database -p 3306:3306 -d mysql:latest

```
- run the following command to start the server

```
go run main.go

```

## Authors

👤 **Akash Pathak**

- Github: [@0prodigy](https://github.com/0prodigy)
- Twitter: [@pathakprodigy](https://twitter.com/pathakprodigy)
- Linkedin: [Akash Pathak](https://www.linkedin.com/in/akash-pathak-0796a7165)
- Email: (pathakvikash9211@gmail.com)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

Feel free to check the [issues page](https://github.com/0prodigy/apple-web-clone/issues).

## Show your support

Give a ⭐️ if you like this project!

## 📝 License

This project is [MIT](./LICENSE) licensed.
