export interface DaySchedule {
    start: string;
    end: string;
    closed: boolean;
}

export interface DaySchedules {
    [key: string]: DaySchedule;
}

export interface WorkingHours {
    delivery: DaySchedules;
    pickup: DaySchedules;
}

export interface ServerTime {
    timestamp: number;
    datetime: string;
    date: string;
    time: string;
    day: string;
    day_number: number;
}

export interface ShopStatus {
    isOpen: boolean;
    isPaused: boolean;
    message: string;
    workingHours: WorkingHours | null;
    serverTime: ServerTime | null;
    deliveryOpen: boolean;
    pickupOpen: boolean;
}
