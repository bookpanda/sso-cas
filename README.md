# SSO: Central Authentication Service

## What is SSO
Single Sign-On (SSO) is a user authentication process that allows a user to access multiple applications with one set of login credentials. This is a common practice in the enterprise, where a user can log in to their computer and, with a single login, gain access to all of the applications they need to do their job.

This project is a very simple implementation of SSO.
## Services involed
### CAS (Central Authentication Service)
This is the main service that handles the authentication of the user. It is responsible for generating the token and validating the token. It also has a user management system.
- TGT (Ticket Granting Ticket)

## Setting up
### Stack
-   golang
-   postgresql
-   react

### Prerequisites
-   💻
-   golang 1.22 or [later](https://go.dev)
-   docker
-   makefile
-   [Go Air](https://github.com/air-verse/air)

### Installation + Running
1. Clone this repo
2. Run `go mod download` to download all the dependencies.
3. Copy `.env.template` and paste it in the same directory as `.env`. Fill in the appropriate values for both `frontend` and `backend`.
4. Run `air` in `backend` and `pnpm dev` in `frontend`.