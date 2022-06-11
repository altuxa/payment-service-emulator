build: 
	docker build -t paymentsys .
run:
	docker run -d -p 8080:8080 --name myproject paymentsys
stop: 
	docker stop myproject
delete: 
	docker rm myproject