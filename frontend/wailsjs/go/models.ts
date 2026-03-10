export namespace main {
	
	export class PausedSessionResponse {
	    sessionType: string;
	    remainingSeconds: number;
	    cyclePosition: number;
	    pausedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new PausedSessionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionType = source["sessionType"];
	        this.remainingSeconds = source["remainingSeconds"];
	        this.cyclePosition = source["cyclePosition"];
	        this.pausedAt = source["pausedAt"];
	    }
	}

}

export namespace storage {
	
	export class DailyStats {
	    Date: string;
	    TotalSessions: number;
	    TotalFocusMinutes: number;
	
	    static createFrom(source: any = {}) {
	        return new DailyStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Date = source["Date"];
	        this.TotalSessions = source["TotalSessions"];
	        this.TotalFocusMinutes = source["TotalFocusMinutes"];
	    }
	}
	export class UserSettings {
	    FocusDuration: number;
	    ShortBreakDuration: number;
	    LongBreakDuration: number;
	    SessionsBeforeLongBreak: number;
	    NotificationSound: string;
	    Mute: boolean;
	    Theme: string;
	
	    static createFrom(source: any = {}) {
	        return new UserSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FocusDuration = source["FocusDuration"];
	        this.ShortBreakDuration = source["ShortBreakDuration"];
	        this.LongBreakDuration = source["LongBreakDuration"];
	        this.SessionsBeforeLongBreak = source["SessionsBeforeLongBreak"];
	        this.NotificationSound = source["NotificationSound"];
	        this.Mute = source["Mute"];
	        this.Theme = source["Theme"];
	    }
	}

}

export namespace timer {
	
	export class TimerState {
	    status: string;
	    sessionType: string;
	    remainingSeconds: number;
	    totalSeconds: number;
	    cyclePosition: number;
	
	    static createFrom(source: any = {}) {
	        return new TimerState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.sessionType = source["sessionType"];
	        this.remainingSeconds = source["remainingSeconds"];
	        this.totalSeconds = source["totalSeconds"];
	        this.cyclePosition = source["cyclePosition"];
	    }
	}

}

