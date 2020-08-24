import {Net} from "./Event/Net";
import {CommonCfg} from "./Common/CommonCfg";
import msgPb = require('./Proto/Common/msg_pb');
const {ccclass, property} = cc._decorator;

@ccclass
export default class Helloworld extends cc.Component {

    @property({
        type: dragonBones.ArmatureDisplay,
        tooltip: "龙骨",
    })
    player = null;

    onLoad() {
        // cc.director.getPhysicsManager().enabled = true; // 开启物理系统
        // cc.director.getPhysicsManager().debugDrawFlags = cc.PhysicsManager.DrawBits.e_aabbBit | // 物理系统调试信息
        //     cc.PhysicsManager.DrawBits.e_pairBit |
        //     cc.PhysicsManager.DrawBits.e_centerOfMassBit |
        //     cc.PhysicsManager.DrawBits.e_jointBit |
        //     cc.PhysicsManager.DrawBits.e_shapeBit
        // ;
        // cc.log(this.player);
        Net.getInstance().connect(CommonCfg.HALL_HOST, () => {
            cc.log("连接上了");
            const data = new msgPb.Code();
            data.setCode(msgPb.CodeType.SUC);
            data.setMsg("测试");
            Net.getInstance().send(msgPb.Event.EVENT_MSG_INFO, data);
        });
    }

    start () {

    }
}
