const {ccclass, property} = cc._decorator;

@ccclass
export default class NewClass extends cc.Component {
    @property({
        type: cc.Label,
        tooltip: "内容",
    })
    note = null;

    onLoad () {
        if ((<any>this.node).time) {
            this.scheduleOnce(() => {
                this.node.destroy();
            }, (<any>this.node).time);
        }
    }

    start () {

    }

    public setNote(str: string): void {
        if (this.note) {
            this.note.string = `${str}...`;
        }
    }
}
