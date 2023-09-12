./hls-cmder.exe download --type m3u8 \
    --m3u8-path "<the url of .m3u8>" \
    --name "<filename>.ts" \
    --worker 2 \
    --origin "<origin of http header>" \
    --referer "<referer of http header>" \
    --directly true