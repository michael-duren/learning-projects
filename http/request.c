#include "request.h"
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

#define MAX_REQ_SIZE (1 << 20)

char *read_from_socket(int sd) {
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

		ssize_t n = recv(sd, buf + total, cap - 1 - total, 0);
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

// must call to free backing buffer
void free_req(req_t *req) { free(req->req_str); }

int read_req(int sd, req_t *req) {
	char *req_str = read_from_socket(sd);
	if (req_str == NULL)
		return -1;

	return parse_req(req_str, req);
}
