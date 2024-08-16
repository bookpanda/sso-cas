# SSO: Central Authentication Service

## What is SSO
Single Sign-On (SSO) is a user authentication process that allows a user to access multiple applications with one set of login credentials. This is a common practice in the enterprise, where a user can log in to their computer and, with a single login, gain access to all of the applications they need to do their job.

This project is a very simple implementation of SSO.
## Services
### CAS (Central Authentication Service)
This is the main service that handles the authentication of the user and also has a user management system. In real world, this would be something like the SSO part of [Okta](https://www.okta.com) or [HENNGE One](https://hennge.com/global/henngeone) service.
#### Components
#### 1. TGT (Ticket Granting Ticket) aka Session
ds
#### 2. ST (Service Ticket)
Every time a user uses a service, if they are already logged in to SSO, they need to prove it to that service. To do so, the user is redirected to CAS to get a new ST and then redirected back to that service. The user then presents the ST to the service so that the service can send the ST back to CAS.

#### 3. Database
The main database that stores the user information and sessions. Here are some tables:
-  `users`: Stores the main user information. For this project, other services don't duplicate data from this table, so they would have to make a request to this service to get the common user information e.g. email, names. For service-specific information, they would be stored in their services' databases.
-  `sessions`: Contains TGT token of each session, expiry time and payload e.g. email, user_id, roles, etc.
- `service_tickets`: Stores ST token, the session it corresponds to, the service it was issued to, user_id, and the expiry time.

#### 4. Google OAuth
Used for Google login. You can use any other OAuth provider like Facebook, Github, etc. or even your own OAuth provider.

### Sample Service (Consumers of CAS SSO)

## SSO Login Flow
1. The user attempts to access a resource or application that is protected by CAS e.g. internal tools.
1. User goes to CAS to login via Google OAuth. If they started from other services, they would be redirected to CAS to do Google login.

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