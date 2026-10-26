#include "request.h"
#include "response.h"
#include <arpa/inet.h>
#include <fcntl.h>
#include <netdb.h>
#include <netinet/in.h>
#include <pthread.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define MAX_PATH 4096

#define MAX_REQ_SIZE 16 * 1024

static void *handle(void *arg) {
	intptr_t fd = (intptr_t)arg;
	int f = -1;
	req_t req;
	int rc = read_req(fd, &req);
	if (rc != 0) {
		rc = send_res(fd, 400, "Malformed Request", NULL, 0);
		if (rc != 0) {
			puts("unable to send 500 response");
		}
		goto clean;
	}
	if (strstr(req.path, "..") != NULL) {
		send_res(fd, 403, "Forbidden", NULL, 0);
	}
	char path[MAX_PATH];
	if (strcmp(req.path, "/") == 0) {
		snprintf(path, sizeof(path), "./static/index.html");
	} else {
		snprintf(path, sizeof(path), "./static%s", req.path);
	}
	printf("opening path: %s\n", path);
	f = open(path, O_RDONLY);
	if (f == -1) {
		puts("error reading path");
		rc = send_res(fd, 404, "Not Found", NULL, 0);
		goto clean;
	}
	char buf[MAX_REQ_SIZE];
	int total_read = 0;

	while (total_read < MAX_REQ_SIZE) {
		int n = read(f, buf + total_read, MAX_PATH - total_read);
		if (n == -1)
			goto clean;
		if (n == 0)
			break;
		total_read += n;
	}

	rc = send_res(fd, 200, "OK", buf, strlen(buf));
	if (rc != 0) {
		puts("unable to send resonse");
	}
clean:
	free_req(&req);
	close(fd);
	if (f >= 0) {
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
