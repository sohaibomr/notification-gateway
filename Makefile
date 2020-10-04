check:
	cd notifications-forwarder && make check
	cd gateway && make check
	cd group-notifier && make check
	cd user-notifier && make check

.PHONY: tidy
tidy:
	cd notifications-forwarder && make tidy
	cd gateway && make tidy
	cd group-notifier && make tidy
	cd user-notifier && make tidy

build:
	$(MAKE) check
	docker-compose build

run:
	docker-compose up