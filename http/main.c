#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define MAX_CONN 10

void parse_port(int argc, char *argv[], char *port) {
	if (argc < 2) {
		port = "8080";
	} else {
		int n;
		if ((n = atoi(port)) > 65535) {
			printf("invalid port: %d, falling back to :8080\n", n);
			port = "8080";
		} else {
			port = argv[1];
		}
	}
}
void *get_in_addr(struct sockaddr *sa) {
	if (sa->sa_family == AF_INET) {
		return &(((struct sockaddr_in *)sa)->sin_addr);
	}

	return &(((struct sockaddr_in6 *)sa)->sin6_addr);
}

int main(int argc, char *argv[]) {
	char port[8];
	char addr_str[INET6_ADDRSTRLEN];
	parse_port(argc, argv, port);

	int status;
	struct addrinfo hints;
	struct addrinfo *servinfo;

	memset(&hints, 0, sizeof hints);

	hints.ai_family = AF_INET;		 // use ipv4 for now
	hints.ai_socktype = SOCK_STREAM; // tcp
	hints.ai_flags = AI_PASSIVE;	 // fill in ip

	if ((status = getaddrinfo(NULL, "8080", &hints, &servinfo) != 0)) {
		perror("getaddrinfo");
		return 2;
	}

	// create socket using prev info
	int s = socket(servinfo->ai_family, servinfo->ai_socktype,
				   servinfo->ai_protocol);
	if (s == -1) {
		perror("socket");
		return 2;
	}
	if (bind(s, servinfo->ai_addr, servinfo->ai_addrlen) != 0) {
		perror("bind");
		close(s);
		return 2;
	}

	if (listen(s, MAX_CONN) != 0) {
		perror("bind");
		close(s);
		return 2;
	}

	printf("listening to incoming connnections at: %s\n",
		   servinfo->ai_addr->sa_data);
	struct sockaddr_storage in_addr;
	socklen_t in_size = sizeof(in_addr);
	// why is casting necessary ?
	int cn = accept(s, (struct sockaddr *)&in_addr, &in_size);
	if (cn == -1) {
		perror("accept");
		puts("couldn't accept connection");
		close(s);
		return 2;
	}

	inet_ntop(in_addr.ss_family, get_in_addr((struct sockaddr *)&in_addr),
			  addr_str, sizeof addr_str);
	printf("request from: %s\n", addr_str);

	const char *res = "HTTP/1.1 200 OK\r\n\r\nHello World!\n";
	ssize_t n = send(cn, res, strlen(res), 0);
	if (n != (long)strlen(res)) {
		perror("send");
		return 1;
	}

	close(cn);
	close(s);
	return 0;
}
