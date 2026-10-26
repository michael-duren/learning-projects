CC=gcc

run:
	rm -r ./bin
	mkdir -p bin
	$(CC) -Wall -Wextra -g \
		-fsanitize=address,undefined \
		main.c request.c handler.c response.c netutils.c \
		-o bin/http
	./bin/http 8081
