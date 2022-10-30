namespace py content_thrift.exception
namespace go content_thrift.exception
namespace java com.content.thrift.exception

exception Error {
    1: required i32 code;   //error no， 整数
    2: required string name;    //error name 错误英文名， 比如 NotFound
    3: required string message; //错误信息，比如可以给用户展示用。
}