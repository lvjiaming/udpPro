export class Utils {
    /**
     *  根据时间搓，获取时间
     * @param timeTamp
     * @returns {string}
     */
    public static getTime(timeTamp: number): string {
        const date = new Date(timeTamp * 1000);
        const o = {
            "y": date.getFullYear(),
            "M": date.getMonth()+1,
            "d": date.getDate(),
            "h": date.getHours(),
            "m": date.getMinutes(),
            "s": date.getSeconds(),
        };
        o.M = (o.M<10)?('0'+o.M):o.M;
        o.d = (o.d<10)?('0'+o.d):o.d;
        o.h = (o.h<10)?('0'+o.h):o.h;
        o.m = (o.m<10)?('0'+o.m):o.m;
        o.s = (o.s<10)?('0'+o.s):o.s;
        return `${o.y}-${o.M}-${o.d} ${o.h}:${o.m}:${o.s}`
    }

    /**
     *  根据时间，获取时间戳
     */
    public static getTimeTamp(time: string): number {
        let timeTamp = Date.parse(new Date(time));
        timeTamp = timeTamp / 1000;
        return timeTamp;
    }

    /**
     *  获取当前时间
     * @returns {string}
     */
    public static getCurTime(): string {
        const date = new Date();
        const o = {
            "y": date.getFullYear(),
            "M": date.getMonth()+1,
            "d": date.getDate(),
            "h": date.getHours(),
            "m": date.getMinutes(),
            "s": date.getSeconds(),
        };
        cc.log(o);
        o.M = (o.M<10)?('0'+o.M):o.M;
        o.d = (o.d<10)?('0'+o.d):o.d;
        o.h = (o.h<10)?('0'+o.h):o.h;
        o.m = (o.m<10)?('0'+o.m):o.m;
        o.s = (o.s<10)?('0'+o.s):o.s;
        return `${o.y}-${o.M}-${o.d} ${o.h}:${o.m}:${o.s}`
    }

    /**
     *  改变钱
     * @param val
     */
    public static changeVal(val: number): number {
        let curVal = val || 0;
        curVal = Math.ceil(curVal * 100) / 100;
        return curVal;
    }

    public static setScreenFit(canvas_node) {
        const winSize: cc.Size = cc.winSize;
        const frameSize: cc.Size = cc.view.getFrameSize();
        const visibleSize = cc.view.getVisibleSize();
        const visibleOrigin = cc.view.getVisibleOrigin();

        const realWH = frameSize.width/frameSize.height;

        if ( realWH < 1.5 ) { //iPad=1.333, iphone4=1.5: 不要拉伸缩放(留黑边)
            canvas_node.getComponent(cc.Canvas).fitWidth = true;
        } else { //其它分辨率: 需拉伸缩放填满屏幕
            canvas_node.scaleY = visibleSize.height/winSize.height;
            canvas_node.scaleX = visibleSize.width/winSize.width;
        }
    }
}