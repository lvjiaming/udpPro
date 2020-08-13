export class EventManager {
    protected _observerList = null;
    constructor() {
        this._observerList = [];
    }

    /**
     *  添加观察者
     * @param target 观察的对象
     */
    public addObserver(target): void {
        this._observerList.forEach((item) => {
            if (item === target) {
                return true;
            }
        });
        if (target) {
            this._observerList.push(target);
        } else {
            console.log(`target is null`)
        }
        console.log(`this.ObserverList.length : ${this._observerList.length}`);
    }

    /**
     *  移除观察者
     * @param target 要移除观察的对象
     */
    public removeObserver(target): void {
        this._observerList.forEach((item, index) => {
            if (item === target) {
                this._observerList.splice(index, 1);
            }
        });
        console.log(`this.ObserverList.length : ${this._observerList.length}`);
    }
    /**
     *  移除所有的观察者
     */
    public removeAllObserver(): void {
        this._observerList = [];
    }
    /**
     *  通知观察者，有地方更改了
     * @param event 事件
     * @param msg 数据
     */
    public notifyEvent(event, msg): void {
        try {
            this._observerList.forEach((item, index) => {
                item.onEventMessage(event, msg);
            });
        } catch (err) {
            console.error(`抛出异常：${err}`);
        }
    }
}