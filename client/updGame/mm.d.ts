/**
 * Created by Administrator on 2020/8/24.
 */
declare namespace msgPb {
    export var Event = {
        EVENT_MSG_INFO: 0,
        EVENT_REGISTER_REQ: 1,
        EVENT_REGISTER_REP: 2,
        EVENT_LOGIN_REQ: 3,
        EVENT_LOGIN_REP: 4,
        EVENT_ADD_ONE_INFO: 5,
        EVENT_QUERY_INFO_REQ: 6,
        EVENT_RETURN_INFO_LIST: 7,
        EVENT_CHANGE_INFO_REQ: 8,
        EVENT_CHANGE_INFO_REP: 9,
        EVENT_DEL_INFO_REQ: 10,
        EVENT_DEL_INFO_REP: 11,
        EVENT_STATISYICAL_INFO_CHANGE: 12,
    };
    export var CodeType = {
        ERR: 0,
        SUC: 1
    };
    class Code {
        setCode(...args);
        getCode(): number;
        setMsg(...args);
        getMsg(): string;
    }
}