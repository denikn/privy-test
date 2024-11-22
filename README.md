# Privy Backend Engineer Test
[![Go Reference](https://pkg.go.dev/badge/github.com/andhikayuana/qiscus-unofficial-go.svg)](https://pkg.go.dev/github.com/syahidfrd/go-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

### Directory Structure
```
├── bin                                         // Directory to save binary file
├── cache                                       // Directory to save cached image data
├── cmd
│   └── main.go                                 // Main applications for this project
├── internal
│   ├── adapter
│   │   ├── http
│   │   │   └── handler.go                      // Encode raw body, query params, make a response to client
│   │   └── storage
│   │       └── filestorage.go                   // Handle saving file and getting path of file
│   ├── application
│   │   └── service.go                          // Contains all of business logic, mostly for image handler
│   ├── domain
│   │   └── model.go                            // Modelling the data
│   ├── middleware
│   │   └── jwt.go                              // Handle authorization with jwt
│   └── port
│       └── interfaces.go                       // Interface all of useful function
├── scripts
│   ├──  face_detection.py                      // Face detection handler
│   └──  haarcascade_frontalface_default.xml    // File to process image library
├── uploads                                     // Directory to save source and processed image
├── Makefile                                     // Set of command
└── .env                                        // Configuration app
```

### Prerequisite and full list what has been used
* [OpenCV](https://github.com/opencv/opencv) - OpenCV
* [golang-jwt](https://github.com/golang-jwt) - Library JWT for Golang
* [Gorilla Mux](https://github.com/gorilla/mux) - HTTP request multiplexer
* [godotenv](https://github.com/joho/godotenv) - Environment variables loader


### Configuration and running program
Copy .env.example to .env
```
cp .env.example .env
```

Fill .env
```
JWT_SECRET_KEY=your_secret
ALLOWED_USERNAME=deni
```

Run unit test
```
Make test
```

Run locally
```
Make run
```

Test then run
```
Make testrun
```

Build
```
Make build
```

### Documentation
<a href="https://documenter.getpostman.com/view/3134681/2sA3duGDLV" target="_blank">
Privy Backend Engineer Mini Project Documentation
</a>