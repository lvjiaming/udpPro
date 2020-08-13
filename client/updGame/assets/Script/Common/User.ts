export class User {
    private static user: User = null;
    public static getInstance(): User {
        if (!this.user) {
            this.user = new User();
        }
        return this.user;
    }
    private _userId  = null;
    public set userid (id: any) {
        this._userId = id;
    }
    public get userid (): any {
        return this._userId;
    }

    private _name  = null;
    public set name (n: string) {
        this._name = n;
    }
    public get name (): string {
        return this._name;
    }

    private _weekInfo: number = null;
    public set weekInfo (val: number) {
        this._weekInfo = val
    }
    public get weekInfo (): number {
        return this._weekInfo;
    }

    private _monInfo: number = null;
    public set monInfo (val: number) {
        this._monInfo = val
    }
    public get monInfo (): number {
        return this._monInfo;
    }

    private _yearInfo: number = null;
    public set yearInfo (val: number) {
        this._yearInfo = val
    }
    public get yearInfo (): number {
        return this._yearInfo;
    }

    /**
     *  初始化玩家信息
     * @param userInfo
     */
    public init(userInfo: any): void {
        this.name = userInfo.getName();
        this.userid = userInfo.getId();
    }

    /**
     *  更新统计信息
     * @param data
     */
    public updateStatisyicalInfo(data: any): void {
        this.weekInfo = data.getWeekval();
        this.monInfo = data.getMonval();
        this.yearInfo = data.getYearval();
    }
}