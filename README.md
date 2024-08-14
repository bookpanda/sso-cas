# SSO: Central Authentication Service

## Mechanism

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