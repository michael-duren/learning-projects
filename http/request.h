#ifndef REQ
typedef struct __req_t {
	const char *path;
	const char *method;
	const char *protocol;
	char **headers;
	int header_count;
	char *req_str;
} req_t;

int read_req(int sd, req_t *req);
void free_req(req_t *req);
#endif // !REQ
