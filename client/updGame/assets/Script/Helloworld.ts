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
    }

    start () {

    }
}
