#include <arpa/inet.h>
#include <assert.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>

typedef struct __req {
	const char *host;
	const char *path;
	char *msg;
	int msg_len;
} req_t;

// write_request writes to the *req struct and returns 0
// if successful, 1 otherwise
int write_request(req_t *req, size_t cap) {
	const char *fmt = "GET %s HTTP/1.1\r\n"
					  "Host: %s\r\n"
					  "Connection: close\r\n"
					  "\r\n";

	int32_t len = snprintf(req->msg, cap, fmt, req->path, req->host);
	if (len < 0 || (size_t)len >= cap)
		return 1; // error or truncated
	req->msg_len = len;
	return 0;
}

int main(int argc, char *argv[]) {
	if (argc != 2) {
		puts("usage: c <url>");
		return 1;
	}
	const char *url = argv[1];

	struct addrinfo hints, *res;
	int status;

	memset(&hints, 0, sizeof hints);
	hints.ai_family = AF_UNSPEC;
	hints.ai_socktype = SOCK_STREAM;
	char ipstr[INET6_ADDRSTRLEN];

	if ((status = getaddrinfo(url, "http", &hints, &res)) != 0) {
		fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(status));
		return 2;
	}

	int s = socket(res->ai_family, res->ai_socktype, res->ai_protocol);
	puts("after opening socket");

	if (s == -1) {
		// TODO: clean up
		freeaddrinfo(res); // free the linked list
		perror("socket");
		return 2;
	}

	int rc = connect(s, res->ai_addr, res->ai_addrlen);
	if (rc != 0) {
		freeaddrinfo(res);
		close(s);
		perror("connect");
		return 2;
	}
	puts("after connect");

	req_t req;
	memset(&req, 0, sizeof(req_t));
	char msg[1024] = {0};
	req.msg = msg;
	req.host = url;
	req.path = "/";

	rc = write_request(&req, 1024);
	int bytes_sent = send(s, req.msg, req.msg_len, 0);
	if (bytes_sent == -1) {
		freeaddrinfo(res); // free the linked list
		perror("send");
		return 2;
	}
	printf("sent %d bytes\n", bytes_sent);

	// TODO: add null terminator to buf
	char buf[1024];
	// TODO: need to loop through since recv might not get the whole response in
	// one call
	int n = recv(s, buf, sizeof buf - 1, 0);
	if (n == -1) {
		perror("recv");
		shutdown(s, SHUT_RDWR);
		freeaddrinfo(res); // free the linked list
		return 2;
	}
	buf[n] = '\0';

	puts("after recv");
	puts(buf);

	shutdown(s, SHUT_RDWR);
	close(s);
	freeaddrinfo(res); // free the linked list

	return EXIT_SUCCESS;
}
