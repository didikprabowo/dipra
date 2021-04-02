push:
	@docker build -t 207979/example-dipra .
	@docker push 207979/example-dipra
	@docker container create --name dipra-example 207979/example-dipra
	@docker container start dipra-example
	@docker container logs -f dipra-example
	@docker container stop dipra-example
	@docker container rm dipra-example