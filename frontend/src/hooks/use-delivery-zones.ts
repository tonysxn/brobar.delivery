import { useState, useEffect } from 'react';
import { useSettings } from '@/contexts/settings-context';

export interface DeliveryZone {
    radius: number;
    innerRadius: number;
    color: string;
    price: number;
    freeOrderPrice: number;
    name: string;
}

const DEFAULT_ZONES: DeliveryZone[] = [
    { radius: 2, innerRadius: 0, color: "#22c55e", price: 150, freeOrderPrice: 600, name: "Зелена зона" },
    { radius: 5, innerRadius: 2, color: "#eab308", price: 200, freeOrderPrice: 1100, name: "Жовта зона" },
    { radius: 7, innerRadius: 5, color: "#f97316", price: 300, freeOrderPrice: 1800, name: "Помаранчева зона" },
    { radius: 10, innerRadius: 7, color: "#ef4444", price: 350, freeOrderPrice: 2400, name: "Червона зона" },
];

export function useDeliveryZones() {
    const { getSetting, loading } = useSettings();
    const [zones, setZones] = useState<DeliveryZone[]>(DEFAULT_ZONES);

    useEffect(() => {
        if (loading) return;

        const setting = getSetting('delivery_zones');
        if (setting && setting.value) {
            try {
                const parsed = JSON.parse(setting.value);
                if (Array.isArray(parsed)) {
                    setZones(parsed);
                }
            } catch (error) {
                console.error("Failed to parse delivery_zones", error);
            }
        }
    }, [loading, getSetting]);

    return { zones, loading };
}
