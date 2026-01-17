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

const DAY_ORDER = ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'];
const DAY_NAMES_UA: Record<string, string> = {
    monday: 'Пн',
    tuesday: 'Вт',
    wednesday: 'Ср',
    thursday: 'Чт',
    friday: 'Пт',
    saturday: 'Сб',
    sunday: 'Нд',
};

function getScheduleKey(schedule: DaySchedule): string {
    if (schedule.closed) return 'closed';
    return `${schedule.start}-${schedule.end}`;
}

function formatSchedule(schedules: DaySchedules | null): string[] {
    if (!schedules) return [];

    // Group consecutive days with the same schedule
    const groups: { days: string[]; schedule: DaySchedule }[] = [];

    for (const day of DAY_ORDER) {
        const schedule = schedules[day];
        if (!schedule) continue;

        const key = getScheduleKey(schedule);
        const lastGroup = groups[groups.length - 1];

        if (lastGroup && getScheduleKey(lastGroup.schedule) === key) {
            lastGroup.days.push(day);
        } else {
            groups.push({ days: [day], schedule });
        }
    }

    // Format each group
    return groups.map(group => {
        const firstDay = group.days[0];
        const lastDay = group.days[group.days.length - 1];

        let dayRange: string;
        if (group.days.length === 7) {
            dayRange = 'Пн-Нд';
        } else if (group.days.length === 1) {
            dayRange = DAY_NAMES_UA[firstDay];
        } else {
            dayRange = `${DAY_NAMES_UA[firstDay]}-${DAY_NAMES_UA[lastDay]}`;
        }

        const timeRange = group.schedule.closed
            ? 'Вихідний'
            : `${group.schedule.start}-${group.schedule.end}`;

        return `${dayRange}: ${timeRange}`;
    });
}

export function formatWorkingHours(workingHours: WorkingHours | null, type: 'delivery' | 'pickup' = 'delivery'): string[] {
    if (!workingHours) return [];
    return formatSchedule(workingHours[type]);
}

export function formatDeliveryHours(workingHours: WorkingHours | null): string[] {
    return formatWorkingHours(workingHours, 'delivery');
}

export function formatPickupHours(workingHours: WorkingHours | null): string[] {
    return formatWorkingHours(workingHours, 'pickup');
}

export type { WorkingHours, DaySchedule, DaySchedules };
