import {GameEventManager} from "./GameEventManager";
import msgPb = require('../Proto/Common/msg_pb');
export class Net extends GameEventManager{
    public static net: Net = null;
    public static getInstance(): Net {
        if (this.net == null) {
            this.net = new Net();
        }
        return this.net;
    }

    public send(msgId, data): void {
        switch (msgId) {
            case msgPb.Event.EVENT_MSG_INFO: {
                this.sendMessage(msgId, data);
                break;
            }
        }
    }

    public onMsg(msgId, body): void {
        cc.log("协议id: ", msgId);
        cc.log(body);
        let data = null;
        switch (msgId) {
            case msgPb.Event.EVENT_MSG_INFO: {
                data = msgPb.Code.deserializeBinary(body);
                break;
            }
        }
        console.log(data);
        this.notifyEvent(msgId, data);
    }
}
