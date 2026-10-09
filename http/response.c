#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#define PACKET_SIZE 4096
#define MAX_RESPONSE_SIZE 16 * 1024

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
	char res[MAX_RESPONSE_SIZE];
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
