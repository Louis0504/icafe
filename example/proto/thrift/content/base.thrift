namespace py content_thrift.base
namespace go content_thrift.base
namespace java com.content.thrift.base


// 内容标准协议
struct Content {
    1: required i64     id;
    2: required string  name;
    3: required string  content;
    7: required string     created_at;  //创建时间
    8: required string     updated_at;  //更新时间
}

