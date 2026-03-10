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
	    focusDuration: number;
	    shortBreakDuration: number;
	    longBreakDuration: number;
	    sessionsBeforeLongBreak: number;
	    notificationSound: string;
	    mute: boolean;
	    theme: string;
	
	    static createFrom(source: any = {}) {
	        return new UserSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.focusDuration = source["focusDuration"];
	        this.shortBreakDuration = source["shortBreakDuration"];
	        this.longBreakDuration = source["longBreakDuration"];
	        this.sessionsBeforeLongBreak = source["sessionsBeforeLongBreak"];
	        this.notificationSound = source["notificationSound"];
	        this.mute = source["mute"];
	        this.theme = source["theme"];
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

