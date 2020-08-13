import {resLoad} from "../Common/ResLoad";
interface InfoClass<T extends InfoBaseClass>{
    _data: any
    showDelUI()
    hideDelUI()
    selectDel()
    unSelectDel();
    showChangeUI()
    hideChangeUI()
}

export class InfoBaseClass extends cc.Component {
    public _data: any = null;
}

export enum INFO_STATE {
    NONE = 1,
    DEL,
    CHANGE,
}

export class InfoManager {
    public static infoMgr: InfoManager = null;
    public static getInstance(): InfoManager {
        if (!this.infoMgr) {
            this.infoMgr = new InfoManager();
        }
        return this.infoMgr;
    }

    private _curInfoState: INFO_STATE = INFO_STATE.NONE; // 当前的消息状态

    private _curDelOptNode: cc.Node = null; // 删除操作的节点

    private _curShowInfoList: Array<InfoClass<T>> = []; // 当前显示的消息节点

    private _curAllInfo: Array<any> = []; // 当前的所有消息

    private _curDelList: Array<number> = []; // 当前删除的列表

    private _curSelectInfo: InfoClass<T> = null;

    /**
     *  显示气泡
     * @param info
     */
    public showTip<T extends InfoBaseClass>(info: InfoClass<T>, cb?: any): void {
        const infoNode: cc.Node = (<any>info).node;
        const qipao = cc.instantiate(resLoad.dirResList["dirRes"]["EDITQIPAO"]);
        infoNode.addChild(qipao);
        const worldPos = cc.find("Canvas").convertToWorldSpaceAR(cc.v2(0, 0));
        const nodePos = infoNode.convertToNodeSpaceAR(worldPos);
        qipao.position = nodePos;
        (<any>qipao).owner = info;
        if (cb && cb instanceof Function) {
            const qp = qipao.getChildByName("Bg");
            cb(qp)
        }
    }

    /**
     *  开始删除信息
     */
    public startDelInfo(): void {
        this._curDelList = [];
        this._curInfoState = INFO_STATE.DEL;
        if (this._curDelOptNode) {
            this._curDelOptNode.active = true;
            const allBtn = this._curDelOptNode.getChildByName("BtnAll");
            if (allBtn) {
                if (allBtn.getComponent(cc.Toggle)) {
                    allBtn.getComponent(cc.Toggle).isChecked = false;
                }
            }
        }
        this._curShowInfoList.forEach((item: InfoClass<T>) => {
            item.showDelUI();
        });
        cc.log(this._curShowInfoList);
    }

    /**
     *  设置所有消息删除
     */
    public setAllInfoDel(state: boolean): void {
        this._curDelList = [];
        if (state) {
            this._curAllInfo.forEach((item) => {
                this._curDelList.push(item.getId());
            });
        }
        this._curShowInfoList.forEach((item: InfoClass<T>) => {
            if (state) {
                item.selectDel();
            } else {
                item.unSelectDel();
            }
        });
    }

    /**
     *  设置当前显示的消息列表
     * @param list
     */
    public addDelInfo(infoId: number): void {
        if (!this.checkIsInDelList(infoId)) {
            this._curDelList.push(infoId);
        }
    }

    /**
     *  移除
     */
    public removeDelInfo(infoId: number): void {
        for (let index = this._curDelList.length - 1; index >= 0 ; index--) {
            const curInfo = this._curDelList[index];
            if (curInfo == infoId) {
                this._curDelList.splice(index, 1);
                break;
            }
        }
    }

    /**
     *  检查是否在列表里面
     */
    public checkIsInDelList(infoId: number): boolean {
        let has = false;
        this._curDelList.forEach((item) => {
            if (item == infoId) {
                has = true;
            }
        });
        return has;
    }

    /**
     *  返回删除列表
     * @returns {Array<InfoClass<T>>|Array}
     */
    public getCurDelList(): Array<number> {
        return this._curDelList || [];
    }

    /**
     *  取消删除信息
     */
    public cancelDelInfo(): void {
        this._curDelList = [];
        this._curInfoState = INFO_STATE.NONE;
        if (this._curDelOptNode) {
            this._curDelOptNode.active = false;
        }
        this._curShowInfoList.forEach((item: InfoClass<T>) => {
            item.hideDelUI();
        });
    }

    /**
     *  设置删除操作节点
     */
    public setDelOptNode(node: cc.Node): void {
        this._curDelOptNode = node;
    }

    /**
     *  设置当前显示的消息列表
     * @param list
     */
    public addShowInfo<T extends InfoBaseClass>(infoClsss: InfoClass<T>): void {
        this._curShowInfoList.push(infoClsss);
        for (let index = this._curShowInfoList.length - 1; index >= 0 ; index--) {
            const curInfo = this._curShowInfoList[index];
            if (!curInfo._data) {
                this._curShowInfoList.splice(index, 1);
            }
        }
        // cc.log(this._curShowInfoList);
    }

    /**
     *  设置当前显示的消息列表
     * @param infoList
     */
    public setShowInfo(infoList: Array<T>): void {
        this._curShowInfoList = infoList;
    }

    /**
     *  移除
     */
    public removeShowInfo<T extends InfoBaseClass>(infoClss: InfoClass<T>): void {
        for (let index = this._curShowInfoList.length - 1; index >= 0 ; index--) {
            const curInfo = this._curShowInfoList[index];
            if (infoClss._data && curInfo._data && infoClss._data.getId() == curInfo._data.getId()) {
                this._curShowInfoList.splice(index, 1);
                cc.log("移除一个");
                break;
            }
        }
    }

    /**
     *  返回状态
     * @returns {INFO_STATE}
     */
    public getInfoState(): INFO_STATE {
        return this._curInfoState;
    }

    /**
     *  返回状态
     * @returns {INFO_STATE}
     */
    public initInfoState(): void {
        this._curInfoState = INFO_STATE.NONE;
        this._curShowInfoList.forEach((item: InfoClass<T>) => {
            item.hideDelUI();
        });
        this._curDelList = [];
    }

    /**
     *  设置当前的所有消息
     * @param infoList
     */
    public setCurAllInfo(infoList: Array<any>): void {
        this._curAllInfo = infoList || [];
    }

    /**
     *  开始修改信息
     */
    public startChange<T extends InfoBaseClass>(info: InfoClass<T>): void {
        this._curInfoState = INFO_STATE.CHANGE;
        if (this._curSelectInfo) {
            this._curSelectInfo.hideChangeUI();
        }
        this._curSelectInfo = info;
        this._curSelectInfo.showChangeUI();
    }

    /**
     *  结束修改
     */
    public endChange(): void {
        this._curInfoState = INFO_STATE.NONE;
        if (this._curSelectInfo) {
            this._curSelectInfo.hideChangeUI();
        }
        this._curSelectInfo = null;
    }

    /**
     *  初始化数据
     */
    public initData(): void {
        this.initInfoState();
        this._curDelList = [];
        this._curShowInfoList = [];
        this._curAllInfo = [];
        this._curDelOptNode = null;
    }
}
