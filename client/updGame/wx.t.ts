/**
 * Created by Administrator on 2020/5/26.
 */
declare module wx {
    export class UDPSocket {
        bind(port?: number): number; //绑定一个系统随机分配的可用端口，或绑定一个指定的端口号
        close(): void; // 关闭 UDP Socket 实例，相当于销毁
        offClose(cb?: (...args) => void): void; // 取消监听关闭事件
        offError(cb?: (...args) => {}): void; // 取消监听错误事件
        offListening(cb?: (...args) => void): void; // 取消监听开始监听数据包消息的事件
        offMessage(cb?: (...args) => void): void; //取消监听收到消息的事件
        onClose(cb?: (...args) => void): void; // 监听关闭事件
        onError(cb?: (...args) => void): void; // 监听错误事件
        onListening(cb?: (...args) => void): void; // 监听开始监听数据包消息的事件
        onMessage(cb?: (...args) => void): void; // 监听收到消息的事件
        send(...args): void; // 监听收到消息的事件
    }
    export function createUDPSocket(): UDPSocket; // 创建一个 UDP Socket 实例
}