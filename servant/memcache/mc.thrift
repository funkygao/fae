namespace go mc
namespace php mc

enum MemcacheOpType {
    GET = 1,
    SET = 2,
    INC = 3,
    ADD
}

service McService { 
    list<string> funCall(1:i64 callTime, 2:string funCode, 3:map<string, string> paramMap),
    oneway void ping()
} 
