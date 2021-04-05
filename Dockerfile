FROM alpine
ADD hash_check /bin/
ENTRYPOINT /bin/hash_check