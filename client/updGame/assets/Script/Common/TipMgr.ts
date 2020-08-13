import {resLoad} from "./ResLoad";
export class TipMgr {
    public static tip = null;
    public static getInstance(): TipMgr {
        if (!this.tip) {
            this.tip = new TipMgr();
        }
        return this.tip;
    }
    private _curPNode: any = null;
    private _popZIndex: any = null;
    constructor() {
        this._popZIndex = 10000;
    }

    /**
     *  初始化
     */
    public init(node: any): void {
        this._curPNode = node;
    }
    /**
     *  显示
     */
    public show(str: string, time?: number) {
        if (!this._curPNode) {
            this._curPNode = cc.find("Canvas");
        }
        if (this._curPNode.getChildByName("CommonTip")) {
            this._curPNode.getChildByName("CommonTip").destroy();
        }
        try {
            const pop = cc.instantiate(resLoad.dirResList["dirRes"]["COMMONTIP"]);
            pop.time = time;
            pop.getComponent("CommonTip").setNote(str);
            this._curPNode.addChild(pop,this._popZIndex);
        } catch (err) {
            console.error(err);
        }
    }
    /**
     *  隐藏
     */
    public hide(): void {
        if (this._curPNode.getChildByName("CommonTip")) {
            this._curPNode.getChildByName("CommonTip").destroy();
        }
    }
}