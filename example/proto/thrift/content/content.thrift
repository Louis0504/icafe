namespace py content_thrift.content
namespace go content_thrift.content
namespace java content.thrift.content

include "exception.thrift"
include "base.thrift"

struct GetContentParam {
    1: required i64 content_id;
}


struct GetContentResponse {
    1: required base.Content content;
}


service ContentService {
    GetContentResponse get_content(1:GetContentParam param)
}
