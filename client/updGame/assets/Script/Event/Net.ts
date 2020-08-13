import {GameEventManager} from "./GameEventManager";
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

        }
    }

    public onMsg(msgId, body): void {
        cc.log("协议id: ", msgId);
        // cc.log(body);
        let data = null;
        switch (msgId) {

        }
        console.log(data);
        this.notifyEvent(msgId, data);
    }
}
