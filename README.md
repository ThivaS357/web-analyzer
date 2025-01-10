# 🌐 Web Page Analyzer

A **Golang-powered web application** that analyzes web pages and provides insights like HTML version, title, heading counts, internal and external links, broken links, and login forms.

---

## 🚀 **Project Overview**

This application accepts a URL, analyzes the web page content, and returns key information:
- **HTML Version:** Detects the HTML version used.
- **Page Title:** Displays the web page's title.
- **Headings:** Counts all heading tags (`h1` to `h6`).
- **Links:** Differentiates between internal and external links and identifies broken ones.
- **Login Form:** Detects login-related forms (`login`, `signin`, etc.).

---

## 🛠️ **Key Components**
- **main.go:** Handles server routing and request handling using Gin.
- **analyzer/analyzer.go:** Core logic for webpage analysis using net/http.
- **templates/index.html:** Frontend HTML page for input and result display.

---


## 📦 **Prerequisites**

Make sure you have the following installed:
- [Golang](https://golang.org/) (version 1.22 or higher)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) (for containerization)
- [Docker Compose](https://docs.docker.com/compose/) (for managing multi-container Docker applications)

---

## 💻 **Setup Instructions**

1. **Clone the Repository:**
    ```bash
    git clone https://github.com/ThivaS357/web-analyzer.git
    cd web-analyzer
    ```

2. **Install Dependencies:**
    ```bash
    go mod tidy
    ```

3. **Run the Application:**
    ```bash
    go run main.go
    ```
4. **Run the Application using Docker:**
    ```bash
    docker build -t web-analyzer .
    docker run -p 8080:8080 web-analyzer
    ```

5. **Run the Application using Docker Compose:**
    ```bash
    docker-compose up --build
    ```

6. **Access the Web Interface:**  
    Open your browser and navigate to:  
    ```
    http://localhost:8080
    ```

---



## 🧪 **Running Tests**

### Note: This project size is small and Logic are tighly coupled (Only for 1 API). So, test coverage only addedd to /analyze API. As complexity grows unit test will be added to each individual logic


### Run to test and create coverage
go test ./... -coverprofile=coverage.out

### Open and see the tect coverage
go test ./... -coverprofile=coverage.out


## Note:
We need to keep HTML traversal in a one thread (Go Routine). But checking external linkes and validating broken links parts can be added to other threads (Go Routine)