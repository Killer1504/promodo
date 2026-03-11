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

export namespace stats {
	
	export class DailyStatsResponse {
	    date: string;
	    totalSessions: number;
	    totalFocusMinutes: number;
	    isToday: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DailyStatsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.totalSessions = source["totalSessions"];
	        this.totalFocusMinutes = source["totalFocusMinutes"];
	        this.isToday = source["isToday"];
	    }
	}

}

export namespace storage {
	
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

