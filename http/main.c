#include <arpa/inet.h>
#include <errno.h>
#include <fcntl.h>
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

// ~ 1MB
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

#define PACKET_SIZE 4096

int write_res(int status, char *reason, char *body, char *buf) {
	int n = snprintf(buf, strlen(buf),
					 "HTTP/1.1 %d %s\r\n"
					 "Content-Type: text/plain\r\n"
					 "Content-Length: %zu\r\n"
					 "Connection: close\r\n"
					 "\r\n"
					 "%s",
					 status, reason, strlen(body), body);

	if (n < 0 || (size_t)n >= sizeof buf)
		return -1;
	return n;
}

// closes socket descriptor
int send_res(int sd, int status, char *reason, char *body, size_t len) {
	char res[MAX_REQ_SIZE];
	write_res(status, reason, body, res);

	size_t written = 0;
	while (written < len) {
		size_t cap = PACKET_SIZE;
		if (len < cap) {
			cap = len;
		}

		ssize_t sent = send(sd, res + written, cap, 0);
		written += sent;
		len -= sent;
		if (!sent)
			goto error;
		if (sent == 0) {
			goto exit;
		}
	}

	goto exit;
error:
	close(sd);
	perror("send");
	return -1;
exit:
	return 0;
}

typedef struct __req_t {
	const char *path;
	const char *method;
	const char *protocol;
	char **headers;
	int header_count;
} req_t;

// GET / HTTP/1.1
// Host: localhost:8081
// User-Agent: curl/8.21.0
// Accept: */*
static int parse_req(char *req_str, req_t *req) {
	char *p = req_str;
	req->method = p;
	p = strchr(req_str, ' ');
	if (!p)
		return -1;
	*p++ = '\0';

	req->path = p;
	// end of path
	p = strchr(p, ' ');
	if (!p)
		return -1;
	*p++ = '\0';

	req->protocol = p;
	p = strstr(p, "\r\n");
	if (!p)
		return -1;
	*p = '\0';

	// TODO: parse headers

	return 0;
}

static void free_req(req_t *req) { free(req->headers); }

#define MAX_PATH 4096

static void *handle(void *arg) {
	intptr_t fd = (intptr_t)arg;
	char *req_str = read_req(fd);

	req_t req;
	int rc = parse_req(req_str, &req);
	if (rc != 0) {
		rc = send_res(fd, 400, "Malformed Request", NULL, 0);
		if (rc != 0) {
			puts("unable to send 500 response");
		}
		goto clean;
	}
	char path[MAX_PATH];
	if (strcmp(req.path, "/") == 0) {
		strcpy(path, "./static/index.html");
	} else {
		strcpy(path, "./static");
		strcpy(path + 9, req.path);
	}
	printf("opening path: %s\n", path);
	int f = open(path, O_RDONLY);
	if (f == -1) {
		rc = send_res(fd, 404, "Not Found", NULL, 0);
		goto clean;
	}
	char buf[MAX_REQ_SIZE];
	int total_read = 0;

	while (1) {
		int n = read(f, &buf + total_read, MAX_PATH);
		if (n == -1)
			goto clean;
		total_read += n;
		if (n == 0 || total_read >= MAX_REQ_SIZE) {
			break;
		}
	}

	free_req(&req);
	rc = send_res(fd, 200, "OK", buf, strlen(buf));
	if (rc != 0) {
		puts("unable to send resonse");
	}
	return NULL;

clean:
	free(req_str);
	close(fd);
	if (f > 0) {
		close(f);
	}
	// TODO: need to send a result of the handler thread
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
