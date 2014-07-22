struct Range {
    1: required i32 StartUid;
    2: required i32 EndUid,
}

struct UsrInfo {
    1: required i32 Uid;
    2: required i32 Weight,
}

service tagSearchService {
    # getRange: 输入标签ID，返回拥有该标签的用户ID范围，用来做下面的查询
    Range getRange(1: i32 tagId)

    # getUsrs: 输入标签ID，查询的用户范围，以及最多查询该范围内的用户数量
    # 将有序返回该范围内的用户，list最大size为limitSize
    list<UsrInfo> getUsrs(1: i32 tagId, 2: Range r, 3: i32 limitSize)
}
