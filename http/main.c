#include <arpa/inet.h>
#include <errno.h>
#include <netdb.h>
#include <netinet/in.h>
#include <pthread.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define MAX_CONN 10

const char *parse_port(int argc, char *argv[]) {
	if (argc == 2) {
		int n = atoi(argv[1]);
		if (n > 0 && n <= 65535) {
			return argv[1];
		}
		printf("invalid port: %s, falling back to :8080\n", argv[1]);
	}
	return "8080";
}

// TODO: graceful shutdown
// static void __shutdown(void) {
// 	printf("waiting for response writers to finish\n");
// 	printf("shutting down\n");
// }

// void sig_handler(int signo) {
// 	if (signo == SIGINT) {
// 		__shutdown();
// 	}
// }

#define MAX_REQ_SIZE (1 << 20)

// caller frees
char *read_req(int fd) {
	size_t cap = 4096, total = 0;
	char *buf = malloc(cap);
	if (buf == NULL) {
		return NULL;
	}

	while (1) {
		if (total == cap - 1) {
			if (cap >= MAX_REQ_SIZE)
				break;
			char *tmp = realloc(buf, cap * 2);
			if (tmp == NULL)
				break;
			buf = tmp;
			cap *= 2;
		}

		ssize_t n = recv(fd, buf + total, cap - 1 - total, 0);
		if (n < 0) {
			if (errno == EINTR)
				continue;
			perror("recv");
			free(buf);
			return NULL;
		}
		if (n == 0)
			break;

		total += (size_t)n;
		buf[total] = '\0';

		if (strstr(buf, "\r\n\r\n"))
			break;
	}

	buf[total] = '\0';
	return buf;
}

static void *handle(void *arg) {
	intptr_t fd = (intptr_t)arg;
	char *req = read_req(fd);
	printf("%s\n", req);

	const char *res = "HTTP/1.1 200 OK\r\n\r\nHello World!\n";
	ssize_t n = send(fd, res, strlen(res), 0);
	if (n != (long)strlen(res)) {
		printf("did not send expected amount\n");
	}

	close(fd);
	return NULL;
}

#define STACK_SIZE 0x100000

int start_handler(int fd) {
	pthread_t p;
	pthread_attr_t attr;
	int rc = pthread_attr_init(&attr);
	if (rc != 0) {
		perror("pthread_attr_init");
		return rc;
	}
	rc = pthread_attr_setstacksize(&attr, STACK_SIZE);
	if (rc != 0) {
		perror("pthread_attr_setstacksize");
		return rc;
	}
	rc = pthread_create(&p, &attr, &handle, (void *)(intptr_t)fd);
	if (rc != 0) {
		perror("pthread_create");
		return rc;
	}
	return 0;
}

void *get_in_addr(struct sockaddr *sa) {
	if (sa->sa_family == AF_INET) {
		return &(((struct sockaddr_in *)sa)->sin_addr);
	}

	return &(((struct sockaddr_in6 *)sa)->sin6_addr);
}

int main(int argc, char *argv[]) {
	// setup server
	// TODO: re add
	// if (signal(SIGINT, sig_handler) == SIG_ERR) {
	// 	perror("signal");
	// 	return 1;
	// }

	char addr_str[INET6_ADDRSTRLEN];

	const char *port = parse_port(argc, argv);

	int rc;
	struct addrinfo hints;
	struct addrinfo *servinfo;

	memset(&hints, 0, sizeof hints);

	hints.ai_family = AF_INET;		 // use ipv4 for now
	hints.ai_socktype = SOCK_STREAM; // tcp
	hints.ai_flags = AI_PASSIVE;	 // fill in ip

	if ((rc = getaddrinfo(NULL, port, &hints, &servinfo)) != 0) {
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
	inet_ntop(servinfo->ai_family, get_in_addr(servinfo->ai_addr), addr_str,
			  sizeof addr_str);

	printf("listening to incoming connnections at: %s:%s\n", addr_str, port);
	struct sockaddr_storage in_addr;
	socklen_t in_size = sizeof(in_addr);

	// TODO: need to handle open fds when ctrl+c happens
	while (1) {
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

		if ((rc = start_handler(cn)) != 0) {
			printf("start handler failed with %d rc", rc);
		}
	}

	// __shutdown();

	return 0;
}
