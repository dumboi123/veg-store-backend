# 🚀 How to Start the Application

> ⚠️ **Notice:** Run all the following commands **from the project root directory**.

## 1. 🧩 Prerequisites
Ensure the following packages are installed:
- **Go** ≥ 1.25.3
- **Docker** (Builder, Engine)
- **Git**
- **make**

---

## 2. ▶️ Initial Setup
```bash
make prepare # Only run for the first time
make start
```
> ⚠️ The first run may take a while — please wait until setup completes.

---

## 3. 🐳 Run Application in Docker

### Run normally (no hot reload)
```bash
make run
```

### Run with Hot Reload
```bash
make run-dev
```
> 💡 If hot reload doesn’t trigger after code changes, stop it using <kbd>Ctrl</kbd> + <kbd>C</kbd> and rerun the command.

---

## 4. ✅ Verify the Application
Visit:  
👉 [http://localhost:2345/api/v2/swagger-ui/index.html#](http://localhost:2345/api/v2/swagger-ui/index.html#)

---

## 5. 🛑 Stop Application
To stop the containers, run:
```bash
make stop
```

---

# ⚙️ Update Environment Variables

1. Update `scripts/helper/env_config.sh`
2. Update `create_env_file()` in `scripts/helper/functions.sh`
3. Update `environment` section of `veg-store-backend` service in `docker/docker-compose.dev.yml`

After changes, restart the environment:
```bash
make restart
```
> 🧠 If you only changed `.env` values, updating (3) and restarting is enough.

---

# 📘 Update Swagger Schemas
To regenerate Swagger documentation:
```bash
make swagger
```
> 🔄 Re-run the app if hot reload is not enabled.

---

# 🧪 Testing Guide

## ▶️ Run All Unit Tests with Coverage
To execute all unit tests and generate a detailed coverage report:
```bash
make coverage
```
> 📄 After running, open ./test/report/index.html in your browser to view the full coverage report.

## 🎯 Run Tests in a Specific Package:
Use the PKG argument to target a specific package:
```bash
# make test PKG=./test/unit/<package-you-want>
# Example:
make test PKG=./test/unit/handler/rest_test
```

## 🧩 Run a Single Test Function:
Use the TEST argument to execute one specific test:
```bash
# make test-one TEST=<TestSuiteName>/<TestName>
# Example:
make test-one TEST=TestUserHandler/TestHello_success
```