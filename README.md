<div align="center">
  <img width=200 alt="logo cpe-calendar" src="/static/favicon.svg">

  # CPE calendar

  Sync your CPE calendar with your personal one.

</div>

# About

The goal of the CPE Calendar is to offer students easy and effortless access to their school schedule by syncing it with their personal calendar. This works on all devices (phones, computers, laptops) and with any calendar provider, including Apple and Google.

The calendar automatically updates every hour, keeping you informed of any schedule changes. This project is open-source, and contributions or issue reports are welcome on GitHub.

# Setup

1. The first step is to set up your own `.env` file. Use `example.env` as a reference.
2. Then run the production version using Docker Compose.

# Known Issues

There can be an issue starting the Docker environment on Windows due to the missing `make-key.sh` script.

# Development

If you want to run the project without the Docker environment, follow these steps:

### Generate the needed key
```bash
openssl genrsa -out secret/private.pem 2048
openssl rsa -in secret/private.pem -pubout > static/public.pem
```

### Start the code
```bash
go mod download
go run main.go
```

# Affiliation

This project is entirely independent and is not affiliated with any school or organization.

