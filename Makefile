pull-latest-mac:
	docker pull --platform linux/x86_64 ghcr.io/bookpanda/sso-cas:latest
	docker pull --platform linux/x86_64 ghcr.io/bookpanda/sso-sample-service:latest

docker-qa:
	docker-compose -f docker-compose.qa.yml up
