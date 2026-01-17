export interface Setting {
    key: string;
    setting_type: string;
    value: string;
}

const GATEWAY_URL = process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8000";

export async function getSettings(): Promise<Setting[]> {
    try {
        const res = await fetch(`${GATEWAY_URL}/settings`);
        if (!res.ok) {
            console.error("Failed to fetch settings:", res.statusText);
            return [];
        }
        const json = await res.json();
        return json.data || [];
    } catch (error) {
        console.error("Error fetching settings:", error);
        return [];
    }
}

export async function getSettingByKey(key: string): Promise<Setting | null> {
    const settings = await getSettings();
    return settings.find(s => s.key === key) || null;
}
