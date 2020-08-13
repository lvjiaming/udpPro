import {EventManager} from "./EventManager";
import {TipMgr} from "../Common/TipMgr";
export class GameEventManager extends EventManager{
    public gameSocket = null;  // 连接
    protected _hostStr = null; // 连接地址
    protected _isSelfClose=  null;  // 是否是自己主动断开连接

    protected _reconnectCurTime = null;  // 当前重连的次数
    protected _reconnectMaxTime = null; // 重连的最大次数
    protected _reconnectTimes = null; // 重连的时间间隔
    protected _reconnectTimer = null;  // 重连的定时器

    protected _curListen = null;  // 当前监听的对象（会发送一系列回调函数）

    protected _isLock = false; // 是否上锁（协议阻塞）
    protected _eventCache = null; // 阻塞的消息列表

    protected _name: null;

    constructor() {
        this.gameSocket = null;
        this._reconnectCurTime = 0;
        this._reconnectMaxTime = 10;
        this._reconnectTimes = 1000;
        this._isSelfClose = false;
        this._reconnectTimer = null;
        this._eventCache = [];
        this._isLock = false;
        super.constructor();
    }

    /**
     *  监听对象（会发送一系列特殊回调）
     *  reconnectFail  // 重连失败
     *  reconnectSuc  // 重连成功
     * @param target
     */
    public setListen(target): void {
        this._curListen = target;
    }

    /**
     *  设置名字
     */
    public setName(str): void {
        this._name = str;
    }
    /**
     *  移除监听
     */
    public removeListen(): void {
        this._curListen = null;
    }
    /**
     *  连接服务器，已经监听服务器一系列事件
     */
    public connect(hostStr, callBack): void {
        this._hostStr = hostStr;
        const self = this;
        cc.log("连接服务器：", this._hostStr);
        this.gameSocket = new WebSocket(hostStr);
        this.gameSocket.onopen = () => {
            console.log(`websocket has connect`);
            if (callBack && callBack instanceof Function) {
                callBack();
            }
        };
        this.gameSocket.onerror = () => {
            console.log(`websocket connect error`);
            // this.reconnect();
        };
        this.gameSocket.onclose = () => {
            console.log(`websocket has close`);
            if (this._isSelfClose) {
                console.log(`玩家主动断开连接，不重连`);
            } else {
                this.reconnect();
            }
        };
        this.gameSocket.onmessage = function (data) {
            // data = JSON.parse(data.data);
            // cc.log("this._isLock: ", self._isLock);
            // if (self._isLock) {
            //     self._eventCache.push({msgId: data.msgId, msgData: data.msgData});
            // } else {
            //     self.onMsg(data.msgId, data.msgData);
            // }
            // return;
            //  todo 以下是用protobuf传输数据写法
            if (cc.sys.isNative) {
                self.handleData(data.data);
            } else {
                const fileReader = new FileReader();  //  在浏览器中读取文件
                fileReader.onload = function (progressEvent) {  //  读取文件完成后触发（成功读取）
                    const utfs = this.result;  //  result就是读取的结果
                    self.handleData(utfs);
                };
                fileReader.readAsArrayBuffer(data.data);
            }
        };
        this.gameSocket.sendMessage = (data) => {
            if (this.gameSocket.readyState === WebSocket.OPEN) {
                this.gameSocket.send(data);
            } else {
                this.reconnect();
            }
        };
    }
    public reconnect(): void {
        this._reconnectCurTime ++;
        if (this._reconnectCurTime > this._reconnectMaxTime) {
            console.log(`重连次数已达最大`);
            TipMgr.getInstance().hide();
            if (this._curListen && this._curListen.reconnectFail) {
                this._curListen.reconnectFail();
            }
            return;
        }
        console.log(`正在进行第次${this._reconnectCurTime}重连`);
        TipMgr.getInstance().show(`正在进行第次${this._reconnectCurTime}重连`);
        this._reconnectTimer = setTimeout(() => {
            this.connect(this._hostStr, () => {
                console.log(`已重连`);
                TipMgr.getInstance().hide();
                TipMgr.getInstance().show(`重连成功`, 2);
                clearTimeout(this._reconnectTimer);
                this._reconnectCurTime = 0;
                if (this._curListen && this._curListen.reconnectSuc) {
                    this._curListen.reconnectSuc();
                }
            });
        }, this._reconnectTimes);
    }
    /**
     *  发送消息给服务端
     * @param msgId 消息的id
     * @param msgData 消息的数据
     */
    public sendMessage(msgId: number, msgData: any): void {
        // if (msgData === null || msgData === undefined) {
        //     msgData = null;
        // }
        // this.gameSocket.sendMessage(msgId, msgData);
        const body = msgData.serializeBinary();
        const uint8 = new Uint8Array(body.length + 1);
        body.forEach((item, index) => {
            uint8[index + 1] = item;
        });
        uint8[0] = msgId;
        // cc.log(uint8);
        this.gameSocket.sendMessage(uint8);
    }
    /**
     *  关闭与服务器的连接
     */
    public close(): void {
        this.gameSocket.close();
        this._isSelfClose = true;
    }
    /**
     *  处理数据（反序列化以及转化）
     * @param data 数据
     */
    public handleData(data): void {
        const self = this;
        const bytes = new Uint8Array(data);  // 转化数据
        let msgId = bytes[0];  //  协议id放在uint8Array的第一位
        const body = new Uint8Array(data, 1, data.byteLength - 1);
        if (this._isLock) {
            this._eventCache.push({msgId: msgId, msgData: body});
            cc.log("消息阻塞中");
        } else {
            this.onMsg(msgId, body);
        }
    }

    /**
     *  设置协议锁状态
     */
    public setEventLockState(state): void {
        cc.log("协议锁开关: ", state);
        this._isLock = state;
        if (!state) {
            this.openCache();
        }
    }

    /**
     *  放开缓存
     */
    public openCache(): void {
        for (let index = this._eventCache.length - 1; index >= 0; index--) {
            if (this._isLock) {
                break;
            } else {
                const data = this._eventCache[index];
                this.onMsg(data.msgId, data.msgData);
                this._eventCache.splice(index, 1);
            }
        }
    }

    public onMsg(msgId, body): void {

    }
}