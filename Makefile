# Generate Swagger Specification documentation in YAML format
swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models
		
swagger-serve:
	GO111MODULE=off swagger serve -F=swagger ./swagger.yaml