import { useState, useEffect } from 'react';
import { useSettings } from '@/contexts/settings-context';

interface DaySchedule {
    start: string;
    end: string;
    closed: boolean;
}

interface DaySchedules {
    [key: string]: DaySchedule;
}

interface WorkingHours {
    delivery: DaySchedules;
    pickup: DaySchedules;
}

interface ServerTime {
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

const DAYS = ['sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];
const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

export function useShopStatus() {
    const { getSetting, loading: settingsLoading } = useSettings();
    const [serverTime, setServerTime] = useState<ServerTime | null>(null);
    const [status, setStatus] = useState<ShopStatus>({
        isOpen: true,
        isPaused: false,
        message: '',
        workingHours: null,
        serverTime: null,
        deliveryOpen: true,
        pickupOpen: true,
    });

    // Fetch server time
    useEffect(() => {
        async function fetchServerTime() {
            try {
                const res = await fetch(`${GATEWAY_URL}/time`);
                if (res.ok) {
                    const json = await res.json();
                    if (json.success) {
                        setServerTime(json.data);
                    }
                }
            } catch (error) {
                console.error("Failed to fetch server time", error);
            }
        }

        fetchServerTime();
        // Poll every minute
        const interval = setInterval(fetchServerTime, 60000);
        return () => clearInterval(interval);
    }, []);

    // Process status when settings and server time are available
    useEffect(() => {
        if (settingsLoading || !serverTime) return;

        const hoursSetting = getSetting('working_hours');
        const pausedSetting = getSetting('sales_paused');

        let isPaused = false;
        if (pausedSetting && pausedSetting.value === 'true') {
            isPaused = true;
        }

        let workingHours: WorkingHours | null = null;
        if (hoursSetting && hoursSetting.value) {
            try {
                workingHours = JSON.parse(hoursSetting.value);
            } catch (e) {
                console.error("Failed to parse working_hours", e);
            }
        }

        const dayName = DAYS[serverTime.day_number];
        let deliveryOpen = true;
        let pickupOpen = true;
        let message = '';

        if (isPaused) {
            deliveryOpen = false;
            pickupOpen = false;
            message = "Вибачте, ми тимчасово не приймаємо замовлення.";
        } else if (workingHours) {
            const currentTime = serverTime.time.substring(0, 5); // "HH:MM"

            // Check delivery hours
            const deliverySchedule = workingHours.delivery?.[dayName];
            if (deliverySchedule) {
                if (deliverySchedule.closed) {
                    deliveryOpen = false;
                } else {
                    if (currentTime < deliverySchedule.start || currentTime > deliverySchedule.end) {
                        deliveryOpen = false;
                    }
                }
            }

            // Check pickup hours
            const pickupSchedule = workingHours.pickup?.[dayName];
            if (pickupSchedule) {
                if (pickupSchedule.closed) {
                    pickupOpen = false;
                } else {
                    if (currentTime < pickupSchedule.start || currentTime > pickupSchedule.end) {
                        pickupOpen = false;
                    }
                }
            }

            if (!deliveryOpen && !pickupOpen) {
                const schedule = deliverySchedule || pickupSchedule;
                if (schedule && !schedule.closed) {
                    message = `Ми працюємо з ${schedule.start} до ${schedule.end}`;
                } else {
                    message = "Сьогодні ми не працюємо";
                }
            }
        }

        const isOpen = deliveryOpen || pickupOpen;

        setStatus({
            isOpen,
            isPaused,
            message,
            workingHours,
            serverTime,
            deliveryOpen,
            pickupOpen
        });
    }, [settingsLoading, serverTime, getSetting]);

    return status;
}
