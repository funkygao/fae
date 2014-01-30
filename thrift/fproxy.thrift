namespace go fproxy
namespace php fproxy

service Fproxy { 
    list<string> funCall(1:i64 callTime, 2:string funCode, 3:map<string, string> paramMap)
} 
