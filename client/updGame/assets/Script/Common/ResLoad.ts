const resLoad = {
    dirResList: [],
    loadDirRes: (path: string, callBack: any): void => {
        cc.loader.loadResDir(path, (err, dir) => {
            if (err) {
                cc.error(`error: ${err}`);
                return;
            }
            const pathList = path.split("/");
            if (!resLoad.dirResList[pathList[pathList.length - 1]]) {
                resLoad.dirResList[pathList[pathList.length - 1]] = {};
            }
            for (let item in dir) {
                if (dir[item] instanceof cc.Prefab) {
                    resLoad.dirResList[pathList[pathList.length - 1]][dir[item].name.toUpperCase()] = dir[item];
                } else if (dir[item] instanceof cc.SpriteAtlas) {
                    xx.sys.objectToArray(dir[item]._spriteFrames).forEach((tex) => {
                        tex._texture.notRelease = true;
                    });
                    resLoad.dirResList[pathList[pathList.length - 1]][dir[item].name.toUpperCase().split(".")[0]] = dir[item];
                }
            }
            if (callBack) {
                callBack();
            }
        });
    },

    releaseDirRes: (path: string): void => {
        const pathList = path.split("/");
        cc.loader.releaseResDir(path);
        resLoad.dirResList[pathList[pathList.length - 1]] = {};
    }
};
export {
    resLoad,
}